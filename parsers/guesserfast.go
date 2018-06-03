package parsers

import (
	"bytes"
	"sort"
	"strings"
	"sync"
)

// Fast method of checking if supplied content contains a licence using
// matching keyword ngrams to find if the licence is a match or not
// returns the maching licences with shortname and the percentage of match.
func keywordGuessLicense(content []byte, licenses []License) []LicenseMatch {
	content = cleanText(content)
	length := len(content)
	lengthFuzzy := length / 100 * 30

	var wg sync.WaitGroup
	output := make(chan LicenseMatch, len(licenses))

	for _, license := range licenses {
		if len(license.LicenseText) >= (length-lengthFuzzy) && len(license.LicenseText) <= (length+lengthFuzzy) {
			wg.Add(1)
			go func(license License) {
				keywordMatch := 0

				for _, keyword := range license.Keywords {
					if bytes.Contains(content, []byte(keyword)) {
						keywordMatch++
					}
				}

				if keywordMatch > 50 { // on the basis of there being 100
					output <- LicenseMatch{LicenseId: license.LicenseId, Percentage: float64(keywordMatch)}
				}
				wg.Done()
			}(license)
		}
	}

	wg.Wait()
	close(output)

	var matchingLicenses []LicenseMatch
	for license := range output {
		matchingLicenses = append(matchingLicenses, license)
	}

	sort.Slice(matchingLicenses, func(i, j int) bool {
		return matchingLicenses[i].Percentage > matchingLicenses[j].Percentage
	})

	matchingLicenses = specialCases(content, matchingLicenses)

	return matchingLicenses
}

func cleanText(content []byte) []byte {
	content = bytes.ToLower(content)

	tmp := alphaNumericRegex.ReplaceAllString(string(content), " ")
	tmp = multipleSpacesRegex.ReplaceAllString(tmp, " ")
	tmp = strings.TrimSpace(tmp)

	return []byte(tmp)
}

func specialCases(content []byte, matchingLicenses []LicenseMatch) []LicenseMatch {
	// Quite often JSON and MIT are confused
	if len(matchingLicenses) > 2 && ((matchingLicenses[0].LicenseId == "JSON" && matchingLicenses[1].LicenseId == "MIT") ||
		(matchingLicenses[0].LicenseId == "MIT" && matchingLicenses[1].LicenseId == "JSON")) {
		if bytes.Contains(content, []byte("not evil")) {
			matchingLicenses = []LicenseMatch{{LicenseId: "JSON", Percentage: 1}}
		} else {
			matchingLicenses = []LicenseMatch{{LicenseId: "MIT", Percentage: 1}}
		}
	}

	// Another one is MIT-feh and MIT
	if len(matchingLicenses) > 2 && matchingLicenses[0].LicenseId == "MIT" {

		if bytes.HasPrefix(content, []byte("mit license")) || bytes.HasPrefix(content, []byte("the mit license")) {
			matchingLicenses = []LicenseMatch{{LicenseId: "MIT", Percentage: 100}}
		} else {
			matchingLicenses = []LicenseMatch{{LicenseId: "MIT-feh", Percentage: 100}}
		}
	}

	return matchingLicenses
}
