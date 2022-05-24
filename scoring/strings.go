package scoring

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/xrash/smetrics"

	"bni.co.id/swift-filtering/config"
	"bni.co.id/swift-filtering/util"
)

var regexSpace = regexp.MustCompile(`\s+`)
var regexNonWord = regexp.MustCompile(`\W`)

const minLengthWord = 3

func splitValue(value string) []string {
	return strings.Split(cleanString(value), " ")
}

func stringSimilarity(val1, val2 string) float64 {
	return util.RoundNumber(smetrics.JaroWinkler(strings.ToUpper(val1), strings.ToUpper(val2), 0.7, 4)*100, 2)
}

func cleanString(value string) string {
	return strings.TrimSpace(regexSpace.ReplaceAllString(value, " "))
}

func removeNonAlphaNumerics(value string) string {
	return strings.TrimSpace(cleanString(regexNonWord.ReplaceAllString(value, " ")))
}

func validMinimumLength(value string) bool {
	var minLen, _ = strconv.Atoi(config.GetSetting("swift.minimum.word.length"))
	if minLen < minLengthWord {
		minLen = minLengthWord
	}
	return len(value) >= minLen
}

// Remove duplicate calculation result
// only take with the higher score
func removeDuplicateResult(result []CalculatedSanctionData) []CalculatedSanctionData {
	if len(result) < 2 {
		return result
	}

	var check = make(map[uuid.UUID]CalculatedSanctionData)

	for _, item := range result {
		prev, isExist := check[item.ID]

		if !isExist {
			check[item.ID] = item
			continue
		}

		if item.Score > prev.Score {
			check[item.ID] = item
		}
	}

	var xcopy []CalculatedSanctionData
	for _, item := range check {
		xcopy = append(xcopy, item)
	}

	return xcopy
}
