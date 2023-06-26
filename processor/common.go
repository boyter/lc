// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"fmt"
	"regexp"
	"strings"
)

var alphaNumericRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
var multipleSpacesRegex = regexp.MustCompile("\\s+")

// Very specific cleaner which is designed to clean to the format lc uses to match things
// so be very careful if you choose to use this
func LcCleanText(content string) string {
	content = strings.ToLower(content)

	content = alphaNumericRegex.ReplaceAllString(content, " ")
	content = multipleSpacesRegex.ReplaceAllString(content, " ")
	content = strings.TrimSpace(content)

	return content
}

const (
	MatchTypeSpdx    = "SPDX"
	MatchTypeKeyword = "Keyword"
	MatchTypeVector  = "Vector"
	MatchTypeBlended = "Blended"
)

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

func ContainsString(ids []string, lst []string) bool {
	for _, j := range ids {
		for _, i := range lst {
			if i == j {
				return true
			}
		}
	}

	return false
}
