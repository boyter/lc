package parsers

import (
	// "io/ioutil"
	"path/filepath"
	// "strings"
	"testing"
)

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

// This is slow but ensures that things work as we expect for fuzzy matching
// func TestProcessFileLicensesFuzzy(t *testing.T) {
// 	files, _ := ioutil.ReadDir("../examples/licenses/")

// 	correctFuzzy := float64(0.00)
// 	correctKeywords := float64(0.00)

// 	for _, file := range files {
// 		actualFuzzy := processFile("../examples/licenses/", file.Name(), []LicenseMatch{})

// 		if len(actualFuzzy.LicenseGuesses) != 0 {
// 			if strings.Replace(file.Name(), ".json", "", 1) == actualFuzzy.LicenseGuesses[0].LicenseId {
// 				correctFuzzy++
// 			}
// 		}

// 		content := readFile(filepath.Join("../examples/licenses/", file.Name()))
// 		actualKeywords := keywordGuessLicense(string(content), loadDatabase())

// 		if len(actualKeywords) != 0 {
// 			if strings.Replace(file.Name(), ".json", "", 1) == actualKeywords[0].LicenseId {
// 				correctKeywords++
// 			}
// 		}
// 	}

// 	totalPercentageFuzzy := (correctFuzzy / float64(len(files))) * 100
// 	if totalPercentageFuzzy < 0.95 {
// 		t.Errorf("Not enough guesses correct fuzzy", correctFuzzy, len(files), totalPercentageFuzzy)
// 	}

// 	totalPercentageKeywords := (correctKeywords / float64(len(files))) * 100
// 	if totalPercentageKeywords < 0.95 {
// 		t.Errorf("Not enough guesses correct keywords", correctKeywords, len(files), totalPercentageKeywords)
// 	}
// }

func TestProcessFileLicensesTop10(t *testing.T) {
	// https://www.blackducksoftware.com/top-open-source-licenses
	files := []string{"MIT", "GPL-2.0", "Apache-2.0", "GPL-3.0-only", "ISC", "Artistic-2.0", "LGPL-2.1", "LGPL-3.0", "EPL-2.0", "MS-PL"}

	for _, file := range files {
		content := readFile(filepath.Join("../examples/licenses/", file+".json"))
		actual := keywordGuessLicense(string(content), loadDatabase())

		if len(actual) == 0 {
			t.Errorf("Expected some guesses for %s", file)
		}

		found := false
		for _, license := range actual {
			if license.LicenseId == file {
				found = true
			}
		}

		if found == false {
			t.Errorf("Expected license to be found %s", file)
		}
	}
}

func TestIdentifierGuessLicence(t *testing.T) {
	actual := identifierGuessLicence("test", loadDatabase())
	if len(actual) != 0 {
		t.Errorf("Should be no matches")
	}

	actual = identifierGuessLicence("# SPDX-License-Identifier: GPL-2.0", loadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = identifierGuessLicence("# SPDX-License-Identifier: GPL-2.0 ", loadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}

	actual = identifierGuessLicence("# SPDX-License-Identifier: GPL-2.0 \n # SPDX-License-Identifier: GPL-3.0+", loadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}
	if actual[1].LicenseId != "GPL-3.0+" {
		t.Errorf("Should match GPL-3.0+")
	}

	actual = identifierGuessLicence(`# SPDX-License-Identifier: GPL-2.0

import cherrypy

class Example(object):

    # SPDX-License-Identifier: GPL-3.0+
    @cherrypy.expose
    def index(self, **params):
        return 'Hello World'


if __name__ == '__main__':
    cherrypy.config.update({
        'server.socket_host': '0.0.0.0',
        'server.socket_port': 8080,
        'server.thread_pool': 30,
    })
    cherrypy.quickstart(Example())
`, loadDatabase())
	if actual[0].LicenseId != "GPL-2.0" {
		t.Errorf("Should match GPL-2.0")
	}
	if actual[1].LicenseId != "GPL-3.0+" {
		t.Errorf("Should match GPL-3.0+")
	}
}
