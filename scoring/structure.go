package scoring

import "bni.co.id/swift-filtering/swift"

// CalculatedSanctionData store calculated sanctiondata
type CalculatedSanctionData struct {
	swift.SanctionData
	Score float64
	MatchedWord string
}

// Processor all scoring processor must be implements this interface
type Processor interface {
	Calculate(value string, minimumScore float64) []CalculatedSanctionData
}
