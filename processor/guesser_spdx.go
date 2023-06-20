// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"fmt"
	"strings"
)

// SpdxIdentify will identify licenses in the text which are using the SPDX indicator
// which is reasonably cheap in terms of looking things up
// This is the only guesser that is 100% accurate as literally everything
// else is slightly fuzzy and "best" effort
func (l *LicenceGuesser) SpdxIdentify(content string) []License {
	// cheap check to see if there might be on in the source code
	if strings.Index(content, spdxLicenceIdentifier) == -1 {
		return nil
	}

	var matchingLicenses []License
	matches := spdxLicenceRegex.FindAllStringSubmatch(content, -1)

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
			for _, license := range l.Database {
				if license.LicenseId == x {
					license.ScorePercentage = 100 // set the score to be 100% IE we are 100% confidence in this guess
					matchingLicenses = append(matchingLicenses, license)
					found = true
				}

				for _, li := range license.LicenseIdDuplicates {
					if license.LicenseId == li {
						license.ScorePercentage = 100 // set the score to be 100% IE we are 100% confidence in this guess
						matchingLicenses = append(matchingLicenses, license)
						found = true
					}
				}

				//GPL-2.0
				for _, li := range license.LicenseIdDuplicates {
					if li == "GPL-2.0" {
						fmt.Println("HERE YA MORNNG")
					}
				}

			}

			// if we didn't find anything try using lower case because hey why not
			// TODO this could probably fit into the above loop for some free performance
			if !found {
				x = strings.ToLower(x)
				for _, license := range l.Database {
					if strings.ToLower(license.LicenseId) == x {
						license.ScorePercentage = 99.99 // set the score to be 99.99% because we are still very confident
						matchingLicenses = append(matchingLicenses, license)
					}

					for _, li := range license.LicenseIdDuplicates {
						if strings.ToLower(li) == x {
							license.ScorePercentage = 99.99 // set the score to be 100% IE we are 100% confidence in this guess
							matchingLicenses = append(matchingLicenses, license)
							found = true
						}
					}
				}
			}
		}
	}

	// filter out duplicates because its possible, but we shouldn't report it
	var found = map[string]bool{}
	var filtered []License

	for _, lic := range matchingLicenses {
		b := found[lic.LicenseId]
		if !b {
			lic.MatchType = MatchTypeSpdx
			filtered = append(filtered, lic)
			found[lic.LicenseId] = true
		}
	}

	return filtered
}
