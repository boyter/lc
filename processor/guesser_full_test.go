package processor

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

type LicenseJson struct {
	LicenseText             string `json:"licenseText"`
	StandardLicenseTemplate string `json:"standardLicenseTemplate"`
	StandardLicenseHeader   string `json:"standardLicenseHeader"`
	Name                    string `json:"name"`
	LicenseId               string `json:"licenseId"`
}

// Note that this takes a long time to run hence goroutines to try and speed it up
func TestFullDatabase(t *testing.T) {
	lg := NewLicenceGuesser()
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

	var wg sync.WaitGroup
	for _, l := range licenses {
		wg.Add(1)
		go func(l LicenseJson) {
			guesses := lg.GuessLicense([]byte(l.LicenseText))

			if guesses[0].LicenseId != l.LicenseId {
				t.Error("expected", l.LicenseId, "got", guesses[0].LicenseId)
			}
			wg.Done()
		}(l)
	}
	wg.Wait()
}
