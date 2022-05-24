package scoring

import (
	"log"
	"strings"
	"time"
	
	"github.com/jinzhu/gorm"
	
	"bni.co.id/swift-filtering/database"
	"bni.co.id/swift-filtering/myjson"
)

// FuzzyProcessor process calculation using similarity algorithm
type FuzzyProcessor struct {
	db *gorm.DB
}

// NewFuzzyProcessor Create new FuzzyProcessor
func NewFuzzyProcessor() *FuzzyProcessor {
	var proc *FuzzyProcessor = &FuzzyProcessor{
		db: database.Open(),
	}
	
	return proc
}

func (proc *FuzzyProcessor) getCache(word string) ([]CalculatedSanctionData, bool) {
	var result []CalculatedSanctionData
	var dataCache string
	var dWord string
	
	row := proc.db.Raw("SELECT word, data_cache FROM swift_fuzzy_cache WHERE word = ?",
		word).Row()
	
	row.Scan(&dWord, &dataCache)
	
	if dWord == "" {
		return result, false
	}
	
	myjson.Decode(dataCache, &result)
	
	return result, true
}

func (proc *FuzzyProcessor) writeCache(word string, dataCache []CalculatedSanctionData) {
	go func() {
		sql := "INSERT INTO swift_fuzzy_cache (word, data_cache, created_at) VALUES(?, ?, ?)"
		proc.db.Exec(sql, word, myjson.Encode(dataCache), time.Now())
	}()
}

// Calculate get fuzzy scores
func (proc *FuzzyProcessor) Calculate(value string, minimumScore float64) []CalculatedSanctionData {
	value = cleanString(value)
	var words = splitValue(removeNonAlphaNumerics(value))
	var result []CalculatedSanctionData
	
	if !validMinimumLength(value) {
		return result
	}
	
	var contains *ContainsProcessor = NewContainsProcessor()
	rs := contains.Calculate(value, minimumScore)
	
	for _, r := range rs {
		l1 := 100.0 * float64(len(r.ItemValue)) / float64(len(r.MatchedWord))
		l2 := 100.0 * float64(len(r.MatchedWord)) / float64(len(r.ItemValue))
		
		if l1 >= minimumScore && l1 >=0 && l1 <= 100.0{
			r.Score = l1
			result = append(result, r)
		} else if l2 >= minimumScore && l2 >= 0 && l2 <= 100.0 {
			r.Score = l2
			result = append(result, r)
		}
	}
	
	for _, word := range words {
		word = strings.ToUpper(word)
		
		if !validMinimumLength(word) {
			continue
		}
		
		cache, cacheExist := proc.getCache(word)
		
		if cacheExist {
			result = append(result, cache...)
			continue
		}
		
		var lists []CalculatedSanctionData = nil
		
		start := time.Now()
		proc.db.Raw(`SELECT d.*,
			? matched_word,
			dbo.SWIFT_StringSimilarity(UPPER(item_value) , ?) score
			FROM swift_sanction_data d
			WHERE dbo.SWIFT_StringSimilarity(UPPER(item_value) , ?) >= ?`,
			word, word, word, minimumScore).
			Scan(&lists)
		
		result = append(result, lists...)
		proc.writeCache(word, lists)
		
		log.Printf("%s took %v \n", word, time.Since(start))
	}
	
	return removeDuplicateResult(result)
}
