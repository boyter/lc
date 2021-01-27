// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"sort"
	"strings"
)

func (l *LicenceGuesser) VectorSpaceGuessLicence(content []byte) []License {
	var matchingLicenses []License

	con := BuildConcordance(strings.Split(LcCleanText(string(content)), " "))

	// Swap out the database to the full one if configured to use it
	db := l.CommonDatabase
	if l.UseFullDatabase {
		db = l.Database
	}

	for _, lic := range db {
		distance := Relation(con, lic.Concordance)
		lic.ScorePercentage = distance * 100
		lic.MatchType = MatchTypeVector
		matchingLicenses = append(matchingLicenses, lic)
	}

	sort.Slice(matchingLicenses, func(i, j int) bool {
		return matchingLicenses[i].ScorePercentage > matchingLicenses[j].ScorePercentage
	})

	return matchingLicenses
}
