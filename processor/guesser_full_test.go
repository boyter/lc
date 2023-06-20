// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"encoding/json"
	"os"
	"testing"
)

// Represents what the JSON looks like on disk enough for loading
type LicenseJson struct {
	LicenseText             string `json:"licenseText"`
	StandardLicenseTemplate string `json:"standardLicenseTemplate"`
	StandardLicenseHeader   string `json:"standardLicenseHeader"`
	Name                    string `json:"name"`
	LicenseId               string `json:"licenseId"`
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
		guesses := lg.KeyWordGuessLicence([]byte(l.LicenseText))

		if len(guesses) == 0 {
			fail++
			t.Error("expected", l.LicenseId)
			continue
		}

		if guesses[0].LicenseId != l.LicenseId {
			t.Error("expected", l.LicenseId, "got", guesses[0].LicenseId)
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
