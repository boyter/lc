package processor

import (
	"encoding/base64"
	"encoding/json"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"regexp"
	"sort"
	"strings"
)

var spdxLicenceIdentifier = "SPDX-License-Identifier:"
var spdxLicenceRegex = regexp.MustCompile(`SPDX-License-Identifier:\s+(.*)[ |\n|\r\n]*?`)

func NewLicenceGuesser() LicenceGuesser {
	l := LicenceGuesser{}
	l.LoadDatabase()
	l.UseFullDatabase = false
	return l
}

type LicenceGuesser struct {
	Database []License
	CommonDatabase []License
	UseFullDatabase bool
}

// LoadDatabase will initialize the database values and should only be called once such as in an init
func (l *LicenceGuesser) LoadDatabase() {
	if len(l.Database) != 0 {
		return
	}

	data, _ := base64.StdEncoding.DecodeString(database_keywords)
	_ = json.Unmarshal(data, &l.Database)

	// Load smaller faster database for checking most common licenses
	common := []string{
		"MIT",
		"Apache-2.0",
		"GPL-3.0",
		"AGPL-3.0",
		"BSD-3-Clause",
		"GPL-2.0",
		"BSD-2-Clause",
		"CC0-1.0",
		"LGPL-3.0",
		"LGPL-2.1",
		"ISC",
		"0BSD",
		"LGPL-2.0",
		"Unlicense",
		"BSD-3-Clause-No-Nuclear-License-2014",
		"MPL-2.0",
		"EPL-1.0",
		"MPL-2.0-no-copyleft-exception",
		"AGPL-1.0",
		"CC-BY-4.0",
		"IPL-1.0",
		"CPL-1.0",
		"CC-BY-3.0",
		"CC-BY-SA-4.0",
		"WTFPL",
		"Zlib",
		"CC-BY-SA-3.0",
		"Cube",
		"JSON",
		"BitTorrent-1.0",
	}

	for _, license := range l.Database {
		for _, com := range common {
			if license.LicenseId == com {
				l.CommonDatabase = append(l.CommonDatabase, license)
			}
		}
	}
}

// Identify licenses in the text which are using the SPDX indicator
// which is reasonably cheap in terms of looking things up
// This is the only guesser that is 100% accurate as literally everything
// else is slightly fuzzy and "best" effort
func (l *LicenceGuesser) SpdxIdentify(content string) []License {
	if strings.Index(content, spdxLicenceIdentifier) == -1 {
		return nil
	}

	var matchingLicenses []License
	matches := spdxLicenceRegex.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		for _, license := range l.Database {
			if license.LicenseId == strings.TrimSpace(val[1]) {
				license.ScorePercentage = 100 // set the score to be 100% IE we are 100% confidence in this guess
				matchingLicenses = append(matchingLicenses, license)
			}

			// TODO should do a lowercase check here with less confidence
		}
	}

	return matchingLicenses
}

// Try to guess the licence for this content based on whatever heuristics are set on the
// guesser itself
// TODO this should guess based on criteria IE try spdx first then keywords etc...
func (l *LicenceGuesser) GuessLicense(content []byte) []License {
	// try via spdx
	// try via keywords
	// try via vector space or otherwise
	return l.KeyWordGuessLicence(content)
	// so if the keyword guess licence found nothing we should consider using vector space
}

// Given some content try to guess what the licence is based on checking for unique keywords
// using the prebuilt licence library
func (l *LicenceGuesser) KeyWordGuessLicence(content []byte) []License {
	haystack := LcCleanText(string(content))

	var matchingLicenses []License
	var totalScore float64
	var maxScore float64

	// Swap out the database to the full one if configured to use it
	db := l.CommonDatabase
	if l.UseFullDatabase {
		db = l.Database
	}

	for _, lic := range db {
		c := l.checkLicenceKeywords(haystack, lic)
		if c != 0 {
			lic.ScorePercentage = float64(c)
			totalScore += lic.ScorePercentage
			matchingLicenses = append(matchingLicenses, lic)
		}
	}

	// Normalise the scores based on the total so we can make a reasonable guess of how confident we are
	for i := 0; i < len(matchingLicenses); i++ {
		matchingLicenses[i].ScorePercentage = (matchingLicenses[i].ScorePercentage / totalScore) * 100

		// keep track of the highest score
		if matchingLicenses[i].ScorePercentage > maxScore {
			maxScore = matchingLicenses[i].ScorePercentage
		}
	}

	// If we have multiple licenses and their scores aren't at least 70% confident do some additional checks
	if len(matchingLicenses) >= 2 && maxScore < 60 {
		for i := 0; i< len(matchingLicenses); i++ {
			distance := levenshtein.DistanceForStrings([]rune(haystack), []rune(LcCleanText(matchingLicenses[i].LicenseText)), levenshtein.DefaultOptions)
			matchingLicenses[i].ScorePercentage = float64(100) / float64(distance)
		}
	}

	// Sort so if someone wants to get the best candidate they can get the 0 index of the return
	sort.Slice(matchingLicenses, func(i, j int) bool {
		return matchingLicenses[i].ScorePercentage > matchingLicenses[j].ScorePercentage
	})

	return matchingLicenses
}

func (l *LicenceGuesser) checkLicenceKeywords(haystack string, lic License) int {
	var count int
	for _, k := range lic.Keywords {
		if strings.Contains(haystack, k) {
			count++
		}
	}

	return count
}