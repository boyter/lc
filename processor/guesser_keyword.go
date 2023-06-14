// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"sort"
)

// KeyWordGuessLicence will some content try to guess what the licence is based on checking for unique keywords
// using the prebuilt licence library which contains what we hope are unique ngrams for each licence
func (l *LicenceGuesser) KeyWordGuessLicence(content []byte) []License {
	haystack := LcCleanText(string(content))

	var matchingLicenses []License
	var totalScore float64
	var maxScore float64
	var highestMatch int

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
			lic.MatchType = MatchTypeKeyword
			matchingLicenses = append(matchingLicenses, lic)

			if len(match) > highestMatch {
				highestMatch = len(match)
			}
		}
	}

	// licences come in with a count of terms found like the below
	// 1 1 1 1 1 13 2 1 1 1
	// what the above is saying is that one option is 13 times more likely to be correct than most of the others
	// because it has a lot more matches...
	// so lets scale in accordance to this, were we say that

	for i := 0; i < len(matchingLicenses); i++ {
		matchingLicenses[i].ScorePercentage = (matchingLicenses[i].ScorePercentage / float64(highestMatch)) * 100

		// keep track of the highest score
		if matchingLicenses[i].ScorePercentage > maxScore {
			maxScore = matchingLicenses[i].ScorePercentage
		}
	}

	// Normalise the scores based on the total so we can make a reasonable guess of how confident we are
	// TODO the problem with this is that when we have more false matches we weigh down the highest match...
	//for i := 0; i < len(matchingLicenses); i++ {
	//	matchingLicenses[i].ScorePercentage = (matchingLicenses[i].ScorePercentage / totalScore) * 100
	//
	//	// keep track of the highest score
	//	if matchingLicenses[i].ScorePercentage > maxScore {
	//		maxScore = matchingLicenses[i].ScorePercentage
	//	}
	//}

	// TODO the problem with this is that its limited based on precanned values which might not be true
	// if we have more than 5 matches consider it a 100% match
	//for i := 0; i < len(matchingLicenses); i++ {
	//	if matchingLicenses[i].ScorePercentage >= 10 {
	//		matchingLicenses[i].ScorePercentage = 100
	//	} else if matchingLicenses[i].ScorePercentage >= 5 {
	//		matchingLicenses[i].ScorePercentage = 80
	//	} else if matchingLicenses[i].ScorePercentage >= 4 {
	//		matchingLicenses[i].ScorePercentage = 70
	//	} else if matchingLicenses[i].ScorePercentage >= 3 {
	//		matchingLicenses[i].ScorePercentage = 50
	//	} else if matchingLicenses[i].ScorePercentage >= 2 {
	//		matchingLicenses[i].ScorePercentage = 30
	//	}
	//
	//	if matchingLicenses[i].ScorePercentage > maxScore {
	//		maxScore = matchingLicenses[i].ScorePercentage
	//	}
	//}

	// only keep those close to the max score so we can ignore anything that isn't even close
	var t []License
	for _, lic := range matchingLicenses {
		if lic.ScorePercentage >= (maxScore * 0.6) {
			t = append(t, lic)
		}
	}
	if len(t) != 0 {
		matchingLicenses = t
	}

	// if we have multiple licences compare them directly to determine which one it actually is
	if len(matchingLicenses) >= 2 {
		for i := 0; i < len(matchingLicenses); i++ {
			distance := levenshtein.DistanceForStrings([]rune(string(content)), []rune(matchingLicenses[i].LicenseText), levenshtein.DefaultOptions)
			if distance == 0 {
				matchingLicenses[i].ScorePercentage = 100
			} else {
				matchingLicenses[i].ScorePercentage = float64(100) / float64(distance)
			}
		}
	}

	// Sort so if someone wants to get the best candidate they can get the 0 index of the return
	sort.Slice(matchingLicenses, func(i, j int) bool {
		return matchingLicenses[i].ScorePercentage > matchingLicenses[j].ScorePercentage
	})

	return matchingLicenses
}
