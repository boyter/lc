package parsers

import (
	"testing"
)

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
