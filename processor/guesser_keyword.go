// SPDX-License-Identifier: MIT
// SPDX-License-Identifier: Unlicense

package processor

import (
	"sort"
)

// Given some content try to guess what the licence is based on checking for unique keywords
// using the prebuilt licence library which contains what we hope are unique ngrams for each
// licence
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
		match := lic.Trie.Match([]byte(haystack))

		if len(match) != 0 {
			lic.ScorePercentage = float64(len(match))
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

	// only keep those close to the max score so we can ignore anything that isn't even close
	var t []License
	for _, lic := range matchingLicenses {
		if lic.ScorePercentage >= (maxScore * 0.8) {
			t = append(t, lic)
		}
	}
	if len(t) != 0 {
		matchingLicenses = t
	}

	// TODO this should be moved out
	// this appears to be horribly slow...
	//// If we have multiple licenses and their scores aren't at least 70% confident do some additional checks
	//if len(matchingLicenses) >= 2 && maxScore < 60 {
	//	for i := 0; i< len(matchingLicenses); i++ {
	//		distance := leven.Distance(haystack, LcCleanText(matchingLicenses[i].LicenseText))
	//		if distance == 0 {
	//			matchingLicenses[i].ScorePercentage = 100
	//		} else {
	//			matchingLicenses[i].ScorePercentage = float64(100) / float64(distance)
	//		}
	//	}
	//}

	// Sort so if someone wants to get the best candidate they can get the 0 index of the return
	sort.Slice(matchingLicenses, func(i, j int) bool {
		return matchingLicenses[i].ScorePercentage > matchingLicenses[j].ScorePercentage
	})

	return matchingLicenses
}
