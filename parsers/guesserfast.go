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
func keywordGuessLicenseFast(content []byte, licenses []License) []LicenseMatch {
	content = cleanTextFast(content)

	var wg sync.WaitGroup
	output := make(chan LicenseMatch, 10)

	for _, license := range licenses {
		wg.Add(1)
		go func(license License) {
			keywordMatch := 0

			for _, keyword := range license.Keywords {
				if bytes.Contains(content, []byte(strings.ToLower(keyword))) {
					keywordMatch++
				}
			}

			if keywordMatch > 0 {
				percentage := (float64(keywordMatch) / float64(len(license.Keywords))) * 100
				if percentage > 70 {
					output <- LicenseMatch{LicenseId: license.LicenseId, Percentage: percentage}
				}
			}
			wg.Done()
		}(license)
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

	if len(matchingLicenses) > 2 && ((matchingLicenses[0].LicenseId == "JSON" && matchingLicenses[1].LicenseId == "MIT") ||
		(matchingLicenses[0].LicenseId == "MIT" && matchingLicenses[1].LicenseId == "JSON")) {
		if bytes.Contains(content, []byte("not evil")) {
			matchingLicenses = []LicenseMatch{{LicenseId: "JSON", Percentage: 1}}
		} else {
			matchingLicenses = []LicenseMatch{{LicenseId: "MIT", Percentage: 1}}
		}
	}

	return matchingLicenses
}

func cleanTextFast(content []byte) []byte {
	content = bytes.ToLower(content)

	tmp := alphaNumericRegex.ReplaceAllString(string(content), " ")
	tmp = multipleSpacesRegex.ReplaceAllString(tmp, " ")

	return []byte(tmp)
}
