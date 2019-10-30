package test

import (
	"bigboofer/helpers"

	"strings"
	"testing"
)

func TestGetRelativeProjPath(t *testing.T) {
	expected := "database/schema.sql"
	actual := helpers.GetRelativeProjPath("database", "schema.sql")

	if !strings.HasSuffix(actual, expected) {
		t.Errorf("Expected actual %v to end with %v", actual, expected)
	}
}
