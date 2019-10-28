package test

import (
	"bigboofer/database"

	"testing"
)

func TestCanOnboardNewDB(t *testing.T) {
	database.OnboardDB()
}

func TestCanCreateNewDBConnection(t *testing.T) {
	database.OnboardDB()
	if database.GetDB() == nil {
		t.Errorf("Returned nil database object")
	}
}
