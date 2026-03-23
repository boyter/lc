// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"encoding/json"
	"os"
	"testing"
)

// Represents what the JSON looks like on disk enough for loading
type LicenseJson struct {
	LicenseTexts []string `json:"licenseTexts"`
	LicenseIds   []string `json:"licenseIds"`
}

func loadLicences() []LicenseJson {
	bytes, _ := os.ReadFile("../database_keywords.json")
	var licenses []LicenseJson
	_ = json.Unmarshal(bytes, &licenses)
	return licenses
}

// Simple test to ensure that the stuff we are about to load actually works as we expect
func TestKeywordCommonDatabase(t *testing.T) {
	lg := NewLicenceGuesser(true, false)
	lg.UseFullDatabase = true
	licenses := loadLicences()
	fail := 0
	pass := 0

	for _, l := range licenses {
		if len(l.LicenseTexts) == 0 || len(l.LicenseIds) == 0 {
			continue
		}
		guesses := lg.KeyWordGuessLicence([]byte(l.LicenseTexts[0]))

		if len(guesses) == 0 {
			fail++
			t.Error("expected", l.LicenseIds[0])
			continue
		}

		if !ContainsString(guesses[0].LicenseIds, l.LicenseIds) {
			t.Error("expected", l.LicenseIds[0], "got", guesses[0].LicenseIds[0])
			fail++
		} else {
			pass++
		}
	}

	if fail != 0 {
		t.Error(pass, "passes")
		t.Error(fail, "fails")
	}
}
