package scenario

import (
	"time"

	"github.com/google/uuid"
)

// Rule Store rule metadata information
type Rule struct {
	ID          uuid.UUID `gorm:"type:uniqueidentifier;primary_key;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Name        string
	Description string
	ActiveFlag  string
}

// MessageTypeRule mapping between  MT XXX and Rule
type MessageTypeRule struct {
	ID          uuid.UUID `gorm:"type:uniqueidentifier;primary_key;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	MessageType string
	RuleID      uuid.UUID `gorm:"type:uniqueidentifier"`
	ActiveFlag  string
}

// RuleItem store detail rules per tags
type RuleItem struct {
	ID            uuid.UUID `gorm:"type:uniqueidentifier;primary_key;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	TagName       string
	ActiveFlag    string
	ScoringMethod string
	MinimumScore  float64
	RuleID        uuid.UUID `gorm:"type:uniqueidentifier"`
}

// Suspected transaction model
type Suspected struct {
	ID            uuid.UUID `gorm:"type:uniqueidentifier;primary_key;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	SwiftID       uuid.UUID  `gorm:"type:uniqueidentifier"`
	RuleID        uuid.UUID  `gorm:"type:uniqueidentifier"`
	TagName       string
	TagValue      string
	Score         float64
	SanctionValue string
	MatchedWord   string
	SanctionType  string
}
