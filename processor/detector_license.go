package processor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/boyter/lc/processor/levenshtein"
	"regexp"
	"sort"
	"strings"
)

type LicenseData struct {
	LicenseTexts []string `json:"licenseTexts"` // examples of text that we have for these licences
	LicenseIds   []string `json:"licenseIds"`   // SPDX ids where licences are considered identical
	Keywords     []string `json:"keywords"`     // keywords that are unique and can be used to identify this group of licences
}

type LicenceDetector struct {
	UseFullDatabase bool
	LicenseData     []LicenseData
}

var licenceIdentifier = "Valid-License-Identifier:"
var licenceRegex = regexp.MustCompile(`Valid-License-Identifier:\s+(.*)[ |\n|\r\n]*?`)
var commonLicences = []string{"MIT", "Apache-2.0", "GPL-3.0", "AGPL-3.0", "BSD-3-Clause", "GPL-2.0", "BSD-2-Clause", "CC0-1.0", "LGPL-3.0", "LGPL-2.1", "ISC", "0BSD", "LGPL-2.0", "Unlicense", "BSD-3-Clause-No-Nuclear-License-2014", "MPL-2.0", "EPL-1.0", "MPL-2.0-no-copyleft-exception", "AGPL-1.0", "CC-BY-4.0", "IPL-1.0", "CPL-1.0", "CC-BY-3.0", "CC-BY-SA-4.0", "WTFPL", "Zlib", "CC-BY-SA-3.0", "Cube", "JSON", "BitTorrent-1.0"}

func NewLicenceDetector(useFullDatabase bool) *LicenceDetector {
	l := LicenceDetector{
		UseFullDatabase: useFullDatabase,
		LicenseData:     []LicenseData{},
	}

	data, _ := base64.StdEncoding.DecodeString(database_keywords)
	_ = json.Unmarshal(data, &l.LicenseData)

	return &l
}

func (l *LicenceDetector) Detect(filename string, content string) []IdentifiedLicense {
	// Step 1. Check if there is a SPDX identifier, and if that is found assume
	// that it is correct because why else would it be there
	spdxIdentified := l.SpdxDetect(content)
	if len(spdxIdentified) != 0 {
		var licenses []IdentifiedLicense
		for _, s := range spdxIdentified {
			licenses = append(licenses, IdentifiedLicense{
				LicenseId:       s,
				ScorePercentage: 100,
			})
		}
		return licenses
	}

	// Step 2. Check the filename to determine if there is something we can use there
	// If the name matches and the length of the content is close to the real one its probably safe
	// say it's that
	for _, lic := range spdxLicenseIds {
		if lic == filename {

			// if we have a potential match, so now find the licence that matches and check the distance
			ld, ok := l.findLicenseById(filename)
			if ok {

				// now we check the content to see if its a similar size then we vector compare
				// to determine how close it is
				// note that we need to do it for all the possible licence texts as things like
				// MIT have multiple
				var bestScore float64

				con := BuildConcordance(strings.Fields(LcCleanText(content)))
				for _, te := range ld.LicenseTexts {
					con2 := BuildConcordance(strings.Fields(LcCleanText(te)))
					score := Relation(con, con2)
					if bestScore < score {
						bestScore = score
					}
				}

				// TODO move this into something configurable
				// TODO need to move to keyword matching and whatever else we do...
				if bestScore >= 0.9 {
					return []IdentifiedLicense{
						{
							LicenseId:       lic,
							ScorePercentage: bestScore * 100,
						},
					}
				}
			}
		}
	}

	// Step 3. We suspect it is a licence but we don't have a clue which one. Start the 3 step program
	// to determine what it might be starting with
	// a. keywords
	//	if the keywords is not conclusive then fall back to vector space
	// b. vector space

	return nil
}

func (l *LicenceDetector) vectorDetect(content string) []IdentifiedLicense {
	con := BuildConcordance(strings.Fields(LcCleanText(content)))

	var possible []IdentifiedLicense
	for _, ld := range l.LicenseData {
		if !l.UseFullDatabase {
			if !ContainsString(ld.Keywords, commonLicences) {
				continue
			}
		}

		for _, lt := range ld.LicenseTexts {
			con2 := BuildConcordance(strings.Fields(LcCleanText(lt)))
			score := Relation(con, con2)

			for _, li := range ld.LicenseIds {
				possible = append(possible, IdentifiedLicense{
					LicenseId:       li,
					ScorePercentage: score,
				})
			}
		}
	}

	sort.Slice(possible, func(i, j int) bool {
		return possible[i].ScorePercentage >= possible[j].ScorePercentage
	})

	average := 0.0
	for _, il := range possible {
		average += il.ScorePercentage
	}
	average = average / float64(len(possible))

	var bestPossible []IdentifiedLicense
	for _, p := range possible {
		if p.ScorePercentage > average {
			bestPossible = append(bestPossible, p)
		}
	}

	return bestPossible
}

