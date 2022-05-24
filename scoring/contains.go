package scoring

import (
	"github.com/jinzhu/gorm"

	"bni.co.id/swift-filtering/database"
)

// ContainsProcessor Scoring engine using substring similarity
type ContainsProcessor struct {
	db *gorm.DB
}

// NewContainsProcessor Create new ContainsProcessor
func NewContainsProcessor() *ContainsProcessor {
	var proc *ContainsProcessor = &ContainsProcessor{
		db: database.Open(),
	}

	return proc
}

// Calculate check given value if contains in sanction lists
func (proc *ContainsProcessor) Calculate(value string, minimumScore float64) []CalculatedSanctionData {
	value = cleanString(value)
	var words = splitValue(removeNonAlphaNumerics(value))
	var result []CalculatedSanctionData

	if !validMinimumLength(value) {
		return result
	}

	var valueLike = "%" + value + "%"
	sql := `SELECT d.*, 100 score, ? matched_word
			FROM swift_sanction_data d
			WHERE  item_value LIKE ?
				OR ? LIKE CONCAT('%', item_value, '%')`

	proc.db.Raw(sql, value, valueLike, value).Scan(&result)
	checkWords := make(map[string]bool)

	for _, word := range words {
		_, wordExist := checkWords[word]
		checkWords[word] = true
		if wordExist || !validMinimumLength(word) {
			continue
		}

		var lists []CalculatedSanctionData

		valueLike = "%" + word + "%"

		proc.db.Raw(sql, word, valueLike, word).Scan(&lists)
		result = append(result, lists...)
	}

	return removeDuplicateResult(result)
}
