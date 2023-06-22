package processor

import (
	"regexp"
	"strings"
)

func NewLicenceDetector(useFullDatabase bool) LicenceDetector {
	l := LicenceDetector{
		UseFullDatabase: useFullDatabase,
	}
	return l
}

type LicenceDetector struct {
	UseFullDatabase bool
}

var licenceIdentifier = "Valid-License-Identifier:"
var licenceRegex = regexp.MustCompile(`Valid-License-Identifier:\s+(.*)[ |\n|\r\n]*?`)

func (l *LicenceDetector) Detect(filename string, content string) []IdentifiedLicense {
	// Step 1. Check if there is a SPDX identifier, and if that is found assume
	// that it is correct because why else would it be there
	spdxIdentified := l.SpdxDetect(content)
	if len(spdxIdentified) != 0 {
		var licenses []IdentifiedLicense
		for _, s := range spdxIdentified {
			licenses = append(licenses, IdentifiedLicense{
				Name:            "", // TODO need lookup to get the name
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

			// now we check the content to see if its a similar size then we vector compare
			// to determine how close it is
			// note that we need to do it for all the possible licence texts as things like
			// MIT have multiple
			return []IdentifiedLicense{
				{
					Name:            "",
					LicenseId:       lic,
					ScorePercentage: 100,
				},
			}
		}
	}

	// Step 3. We suspect it is a licence but we don't have a clue which one. Start the 3 step program
	// to determine what it might be starting with
	// a. keywords
	// b. vector space

	return nil
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
