// SPDX-License-Identifier: MIT
// SPDX-License-Identifier: Unlicense

package processor

import ahocorasick "github.com/BobuSumisu/aho-corasick"

// Represents a license inside the JSON which allows us to hopefully match against
type License struct {
	LicenseText             string   `json:"licenseText"`
	StandardLicenseTemplate string   `json:"standardLicenseTemplate"`
	Name                    string   `json:"name"`
	LicenseId               string   `json:"licenseId"`
	Keywords                []string `json:"keywords"`
	ScorePercentage         float64  `json:"scorePercentage"` // this is used so we don't have a new struct
	Trie					*ahocorasick.Trie // used for faster matching
	Concordance Concordance // used for vector matching
}