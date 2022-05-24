package swift

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"bni.co.id/swift-filtering/database"
)

// Repository  repository to manage swift data
type Repository struct {
	db *gorm.DB
}

// NewRepository Create new instance of Repository
func NewRepository() *Repository {
	var repo *Repository = &Repository{
		db: database.Open(),
	}

	return repo
}

// Save save Swift Transaction data
func (repo *Repository) Save(sm *Message, isSuspected int) uuid.UUID {
	var trx = &Transaction{
		MessageType:  sm.Block2.MessageType,
		Header:       sm.Block1.Values,
		Body:         sm.Body,
		IsFollowedUp: 0,
		IsSuspected:  isSuspected,
	}

	repo.db.Create(trx)
	//
	// for tag, values := range sm.Tags {
	// 	for _, value := range values {
	// 		var dtl = &TransactionValue{
	// 			SwiftID:  trx.ID,
	// 			TagName:  tag,
	// 			TagValue: value,
	// 		}
	// 		repo.db.Create(dtl)
	// 	}
	// }

	return trx.ID
}

// SaveLog save message to log
func (repo *Repository) SaveLog(sm *Message) {
	mdb := database.OpenMongoDb()

	mdb.Collection("transaction").InsertOne(context.Background(), sm)
}
