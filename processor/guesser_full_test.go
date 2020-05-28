// SPDX-License-Identifier: MIT
// SPDX-License-Identifier: Unlicense

package processor

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
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

// Note that this takes a long time to run hence goroutines to try and speed it up
func TestKeywordFullDatabase(t *testing.T) {
	lg := NewLicenceGuesser(true, false)
	lg.UseFullDatabase = true

	files, _ := ioutil.ReadDir("../assets/database/licenses/")

	var licenses []LicenseJson
	// Load all of the licenses from disk
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

	fail := 0

	var wg sync.WaitGroup
	for _, l := range licenses {
		wg.Add(1)
		go func(l LicenseJson) {
			guesses := lg.KeyWordGuessLicence([]byte(l.LicenseText))

			if guesses[0].LicenseId != l.LicenseId {
				t.Error("expected", l.LicenseId, "got", guesses[0].LicenseId)
				fail++
			}
			wg.Done()
		}(l)
	}
	wg.Wait()

	if fail != 0 {
		t.Error(fail, "fails")
	}
}

// Note that this takes a long time to run hence goroutines to try and speed it up
func TestVectorSpaceFullDatabase(t *testing.T) {
	lg := NewLicenceGuesser(false, true)
	lg.UseFullDatabase = true

	files, _ := ioutil.ReadDir("../assets/database/licenses/")

	var licenses []LicenseJson
	// Load all of the licenses from disk
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

	fail := 0

	var wg sync.WaitGroup
	for _, l := range licenses {
		wg.Add(1)
		go func(l LicenseJson) {
			guesses := lg.VectorSpaceGuessLicence([]byte(l.LicenseText))

			if guesses[0].LicenseId != l.LicenseId {
				t.Error("expected", l.LicenseId, "got", guesses[0].LicenseId)
				fail++
			}
			wg.Done()
		}(l)
	}
	wg.Wait()

	if fail != 0 {
		t.Error(fail, "fails")
	}
}
