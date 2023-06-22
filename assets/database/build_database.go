// SPDX-License-Identifier: MIT OR Unlicense

package main

import (
	"encoding/json"
	"fmt"
	"github.com/boyter/lc/processor"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// License contains details from loading off the SPDX list and holds state for changes
// we make before saving to disk so it can be used in lc itself
type License struct {
	LicenseText             string `json:"licenseText"`
	StandardLicenseTemplate string `json:"standardLicenseTemplate"`
	StandardLicenseHeader   string `json:"standardLicenseHeader"`
	Name                    string `json:"name"`
	LicenseId               string `json:"licenseId"`

	// the below are used internally when processing
	Ngrams           []string   // the unique bits of string we used to identify a licence
	LicenceNgrams    [][]string // for each extra text + the main ngrams so we can get a selection from each
	LicenseIds       []string   // indicates that this has 100% the same LicenseText such as AGPL-3.0-only and AGPL-3.0-or-later
	ExtraLicenseText []string   // things like MIT have multiple variants of the same license and we need to track that
}

// LicenseOutput is the output format that we save to disk and import into lc
type LicenseOutput struct {
	LicenseTexts []string `json:"licenseTexts"` // examples of text that we have for these licences
	LicenseIds   []string `json:"licenseIds"`   // SPDX ids where licences are considered identical
	Keywords     []string `json:"keywords"`     // keywords that are unique and can be used to identify this group of licences
}

// returns all the ngrams of a supplied size for supplied list
func findNgrams(list []string, size int) []string {
	var ngrams []string

	for i := 0; i < len(list); i++ {
		if i+size < len(list)+1 {
			ngram := list[i : i+size]
			ngrams = append(ngrams, strings.Join(ngram, " "))
		}
	}

	return ngrams
}

var startNgrams = 3
var endNgrams = 6
var keepNgrams = 200

func main() {
	// find the licence files that we need to compare against as a starting
	// point based on the SPDX
	files, _ := os.ReadDir("./licenses/")

	fmt.Println("loading licenses")

	var allLicenseIds []string
	licenseTextCount := map[string]int{}
	var licenses []License
	// Load all the licenses from disk and keep track of duplicates
	for _, f := range files {
		//if !strings.HasPrefix(f.Name(), "GPL") {
		//	continue
		//}
		//fmt.Println(f.Name())

		if strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := os.ReadFile(filepath.Join("./licenses/", f.Name()))

			var license License
			_ = json.Unmarshal(bytes, &license)
			license.Ngrams = []string{}
			license.ExtraLicenseText = append(license.ExtraLicenseText, license.LicenseText)

			// if MIT add in the other example so we can match it better...
			if license.LicenseId == "MIT" {
				fmt.Println("adding extra to MIT")
				license.ExtraLicenseText = append(license.ExtraLicenseText, mitExtra...)
			}

			if license.LicenseId == "BSD-3-Clause" {
				fmt.Println("adding extra to BSD-3-Clause")
				license.ExtraLicenseText = append(license.ExtraLicenseText, bsd3ClauseExtra...)
			}

			allLicenseIds = append(allLicenseIds, license.LicenseId)
			license.LicenseIds = append(license.LicenseIds, license.LicenseId)
			licenses = append(licenses, license)

			// track where the license text is the licenseTextCount
			licenseTextCount[processor.LcCleanText(license.LicenseText)] = licenseTextCount[processor.LcCleanText(license.LicenseText)] + 1
		}
	}

	fmt.Println("the following are duplicates...")
	for k, v := range licenseTextCount {
		// if we have only 1 there is no point doing anything else
		if v == 1 {
			continue
		}

		// go through each licence we have and if its text matches what we are checking keep track of which one it is
		// such that we can mark them as duplicates
		var d []string
		for _, lic := range licenses {
			if processor.LcCleanText(lic.LicenseText) == k {
				d = append(d, lic.LicenseId)
			}
		}

		// update any license with this text to tell it about all the duplicates
		for i := 0; i < len(licenses); i++ {
			if processor.LcCleanText(licenses[i].LicenseText) == k {
				licenses[i].LicenseIds = d
				fmt.Println(fmt.Sprintf("	%v %v duplicates %v", licenses[i].LicenseId, len(d), d))
			}
		}
	}

	fmt.Println("removing duplicates...")
	// now we have identified the duplicates lets remove them such that we only have a single item
	var deduplicatedLicences []License
	var seenIds []string
	for _, l := range licenses {
		seen := false
		if len(l.LicenseIds) != 0 {
			// ok have we already kept this one?
			for _, s := range seenIds {
				if l.LicenseId == s {
					seen = true
				}
			}
		}

		seenIds = append(seenIds, l.LicenseId)
		seenIds = append(seenIds, l.LicenseIds...)

		if !seen {
			fmt.Println("	keeping", l.LicenseId, l.LicenseIds)
			deduplicatedLicences = append(deduplicatedLicences, l)
		} else {
			fmt.Println("	duplicate! removing", l.LicenseId, l.LicenseIds)
		}
	}

	fmt.Println("from", len(licenses), "licences to", len(deduplicatedLicences))
	licenses = deduplicatedLicences

	fmt.Println("building ngrams for each license")
	for j := 0; j < len(licenses); j++ {
		split := strings.Fields(processor.LcCleanText(licenses[j].LicenseText))
		for i := startNgrams; i < endNgrams; i++ {
			ngrams := findNgrams(split, i)
			licenses[j].Ngrams = append(licenses[j].Ngrams, ngrams...)
		}
		licenses[j].LicenceNgrams = append(licenses[j].LicenceNgrams, licenses[j].Ngrams)

		// now calculate all the ngrams for the extra licence examples that we have
		for _, v := range licenses[j].ExtraLicenseText {
			split := strings.Fields(processor.LcCleanText(v))
			ngramsSlice := []string{}
			for i := startNgrams; i < endNgrams; i++ {
				ngrams := findNgrams(split, i)
				ngramsSlice = append(ngramsSlice, ngrams...)
			}
			licenses[j].LicenceNgrams = append(licenses[j].LicenceNgrams, licenses[j].Ngrams)
		}
	}

	// put every ngram into a huge map with a incrementing count, so if a ngram exists
	// and only has a count of 1 then we know it to be unique
	// TODO put this into the above loop so we can track uniqueness PER licence not just on keywords
	ngramCountMap := map[string]int{}
	for _, lic := range licenses {
		for _, ng := range lic.Ngrams {
			ngramCountMap[ng] = ngramCountMap[ng] + 1
		}
	}

	fmt.Println("finding unique ngrams")
	// store what we want to save here
	outputLicenses := []LicenseOutput{}

	// For each licence, check each ngram and see if it is unique
	for _, currentLicense := range licenses {
		// go through every ngram for this currentLicense and check that it does not occur anywhere else
		var uniqueNgrams []string
		for _, ngram := range currentLicense.Ngrams {
			// if its count is 1 that means its globally unique because it only exists in this
			// license as there is only a count of 1
			if ngramCountMap[ngram] == 1 {
				uniqueNgrams = append(uniqueNgrams, ngram)
			}
		}

		fmt.Println("	", currentLicense.LicenseId, "ngrams", len(currentLicense.Ngrams), "unique ngrams", len(uniqueNgrams))

		if len(uniqueNgrams) > keepNgrams {
			uniqueNgrams = uniqueNgrams[:keepNgrams]
		}

		outputLicenses = append(outputLicenses, LicenseOutput{
			LicenseTexts: currentLicense.ExtraLicenseText,
			Keywords:     uniqueNgrams,
			LicenseIds:   DedupeString(currentLicense.LicenseIds),
		})
	}

	out, _ := os.Create("database_keywords.json")

	data, _ := json.Marshal(outputLicenses)
	_, _ = out.Write(data)
	_ = out.Close()

	// now write out a list of every licence id that can be used for SPDX identification

	fmt.Println(fmt.Sprintf(`var spdxLicenseIds = []string{"%v"}`, strings.Join(allLicenseIds, `", "`)))
}

func DedupeString(s []string) []string {
	if len(s) < 2 {
		return s
	}
	sort.Slice(s, func(x, y int) bool { return s[x] > s[y] })
	var e = 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			continue
		}
		s[e] = s[i]
		e++
	}

	return s[:e]
}
