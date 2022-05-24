package scoring

import (
	"testing"

	"bni.co.id/swift-filtering/database"
)

func init() {
	database.SetupFactory("mssql", "sqlserver://sa:Uat46@10.70.152.54:1433?database=SWIFT_FILTERING")
	defer database.Close()
}

func TestContainsScoring(t *testing.T) {
	proc := NewContainsProcessor()
	result := proc.Calculate("IRAN", 80)

	t.Log(result)

	result = proc.Calculate("BELARUS", 80)

	t.Log(result)
}
