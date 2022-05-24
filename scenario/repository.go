package scenario

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"bni.co.id/swift-filtering/database"
)

// Repository Scenario repository is bridge to access scenario database
type Repository struct {
	db *gorm.DB
}

// RuleVO Rule and items model
type RuleVO struct {
	Rule
	RuleItems *[]RuleItem
}

var rulesCache = make(map[string][]RuleVO)

// NewRepository Create new instance of Repository
func NewRepository() *Repository {
	var repo *Repository = &Repository{
		db: database.Open(),
	}

	return repo
}

// GetRuleItemsByRuleID get Active RuleItem by RuleID
func (repo *Repository) GetRuleItemsByRuleID(ruleID uuid.UUID) []RuleItem {
	var result []RuleItem

	repo.db.Raw(`SELECT ri.*
				FROM swift_rule_items ri
				WHERE ri.rule_id = ?
				AND active_flag = 'Y'
				AND scoring_method IS NOT NULL`, ruleID).
		Scan(&result)

	return result
}

// GetRulesByMessageType Get Active Rules and Details by Message Type
func (repo *Repository) GetRulesByMessageType(mt string) []RuleVO {
	var rules, exist = rulesCache[mt]

	if exist {
		return rules
	}

	var result []RuleVO

	repo.db.Raw(`SELECT r.*
				FROM swift_rules r
				WHERE r.active_flag = 'Y'
					AND r.id IN(
						SELECT rm.rule_id
						FROM swift_message_type_rules rm
						WHERE rm.message_type = ?
							AND rm.active_flag = 'Y')`, mt).
		Scan(&result)

	for i := range result {
		items := repo.GetRuleItemsByRuleID(result[i].ID)
		result[i].RuleItems = &items
	}

	rulesCache[mt] = result

	return result
}

// ClearRulesCache release rules cache from memory
func ClearRulesCache() {
	rulesCache = make(map[string][]RuleVO)
}

// SaveSuspectedResult save suspected result data into db
func (repo *Repository) SaveSuspectedResult(suspected *Suspected) {
	repo.db.Create(suspected)
}