func (l *LicenceDetector) levenshteinDetect(content string) []IdentifiedLicense {
	lev1 := LcCleanText(content)

	var possible []IdentifiedLicense
	for _, ld := range l.LicenseData {
		if !l.UseFullDatabase {
			if !ContainsString(ld.Keywords, commonLicences) {
				continue
			}
		}

		for _, li := range ld.LicenseTexts {
			lev2 := LcCleanText(li)

			possible = append(possible, IdentifiedLicense{
				LicenseId:       li,
				ScorePercentage: float64(levenshtein.DistanceForStrings([]rune(lev1), []rune(lev2), levenshtein.DefaultOptions)),
			})
		}
	}

	sort.Slice(possible, func(i, j int) bool {
		return possible[i].ScorePercentage <= possible[j].ScorePercentage
	})

	// only take the top quartile 

	average := 0.0
	for _, il := range possible {
		average += il.ScorePercentage
	}
	average = average / float64(len(possible))

	var bestPossible []IdentifiedLicense
	for _, p := range possible {
		if p.ScorePercentage < average {
			bestPossible = append(bestPossible, p)
		}
	}

	fmt.Println(bestPossible)

	return nil
}

func (l *LicenceDetector) keywordDetect(content string) []IdentifiedLicense {
	cleaned := LcCleanText(content)
	var possible []IdentifiedLicense

	for _, ld := range l.LicenseData {
		if !l.UseFullDatabase {
			if !ContainsString(ld.Keywords, commonLicences) {
				continue
			}
		}

		count := 0
		for _, v := range ld.Keywords {
			if strings.Contains(cleaned, v) {
				count++
			}
		}

		if count != 0 {
			for _, l := range ld.LicenseIds {
				possible = append(possible, IdentifiedLicense{
					LicenseId:       l,
					ScorePercentage: float64(count),
				})
			}
		}
	}

	// now that we have some, lets find the average and remove anything that's less than the average
	average := 0.0
	for _, il := range possible {
		average += il.ScorePercentage
	}
	average = average / float64(len(possible))

	var bestPossible []IdentifiedLicense
	for _, p := range possible {
		if p.ScorePercentage > average {
			bestPossible = append(bestPossible, p)
		}
	}

	return bestPossible
}

func (l *LicenceDetector) findLicenseById(id string) (LicenseData, bool) {
	for _, ld := range l.LicenseData {
		for _, ln := range ld.LicenseIds {
			if id == ln {
				return ld, true
			}
		}
	}

	return LicenseData{}, false
}

// SpdxDetect will identify licenses in the text which are using the SPDX indicator for licences
// https://www.kernel.org/doc/html/latest/process/license-rules.html
// which is reasonably cheap in terms of looking things up
func (l *LicenceDetector) SpdxDetect(content string) []string {
	// cheap check to see if there might be on in the source code
	if strings.Index(content, licenceIdentifier) == -1 {
		return nil
	}

	var matchingLicenses []string
	matches := licenceRegex.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		var toCheck []string
		t := strings.TrimSpace(val[1])
		if strings.Contains(val[1], " ") {
			// deal with multiple with an OR or some such
			for _, x := range strings.Split(t, " ") {
				x = strings.TrimSpace(x)
				if x != "" {
					toCheck = append(toCheck, x)
				}
			}
		} else {
			toCheck = append(toCheck, t)
		}

		for _, x := range toCheck {
			found := false
			// Check the full database because there is so little cost to do so
			for _, license := range spdxLicenseIds {
				if license == x {
					matchingLicenses = append(matchingLicenses, license)
					found = true
					// we should only ever find a single per what we are checking
					break
				}
			}

			// if we didn't find anything try using lower case because hey why not
			if !found {
				x = strings.ToLower(x)
				for _, license := range spdxLicenseIds {
					if strings.ToLower(license) == x {
						matchingLicenses = append(matchingLicenses, license)
						// we should only ever find a single per what we are checking
						break
					}
				}
			}
		}
	}

	// filter out duplicates because its possible someone put in multiple markers of the same
	var found = map[string]bool{}
	var filtered []string

	for _, lic := range matchingLicenses {
		b := found[lic]
		if !b {
			filtered = append(filtered, lic)
			found[lic] = true
		}
	}

	return filtered
}
