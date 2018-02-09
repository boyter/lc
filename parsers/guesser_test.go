package parsers

import (
	"testing"
)

func expect(t *testing.T, expected string, actual string) {

}

func TestCleanText(t *testing.T) {
	actual := cleanText("ToLower")
	expected := "tolower"

	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	actual = cleanText("   ToLower999$%")
	expected = " tolower999 "

	if expected != actual {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestLoadDatabase(t *testing.T) {
	actual := loadDatabase()

	if len(actual) == 0 {
		t.Errorf("Expected database to not be empty")
	}
}
