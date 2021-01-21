// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
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

func loadLicences(t *testing.T) []LicenseJson {
	files, _ := ioutil.ReadDir("../assets/database/licenses/")

	var licenses []LicenseJson
	// Load all of the licenses from disk to use as a comparison for the check
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := ioutil.ReadFile(filepath.Join("../assets/database/licenses/", f.Name()))

			var license LicenseJson
			err := json.Unmarshal(bytes, &license)
			if err != nil {
				t.Error(err)
			}

			licenses = append(licenses, license)
		}
	}
	return licenses
}

func TestTesty(t *testing.T) {
	if 0 != 0 {
		t.Error("error...")
	}
}

//
//func TestKeywordCommonDatabase(t *testing.T) {
//	lg := NewLicenceGuesser(true, false)
//	lg.UseFullDatabase = false
//	licenses := loadLicences(t)
//	fail := 0
//	pass := 0
//
//	for _, l := range licenses {
//		guesses := lg.KeyWordGuessLicence([]byte(l.LicenseText))
//
//		if len(guesses) == 0 {
//			fail++
//			t.Error("expected", l.LicenseId)
//			continue
//		}
//
//		if guesses[0].LicenseId != l.LicenseId {
//			t.Error("expected", l.LicenseId, "got", guesses[0].LicenseId)
//			fail++
//		} else {
//			pass++
//		}
//	}
//
//	if fail != 0 {
//		t.Error(pass, "passes")
//		t.Error(fail, "fails")
//	}
//}
//
//func TestKeywordFullDatabase(t *testing.T) {
//	lg := NewLicenceGuesser(true, false)
//	lg.UseFullDatabase = true
//
//	licenses := loadLicences(t)
//	fail := 0
//	pass := 0
//
//	for _, l := range licenses {
//		guesses := lg.KeyWordGuessLicence([]byte(l.LicenseText))
//
//		if guesses[0].LicenseId != l.LicenseId {
//			t.Error("expected", l.LicenseId, "got", guesses[0].LicenseId)
//			fail++
//		} else {
//			pass++
//		}
//	}
//
//	if fail != 0 {
//		t.Error(pass, "passes")
//		t.Error(fail, "fails")
//	}
//}
//
//func TestVectorSpaceCommonDatabase(t *testing.T) {
//	lg := NewLicenceGuesser(false, true)
//	lg.UseFullDatabase = false
//
//	licenses := loadLicences(t)
//	fail := 0
//	pass := 0
//
//	for _, l := range licenses {
//		guesses := lg.VectorSpaceGuessLicence([]byte(l.LicenseText))
//
//		if len(guesses) == 0 {
//			fail++
//			t.Error("expected", l.LicenseId)
//			continue
//		}
//
//		if guesses[0].LicenseId != l.LicenseId {
//			t.Error("expected", l.LicenseId, "got", guesses[0].LicenseId)
//			fail++
//		} else {
//			pass++
//		}
//	}
//
//	if fail != 0 {
//		t.Error(pass, "passes")
//		t.Error(fail, "fails")
//	}
//}
//
//func TestVectorSpaceFullDatabase(t *testing.T) {
//	lg := NewLicenceGuesser(false, true)
//	lg.UseFullDatabase = true
//
//	licenses := loadLicences(t)
//	fail := 0
//	pass := 0
//
//	for _, l := range licenses {
//		guesses := lg.VectorSpaceGuessLicence([]byte(l.LicenseText))
//
//		if guesses[0].LicenseId != l.LicenseId {
//			t.Error("expected", l.LicenseId, "got", guesses[0].LicenseId)
//			fail++
//		} else {
//			pass++
//		}
//	}
//
//	if fail != 0 {
//		t.Error(pass, "passes")
//		t.Error(fail, "fails")
//	}
//}
