package scoring

import (
	"testing"

	"bni.co.id/swift-filtering/database"
)

func init() {
	database.SetupFactory("mssql", "sqlserver://sa:Uat46@10.70.152.54:1433?database=SWIFT_FILTERING")
	defer database.Close()
}

func TestFuzzyMatching(t *testing.T) {
	proc := NewFuzzyProcessor()
	result := proc.Calculate("republic of belarusia", 70)

	t.Log(result)

	result = proc.Calculate("islamic republic of iran", 70)

	t.Log(result)
}
