package swift

import (
	"time"

	"github.com/google/uuid"
)

// SanctionData list of Sanction data
type SanctionData struct {
	ID           uuid.UUID `gorm:"type:uniqueidentifier;primary_key;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
	SanctionType string
	ItemValue    string
}

// Transaction Swift Transaction on database
type Transaction struct {
	ID           uuid.UUID `gorm:"type:uniqueidentifier;primary_key;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
	MessageType  string
	Header       string
	Body         string
	IsSuspected  int
	IsFollowedUp int
}

// TransactionValue Swift transaction detail values
type TransactionValue struct {
	ID        uuid.UUID `gorm:"type:uniqueidentifier;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	SwiftID   uint
	TagName   string
	TagValue  string
}
