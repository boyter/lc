// SPDX-License-Identifier: MIT OR Unlicense

package processor

import ahocorasick "github.com/BobuSumisu/aho-corasick"

// Represents a license inside the JSON which allows us to hopefully match against
type License struct {
	LicenseTexts    []string          `json:"licenseTexts"`
	LicenseIds      []string          `json:"licenseIds"`
	Keywords        []string          `json:"keywords"`
	ScorePercentage float64           // this is used so we don't have a new struct
	Trie            *ahocorasick.Trie // used for faster matching
	Concordance     Concordance       // used for vector matching
	WordSet         map[string]struct{} // precomputed word set for Jaccard similarity
	MatchType       string            // indicates how this license match was made
}

// IdentifiedLicense
type IdentifiedLicense struct {
	LicenseId       string  `json:"licenseId"`
	ScorePercentage float64 `json:"scorePercentage"` // this is used so we don't have a new struct
}
