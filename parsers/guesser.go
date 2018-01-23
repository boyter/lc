package parsers

import (
	vectorspace "github.com/boyter/golangvectorspace"
	"sort"
	"strings"
	"unicode/utf8"
)

var confidence = 0.85

type License struct {
	Keywords    []string `json:"keywords"`
	Text        string   `json:"text"`
	Fullname    string   `json:"fullname"`
	Shortname   string   `json:"shortname"`
	Header      string   `json:"header"`
	Concordance vectorspace.Concordance
}

type LicenseMatch struct {
	Shortname  string
	Percentage float64
}

// Fast method of checking if supplied content contains a licence using
// matching keyword ngrams to find if the licence is a match or not
// returns the maching licences with shortname and the percentage of match.
func KeywordGuessLicense(content string, licenses []License) []LicenseMatch {
	content = strings.ToLower(content)
	matchingLicenses := []LicenseMatch{}

	for _, license := range licenses {
		keywordmatch := 0
		contains := false

		for _, keyword := range license.Keywords {
			contains = strings.Contains(content, keyword)
			if contains {
				keywordmatch++
			}
		}

		if keywordmatch > 0 {
			percentage := (float64(keywordmatch) / float64(len(license.Keywords))) * 100
			matchingLicenses = append(matchingLicenses, LicenseMatch{Shortname: license.Shortname, Percentage: percentage})
		}
	}

	return matchingLicenses
}

// Parses the supplied file content against the list of licences and
// returns the matching licences with the shortname and the percentage of match.
// If fast lookup methods fail it will try deep matching which is slower.
func GuessLicense(content string, licenses []License) []LicenseMatch {
	matchingLicenses := []LicenseMatch{}

	for _, license := range KeywordGuessLicense(content, licenses) {

		matchingLicense := License{}

		for _, l := range licenses {
			if l.Shortname == license.Shortname {
				matchingLicense = l
				break
			}
		}

		runecontent := []rune(content)
		trimto := utf8.RuneCountInString(matchingLicense.Text)

		if trimto > len(runecontent) {
			trimto = len(runecontent)
		}

		contentConcordance := vectorspace.BuildConcordance(string(runecontent[:trimto]))
		relation := vectorspace.Relation(matchingLicense.Concordance, contentConcordance)

		if relation >= confidence {
			matchingLicenses = append(matchingLicenses, LicenseMatch{Shortname: license.Shortname, Percentage: relation})
		}
	}

	if len(matchingLicenses) == 0 {
		for _, license := range licenses {
			runecontent := []rune(content)
			trimto := utf8.RuneCountInString(license.Text)

			if trimto > len(runecontent) {
				trimto = len(runecontent)
			}

			contentConcordance := vectorspace.BuildConcordance(string(runecontent[:trimto]))
			relation := vectorspace.Relation(license.Concordance, contentConcordance)

			if relation >= confidence {
				matchingLicenses = append(matchingLicenses, LicenseMatch{Shortname: license.Shortname, Percentage: relation})
			}
		}
	}

	sort.Slice(matchingLicenses, func(i, j int) bool {
		return matchingLicenses[i].Percentage > matchingLicenses[j].Percentage
	})

	return matchingLicenses
}
