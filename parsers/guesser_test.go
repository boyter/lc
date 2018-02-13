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

	if len(actual) != 370 {
		t.Log(len(actual))
		t.Errorf("Expected database to have 370 elements")
	}
}

func TestWalkDirectory(t *testing.T) {
	actual := walkDirectory("../examples/identifier/", [][]LicenseMatch{})

	if len(actual) != 3 {
		t.Errorf("Expected 3 results for directory")
	}
}

func TestProcessFile(t *testing.T) {
	actual := processFile("../examples/identifier/", "has_identifier.py", []LicenseMatch{})

	if actual.Md5Hash != "fc7d75e0bc0275841de8426b18791fa4" {
		t.Errorf("Expected MD5 to match")
	}

	if actual.Sha1Hash != "03a614cc51e9a783a695bcf99ec4adcdac34e1cc" {
		t.Errorf("Expected SHA1 to match")
	}

	if actual.Sha256Hash != "5b6bf8d45b25a0dab4f0817324618e425037f156b67cdc7503da01e6d9beb652" {
		t.Errorf("Expected SHA256 to match")
	}

	if len(actual.LicenseIdentified) != 2 {
		t.Errorf("Expected 2 identified licenses")
	}

	if actual.LicenseIdentified[0].LicenseId != "GPL-2.0" {
		t.Errorf("Expected license not identified")
	}

	if actual.LicenseIdentified[1].LicenseId != "GPL-3.0+" {
		t.Errorf("Expected license not identified")
	}
}
