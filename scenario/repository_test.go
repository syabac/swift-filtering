package scenario

import (
	"fmt"
	"testing"

	"bni.co.id/swift-filtering/database"
	"bni.co.id/swift-filtering/myjson"
	"github.com/google/uuid"
)

func init() {
	database.SetupFactory("mssql", "sqlserver://sa:Uat46@10.70.152.54:1433?database=SWIFT_FILTERING")
	defer database.Close()
}

var repo = NewRepository()

func TestGetActiveRules(t *testing.T) {
	rules := repo.GetRulesByMessageType("103")

	fmt.Println(myjson.EncodePretty(rules))
}

func TestGetActiveRuleItems(t *testing.T) {
	rules := repo.GetRuleItemsByRuleID(uuid.Nil)

	fmt.Println(myjson.EncodePretty(rules))
}
