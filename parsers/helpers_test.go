package parsers

import (
	"testing"
)

func init() {
	LoadDatabase()
}
func TestGetMd5Hash(t *testing.T) {
	actual := getMd5Hash([]byte("ToLower"))
	expected := "82b2da23d045588f0a386e035a43effd"

	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestGetSha1Hash(t *testing.T) {
	actual := getSha1Hash([]byte("ToLower"))
	expected := "2458b21542ecbc0f90f413b7ee46521686d711b0"

	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestGetSha256Hash(t *testing.T) {
	actual := getSha256Hash([]byte("ToLower"))
	expected := "c64acc48c7d802c9418269506dc0efc01c68b18209ee2f61e518702dc5c135ac"

	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestLicenceListHasLicense(t *testing.T) {
	licenseList := []LicenseMatch{}
	licenseMatch := LicenseMatch{LicenseId: "test", Score: 0.0}

	actual := licenceListHasLicense(licenseMatch, licenseList)

	if actual != false {
		t.Errorf("Expected false but got true")
	}

	licenseList = append(licenseList, licenseMatch)

	actual = licenceListHasLicense(licenseMatch, licenseList)

	if actual != true {
		t.Errorf("Expected true but got false")
	}
}

func TestUniqLicenseMatch(t *testing.T) {
	licenseList := []LicenseMatch{}
	licenseMatch1 := LicenseMatch{LicenseId: "test", Score: 0.0}
	licenseMatch2 := LicenseMatch{LicenseId: "test", Score: 0.0}

	licenseList = append(licenseList, licenseMatch1)
	licenseList = append(licenseList, licenseMatch2)

	actual := uniqLicenseMatch(licenseList)

	if len(actual) != 1 {
		t.Errorf("Expected single result")
	}

	// if even one portion changes then it should be included
	licenseMatch3 := LicenseMatch{LicenseId: "test", Score: 0.1}
	licenseMatch4 := LicenseMatch{LicenseId: "test2", Score: 0.0}

	licenseList = append(licenseList, licenseMatch3)
	licenseList = append(licenseList, licenseMatch4)

	actual = uniqLicenseMatch(licenseList)

	if len(actual) != 3 {
		t.Errorf("Expected three results")
	}
}

func TestRandStringBytes(t *testing.T) {
	actual := randStringBytes(5)

	if len(actual) != 5 {
		t.Errorf("Expected length of 5 %s", actual)
	}

	actual = randStringBytes(10)

	if len(actual) != 10 {
		t.Errorf("Expected length of 10 %s", actual)
	}
}

func TestBytesToHuman(t *testing.T) {
	if "123B" != bytesToHuman(123) {
		t.Errorf("Expected 123B")
	}

	if "1K" != bytesToHuman(1025) {
		t.Errorf("Expected 1K")
	}

	if "1M" != bytesToHuman(1048576) {
		t.Errorf("Expected 1M")
	}

	if "1.1M" != bytesToHuman(1178576) {
		t.Errorf("Expected 1.1M")
	}

	if "1024M" != bytesToHuman(1073741823) {
		t.Errorf("Expected 1024M")
	}

	if "1G" != bytesToHuman(1073741824) {
		t.Errorf("Expected 1G")
	}

	if "1T" != bytesToHuman(1099511627776) {
		t.Errorf("Expected 1T")
	}
}
