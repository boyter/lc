package parsers

import (
	"testing"
)

func TestWalkDirectory(t *testing.T) {
	actual := walkDirectory("../examples/identifier/", [][]LicenseMatch{})

	if len(actual) != 3 {
		t.Errorf("Expected 3 results for directory")
	}
}
