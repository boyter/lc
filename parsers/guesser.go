package parsers

import (
	vectorspace "github.com/boyter/golangvectorspace"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

var spdxLicenceRegex = regexp.MustCompile(`SPDX-License-Identifier:\s+(.*)[ |\n|\r\n]*?`)
var alphaNumericRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
var multipleSpacesRegex = regexp.MustCompile("\\s+")

func cleanText(content string) string {
	content = strings.ToLower(content)

	content = alphaNumericRegex.ReplaceAllString(content, " ")
	content = multipleSpacesRegex.ReplaceAllString(content, " ")

	return content
}

// Identify licenses in the text which is using the SPDX indicator
// Can return multiple license matches
func identifierGuessLicence(content string, licenses []License) []LicenseMatch {
	matchingLicenses := []LicenseMatch{}
	matches := spdxLicenceRegex.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		for _, license := range licenses {
			if license.LicenseId == strings.TrimSpace(val[1]) {
				matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: 1})
			}
		}
	}

	return matchingLicenses
}

// Fast method of checking if supplied content contains a licence using
// matching keyword ngrams to find if the licence is a match or not
// returns the maching licences with shortname and the percentage of match.
func keywordGuessLicense(content string, licenses []License) []LicenseMatch {
	content = cleanText(content)

	matchingLicenses := []LicenseMatch{}

	for _, license := range licenses {
		keywordmatch := 0
		contains := false

		for _, keyword := range license.Keywords {
			contains = strings.Contains(content, strings.ToLower(keyword))

			if contains == true {
				keywordmatch++
			}
		}

		if keywordmatch > 0 {
			percentage := (float64(keywordmatch) / float64(len(license.Keywords))) * 100
			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: percentage})
		}
	}

	return matchingLicenses
}

// Parses the supplied file content against the list of licences and
// returns the matching licences with the shortname and the percentage of match.
// If fast lookup methods fail it will try deep matching which is slower.
func guessLicense(content string, deepguess bool, licenses []License) []LicenseMatch {
	matchingLicenses := []LicenseMatch{}

	for _, license := range keywordGuessLicense(content, licenses) {
		matchingLicense := License{}

		for _, l := range licenses {
			if l.LicenseId == license.LicenseId {
				matchingLicense = l
				break
			}
		}

		runecontent := []rune(content)
		trimto := utf8.RuneCountInString(matchingLicense.LicenseText)

		if trimto > len(runecontent) {
			trimto = len(runecontent)
		}

		contentConcordance := vectorspace.BuildConcordance(string(runecontent[:trimto]))
		relation := vectorspace.Relation(matchingLicense.Concordance, contentConcordance)

		// Average out the vector calc against the keyword percentage
		relation = (relation + math.Min(1, (license.Percentage/100)+0.5)) / 2

		if relation >= confidence {
			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: relation})
		}
	}

	if len(matchingLicenses) == 0 && deepguess == true {
		for _, license := range licenses {
			runecontent := []rune(content)
			trimto := utf8.RuneCountInString(license.LicenseText)

			if trimto > len(runecontent) {
				trimto = len(runecontent)
			}

			contentConcordance := vectorspace.BuildConcordance(string(runecontent[:trimto]))
			relation := vectorspace.Relation(license.Concordance, contentConcordance)

			if relation >= confidence {
				matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: relation})
			}
		}
	}

	sort.Slice(matchingLicenses, func(i, j int) bool {
		return matchingLicenses[i].Percentage > matchingLicenses[j].Percentage
	})

	// Special cases such as MIT and JSON here
	if len(matchingLicenses) > 2 && ((matchingLicenses[0].LicenseId == "JSON" && matchingLicenses[1].LicenseId == "MIT") ||
		(matchingLicenses[0].LicenseId == "MIT" && matchingLicenses[1].LicenseId == "JSON")) {
		if strings.Contains(strings.ToLower(content), "not evil") {
			// Its JSON
			matchingLicenses = []LicenseMatch{}
			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: "JSON", Percentage: 1})
		} else {
			// Its MIT
			matchingLicenses = []LicenseMatch{}
			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: "MIT", Percentage: 1})
		}
	}

	return matchingLicenses
}

func processFile(directory string, file string, rootLicenses []LicenseMatch) FileResult {
	process := true

	for _, ext := range strings.Split(ExtentionBlacklist, ",") {
		if strings.HasSuffix(file, "."+ext) {
			// Needs to be smarter we should skip reading the contents but it should still be under the license in the root folders
			process = false
		}
	}

	content, _ := ioutil.ReadFile(filepath.Join(directory, file))
	licenseGuesses := []LicenseMatch{}
	licenseIdentified := []LicenseMatch{}

	if len(content) > maxSize {
		process = false
	}

	if process == true {
		licenseGuesses = guessLicense(string(content), deepGuess, loadDatabase())
		licenseIdentified = identifierGuessLicence(string(content), loadDatabase())
	}

	fileResult := FileResult{
		Directory:         directory,
		Filename:          file,
		LicenseGuesses:    licenseGuesses,
		LicenseRoots:      rootLicenses,
		LicenseIdentified: licenseIdentified,
		Md5Hash:           getMd5Hash(content),
		Sha1Hash:          getSha1Hash(content),
		Sha256Hash:        getSha256Hash(content),
		BytesHuman:        bytesToHuman(int64(len(content))),
		Bytes:             len(content)}

	return fileResult
}
