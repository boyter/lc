// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"encoding/base64"
	"encoding/json"
	corasick "github.com/BobuSumisu/aho-corasick"
	"strings"
)

func NewLicenceGuesser(keyword bool, vectorspace bool) LicenceGuesser {
	l := LicenceGuesser{}
	l.keyword = keyword
	l.vectorspace = vectorspace
	l.LoadDatabase()
	l.UseFullDatabase = false
	l.cutoffPercentage = 60
	return l
}

type LicenceGuesser struct {
	Database         []License
	CommonDatabase   []License
	UseFullDatabase  bool
	keyword          bool
	vectorspace      bool
	cutoffPercentage float64
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

	if l.keyword {
		for i := 0; i < len(l.Database); i++ {
			l.Database[i].Trie = corasick.NewTrieBuilder().
				AddStrings(l.Database[i].Keywords).
				Build()
			// Precompute word sets for Jaccard similarity
			if len(l.Database[i].LicenseTexts) > 0 {
				words := strings.Fields(LcCleanText(l.Database[i].LicenseTexts[0]))
				ws := make(map[string]struct{}, len(words))
				for _, w := range words {
					ws[w] = struct{}{}
				}
				l.Database[i].WordSet = ws
			}
		}
	}

	if l.vectorspace {
		for i := 0; i < len(l.Database); i++ {
			if len(l.Database[i].LicenseTexts) > 0 {
				l.Database[i].Concordance = BuildConcordance(strings.Split(LcCleanText(l.Database[i].LicenseTexts[0]), " "))
			}
		}
	}

	for _, license := range l.Database {
		for _, com := range common {
			for _, lid := range license.LicenseIds {
				if lid == com {
					l.CommonDatabase = append(l.CommonDatabase, license)
				}
			}
		}
	}
}
