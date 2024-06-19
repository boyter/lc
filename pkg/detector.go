// SPDX-License-Identifier: AGPL-3.0

package pkg

import (
	"encoding/base64"
	"encoding/json"
	"github.com/boyter/lc/pkg/levenshtein"
	"math"
	"strings"
)

type LicenceDetector struct {
	Database []License
}

func NewLicenceDetector() *LicenceDetector {
	l := new(LicenceDetector)
	l.LoadDatabase()
	return l
}

type License struct {
	LicenseTexts []string `json:"licenseTexts"` // examples of text that we have for these licences
	LicenseIds   []string `json:"licenseIds"`   // SPDX ids where licences are considered identical
	Keywords     []string `json:"keywords"`     // keywords that are unique and can be used to identify this group of licences
}

// LoadDatabase will initialize the database values and should only be called once such as in an init
func (l *LicenceDetector) LoadDatabase() {
	if len(l.Database) != 0 {
		return
	}

	data, _ := base64.StdEncoding.DecodeString(database_keywords)
	_ = json.Unmarshal(data, &l.Database)
}

type LicenseGuess struct {
	Name string
}

func (l *LicenceDetector) Guess(filename string, content string) []LicenseGuess {
	if IsLicenceFile(filename) {
		// Check if the filename matches on of the common licences in which case return that
		// since it seems unlikely someone would add a file called LGPL-2.0 without
		// it actually being that licence
		for _, li := range commonLicences {
			if strings.EqualFold(filename, li) {
				return []LicenseGuess{
					{
						Name: li,
					},
				}
			}
		}

		// at this point we are confident we have a licence file, but we don't know which one, so lets
		// start by firstly assuming there is only 1 license in the file
		// and then try to determine what is actually inside the file
		var bestGuess License
		bestMatch := math.MaxInt
		con := []rune(compareOptimize(content))
		for _, li := range l.Database {
			for _, lic := range li.LicenseTexts {
				m := levenshtein.DistanceForStrings([]rune(compareOptimize(lic)), con, levenshtein.DefaultOptions)
				if m < bestMatch {
					bestGuess = li
					bestMatch = m
				}
			}
		}

		if len(bestGuess.LicenseIds) != 0 {
			return []LicenseGuess{
				{
					Name: bestGuess.LicenseIds[0],
				},
			}
		}

		return nil
	}

	if IsReadmeFile(filename) {
		// at this point we are confident we have a licence file, but we don't know which one, so lets
		// start by firstly assuming there is only 1 license in the file
		// and then try to determine what is actually inside the file
		var bestGuess License
		bestMatch := math.MaxInt
		con := []rune(compareOptimize(content))
		for _, li := range l.Database {
			for _, lic := range li.LicenseTexts {
				m := levenshtein.DistanceForStrings([]rune(compareOptimize(lic)), con, levenshtein.DefaultOptions)
				if m < bestMatch {
					bestGuess = li
					bestMatch = m
				}
			}
		}

		if len(bestGuess.LicenseIds) != 0 {
			return []LicenseGuess{
				{
					Name: bestGuess.LicenseIds[0],
				},
			}
		}
	}

	return nil
}
