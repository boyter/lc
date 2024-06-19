// SPDX-License-Identifier: AGPL-3.0

package pkg

import (
	"fmt"
	"regexp"
	"strings"
)

var commonLicences = []string{"MIT", "Apache-2.0", "GPL-3.0", "AGPL-3.0", "BSD-3-Clause", "GPL-2.0", "BSD-2-Clause", "CC0-1.0", "LGPL-3.0", "LGPL-2.1", "ISC", "0BSD", "LGPL-2.0", "Unlicense", "BSD-3-Clause-No-Nuclear-License-2014", "MPL-2.0", "EPL-1.0", "MPL-2.0-no-copyleft-exception", "AGPL-1.0", "CC-BY-4.0", "IPL-1.0", "CPL-1.0", "CC-BY-3.0", "CC-BY-SA-4.0", "WTFPL", "Zlib", "CC-BY-SA-3.0", "Cube", "JSON", "BitTorrent-1.0"}

// Lifted from https://github.com/go-enry/go-license-detector/blob/580c5627556917dee649cdb2b179cb42d6c56a60/licensedb/internal/investigation.go#L29
// SPDX-License-Identifier: Apache-2.0
var (
	// Base names of guessable license files
	licenseFileNames = []string{
		"li[cs]en[cs]e(s?)",
		"legal",
		"copy(left|right|ing)",
		"unlicense",
		"l?gpl([-_ v]?)(\\d\\.?\\d)?",
		"bsd",
		"mit",
		"apache",
	}

	// License file extensions. Combined with the fileNames slice
	// to create a set of files we can reasonably assume contain
	// licensing information.
	fileExtensions = []string{
		"",
		".md",
		".rst",
		".html",
		".txt",
	}

	licenseFileRe = regexp.MustCompile(
		fmt.Sprintf("^(|.*[-_. ])(%s)(|[-_. ].*)$",
			strings.Join(licenseFileNames, "|")))

	readmeFileRe = regexp.MustCompile(fmt.Sprintf("^(readme|guidelines)(%s)$",
		strings.Replace(strings.Join(fileExtensions, "|"), ".", "\\.", -1)))
)

func IsLicenceFile(filename string) bool {
	// attempt to filter out false positives that come from java due to filenames
	if strings.Count(filename, ".") > 2 {
		return false
	}

	return licenseFileRe.Match([]byte(strings.ToLower(filename)))
}

func IsReadmeFile(filename string) bool {
	return readmeFileRe.Match([]byte(strings.ToLower(filename)))
}

func compareOptimize(input string) string {
	tokens := strings.Fields(input)
	var sb strings.Builder
	skipTokens := map[string]struct{}{}
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		foundLonger := false

		// if we have already looked at this token, skip it, important for performance
		_, ok := skipTokens[tok]
		if ok {
			continue
		}
		skipTokens[tok] = struct{}{}

		for j := i; j < len(tokens); j++ {
			tok2 := tokens[j]
			if tok == tok2 {
				continue
			}

			if len(tok2) <= len(tok) {
				continue
			}

			if strings.Contains(tok2, tok) {
				foundLonger = true
			}
		}

		if !foundLonger {
			sb.WriteString(tok)
		}
	}

	return sb.String()
}
