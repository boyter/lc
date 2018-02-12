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

func TestWalkDirectory(t *testing.T) {
	actual := walkDirectory("../examples/identifier/", [][]LicenseMatch{})

	if len(actual) != 3 {
		t.Errorf("Expected 3 results for directory")
	}
}

func TestProcessFile(t *testing.T) {
	actual := processFile("../examples/identifier/", "has_identifier.py", []LicenseMatch{})

	if actual.Md5Hash != "0ad2e6786423fa6933a49ae4f97ae79e" {
		t.Errorf("Expected MD5 to match")
	}

	if actual.Sha1Hash != "64904ca8a945009f95734d19198aaacd5e2db959" {
		t.Errorf("Expected SHA1 to match")
	}

	if actual.Sha256Hash != "bf201f35c6a8504b4d956b4403ff5a7fef490889d5166a34f01b653e4ce08a06" {
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
