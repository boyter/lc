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

// License contains details from loading off the SPDX list
type License struct {
	LicenseText             string `json:"licenseText"`
	StandardLicenseTemplate string `json:"standardLicenseTemplate"`
	StandardLicenseHeader   string `json:"standardLicenseHeader"`
	Name                    string `json:"name"`
	LicenseId               string `json:"licenseId"`

	// the below are used internally when processing
	//Ngrams           []string   // the unique bits of string we used to identify a licence
	//LicenceNgrams    [][]string // for each extra text + the main ngrams so we can get a selection from each
	//LicenseIds       []string   // indicates that this has 100% the same LicenseText such as AGPL-3.0-only and AGPL-3.0-or-later
	//ExtraLicenseText []string   // things like MIT have multiple variants of the same license and we need to track that
}

// LicenseConverted holds state that we use for processing before converted to output
type LicenseConverted struct {
	LicenceNgrams [][]string // for each extra text + the main ngrams so we can get a selection from each
	LicenseIds    []string   // indicates that this has 100% the same LicenseTexts such as AGPL-3.0-only and AGPL-3.0-or-later
	LicenseTexts  []string   // things like MIT have multiple variants of the same license and we need to track that
}

// LicenseOutput is the output format that we save to disk and import into lc
//type LicenseOutput struct {
//	LicenseTexts []string `json:"licenseTexts"` // examples of text that we have for these licences
//	LicenseIds   []string `json:"licenseIds"`   // SPDX ids where licences are considered identical
//	Keywords     []string `json:"keywords"`     // keywords that are unique and can be used to identify this group of licences
//}

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
	//var licenses []License
	var licenses []LicenseConverted
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

			// in flight representation of what we are doing
			licenseConverted := LicenseConverted{
				LicenceNgrams: [][]string{},
				LicenseIds:    []string{license.LicenseId},
				LicenseTexts:  []string{license.LicenseText},
			}

			// if MIT add in the other example so we can match it better...
			if license.LicenseId == "MIT" {
				fmt.Println("adding extra to MIT")
				licenseConverted.LicenseTexts = append(licenseConverted.LicenseTexts, mitExtra...)
			}

			if license.LicenseId == "BSD-3-Clause" {
				fmt.Println("adding extra to BSD-3-Clause")
				licenseConverted.LicenseTexts = append(licenseConverted.LicenseTexts, bsd3ClauseExtra...)
			}

			allLicenseIds = append(allLicenseIds, license.LicenseId)
			licenses = append(licenses, licenseConverted)

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
			// the first is the main one we are checking as that should always be the one from SPDX
			if processor.LcCleanText(lic.LicenseTexts[0]) == k {
				d = append(d, lic.LicenseIds[0])
			}
		}

		// update any license with this text to tell it about all the duplicates
		for i := 0; i < len(licenses); i++ {
			if processor.LcCleanText(licenses[i].LicenseTexts[0]) == k {
				licenses[i].LicenseIds = d
				fmt.Println(fmt.Sprintf("	%v %v duplicates %v", licenses[i].LicenseIds[0], len(d), d))
			}
		}
	}

	fmt.Println("removing duplicates...")
	// now we have identified the duplicates lets remove them such that we only have a single item
	var deduplicatedLicences []LicenseConverted
	var seenIds []string
	for _, l := range licenses {
		seen := false
		if len(l.LicenseIds) != 0 {
			// ok have we already kept this one?
			for _, s := range seenIds {
				if l.LicenseIds[0] == s {
					seen = true
				}
			}
		}

		seenIds = append(seenIds, l.LicenseIds[0])
		seenIds = append(seenIds, l.LicenseIds...)

		if !seen {
			fmt.Println("	keeping", l.LicenseIds[0], l.LicenseIds)
			deduplicatedLicences = append(deduplicatedLicences, l)
		} else {
			fmt.Println("	duplicate! removing", l.LicenseIds[0], l.LicenseIds)
		}
	}

	fmt.Println("from", len(licenses), "licences to", len(deduplicatedLicences))
	licenses = deduplicatedLicences

	fmt.Println("building ngrams for each license")
	for j := 0; j < len(licenses); j++ {
		// now calculate all the ngrams for the extra licence examples that we have
		for _, v := range licenses[j].LicenseTexts {
			split := strings.Fields(processor.LcCleanText(v))
			ngramsSlice := []string{}
			for i := startNgrams; i < endNgrams; i++ {
				ngrams := findNgrams(split, i)
				ngramsSlice = append(ngramsSlice, ngrams...)
			}
			licenses[j].LicenceNgrams = append(licenses[j].LicenceNgrams, ngramsSlice)
		}
	}

	// put every ngram into a huge map with a incrementing count, so if a ngram exists
	// and only has a count of 1 then we know it to be unique
	// TODO try to track uniqueness PER licence not just on keywords?
	ngramCountMap := map[string]int{}
	for _, lic := range licenses {
		for _, ngrams := range lic.LicenceNgrams {
			for _, ng := range ngrams {
				ngramCountMap[ng] = ngramCountMap[ng] + 1
			}
		}
	}

	fmt.Println("finding unique ngrams")
	// store what we want to save here
	outputLicenses := []processor.LicenseData{}

	// For each licence, check each ngram and see if it is unique
	for _, currentLicense := range licenses {
		// go through every ngram for this currentLicense and check that it does not occur anywhere else
		var uniqueNgrams []string
		for _, ngrams := range currentLicense.LicenceNgrams {
			for _, ngram := range ngrams {
				// if its count is 1 that means its globally unique because it only exists in this
				// license as there is only a count of 1
				if ngramCountMap[ngram] == 1 {
					uniqueNgrams = append(uniqueNgrams, ngram)
				}
			}
		}

		fmt.Println("	", currentLicense.LicenseIds, "ngrams", len(currentLicense.LicenceNgrams[0]), "unique ngrams", len(uniqueNgrams))

		if len(uniqueNgrams) > keepNgrams {
			uniqueNgrams = uniqueNgrams[:keepNgrams]
		}

		outputLicenses = append(outputLicenses, processor.LicenseData{
			LicenseTexts: currentLicense.LicenseTexts,
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
