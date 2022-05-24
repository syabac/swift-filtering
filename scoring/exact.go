package scoring

import (
	"strings"

	"bni.co.id/swift-filtering/database"
	"github.com/jinzhu/gorm"
)

// ExactProcessor exact match algorithm
type ExactProcessor struct {
	db *gorm.DB
}

// NewExactProcessor create new ExactProcessor
func NewExactProcessor() *ExactProcessor {
	var proc *ExactProcessor = &ExactProcessor{
		db: database.Open(),
	}

	return proc
}

// Calculate get fuzzy scores
func (proc *ExactProcessor) Calculate(value string, minimumScore float64) []CalculatedSanctionData {
	var result []CalculatedSanctionData

	proc.db.Raw(`SELECT d.*, 100 score
		FROM swift_sanction_data d
		WHERE item_value=?`,
		strings.TrimSpace(strings.ToUpper(value))).Scan(&result)

	return result
}
