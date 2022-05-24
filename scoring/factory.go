package scoring

import (
	"strings"
)

// GetScoringEngineByName get scoring engine based on methods (fuzzy, exact, contains)
func GetScoringEngineByName(name string) Processor {
	var processor Processor

	switch strings.ToUpper(name) {
	case "FUZZY":
		processor = NewFuzzyProcessor()
	case "EXACT":
		processor = NewExactProcessor()
	case "CONTAINS":
		processor = NewContainsProcessor()
	default:
		return nil
	}

	return processor
}
