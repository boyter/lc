// SPDX-License-Identifier: MIT OR Unlicense

package processor

import ahocorasick "github.com/BobuSumisu/aho-corasick"

// Represents a license inside the JSON which allows us to hopefully match against
type License struct {
	LicenseText             string            `json:"licenseText"`
	StandardLicenseTemplate string            `json:"standardLicenseTemplate"`
	Name                    string            `json:"name"`
	LicenseId               string            `json:"licenseId"`
	Keywords                []string          `json:"keywords"`
	LicenseIdDuplicates     []string          `json:"duplicates"`
	ExtraLicenseText        []string          `json:"extraLicenseText"`
	ScorePercentage         float64           `json:"scorePercentage"` // this is used so we don't have a new struct
	Trie                    *ahocorasick.Trie // used for faster matching
	Concordance             Concordance       // used for vector matching
	MatchType               string            // indicates how this license match was made
}

// IdentifiedLicense
type IdentifiedLicense struct {
	LicenseId       string  `json:"licenseId"`
	ScorePercentage float64 `json:"scorePercentage"` // this is used so we don't have a new struct
}
