// SPDX-License-Identifier: MIT OR Unlicense

package main

import (
	"encoding/json"
	"fmt"
	"github.com/boyter/lc/processor"
	"os"
	"path/filepath"
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
	Ngrams              []string   // the unique bits of string we used to identify a licence
	LicenceNgrams       [][]string // for each extra text + the main ngrams so we can get a selection from each
	LicenseIdDuplicates []string   `json:"duplicates"` // indicates that this has 100% the same LicenseText such as AGPL-3.0-only and AGPL-3.0-or-later
	ExtraLicenseText    []string   // things like MIT have multiple variants of the same license and we need to track that
}

// LicenseOutput is the output format that we save to disk and import into lc
type LicenseOutput struct {
	LicenseText             string   `json:"licenseText"`
	StandardLicenseTemplate string   `json:"standardLicenseTemplate"`
	Name                    string   `json:"name"`
	LicenseId               string   `json:"licenseId"`
	Keywords                []string `json:"keywords"`
	LicenseIdDuplicates     []string `json:"duplicates"`
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
var endNgrams = 17
var keepNgrams = 200

func main() {
	// find the licence files that we need to compare against as a starting
	// point based on the SPDX
	files, _ := os.ReadDir("./licenses/")

	fmt.Println("loading licenses")

	licenseTextCount := map[string]int{}
	var licenses []License
	// Load all of the licenses from disk and keep track of duplicates
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := os.ReadFile(filepath.Join("./licenses/", f.Name()))

			var license License
			_ = json.Unmarshal(bytes, &license)
			license.Ngrams = []string{}

			// if MIT add in the other example so we can match it better...
			if license.LicenseId == "MIT" {
				fmt.Println("adding extra to MIT")
				license.ExtraLicenseText = mitExtra
			}

			if license.LicenseId == "BSD-3-Clause" {
				fmt.Println("adding extra to BSD-3-Clause")
				license.ExtraLicenseText = bsd3ClauseExtra
			}

			licenses = append(licenses, license)

			// track where the license text is the licenseTextCount
			licenseTextCount[processor.LcCleanText(license.LicenseText)] = licenseTextCount[processor.LcCleanText(license.LicenseText)] + 1
		}
	}

	fmt.Println("the following are duplicates...")
	for k, v := range licenseTextCount {
		if v != 1 {
			var d []string
			for _, lic := range licenses {
				if processor.LcCleanText(lic.LicenseText) == k {
					d = append(d, lic.LicenseId)
				}
			}

			// update any license with this text to tell it about all the duplicates
			for i := 0; i < len(licenses); i++ {
				if processor.LcCleanText(licenses[i].LicenseText) == k {
					licenses[i].LicenseIdDuplicates = d
					fmt.Println(fmt.Sprintf("	%v %v duplicates %v", licenses[i].LicenseId, len(d), d))
				}
			}
		}
	}

	fmt.Println("removing duplicates...")
	// now we have identified the duplicates lets remove them such that we only have a single item
	var deduplicatedLicences []License
	var seenIds []string
	for _, l := range licenses {
		seen := false
		if len(l.LicenseIdDuplicates) != 0 {
			// ok have we already kept this one?
			for _, s := range seenIds {
				if l.LicenseId == s {
					seen = true
				}
			}
		}

		seenIds = append(seenIds, l.LicenseId)
		seenIds = append(seenIds, l.LicenseIdDuplicates...)

		if !seen {
			fmt.Println("	keeping", l.LicenseId)
			deduplicatedLicences = append(deduplicatedLicences, l)
		} else {
			fmt.Println("	duplicate! removing", l.LicenseId)
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

		fmt.Println(currentLicense.LicenseId, "ngrams", len(currentLicense.Ngrams), "unique ngrams", len(uniqueNgrams))

		if len(uniqueNgrams) != 0 {
			if len(uniqueNgrams) > keepNgrams {
				uniqueNgrams = uniqueNgrams[:keepNgrams]
			}

			outputLicenses = append(outputLicenses, LicenseOutput{
				LicenseId:               currentLicense.LicenseId,
				Name:                    currentLicense.Name,
				LicenseText:             currentLicense.LicenseText,
				StandardLicenseTemplate: currentLicense.StandardLicenseTemplate,
				Keywords:                uniqueNgrams,
			})
		}
	}

	out, _ := os.Create("database_keywords.json")

	data, _ := json.Marshal(outputLicenses)
	_, _ = out.Write(data)
	_ = out.Close()

	/*
		//// Write out
		//files, _ = ioutil.ReadDir(".")
		//out, _ = os.Create("./database.go")
		//
		//// Open constants
		//out.Write([]byte("package processor \n\nvar LicenseDatabase = []License{\n"))
		//for _, f := range outputLicenses {
		//
		//	key := ""
		//	if len(f.Keywords) != 0 {
		//		for _, k := range f.Keywords {
		//			key += fmt.Sprintf(`"%s",`, k)
		//		}
		//	}
		//
		//	out.Write(bytes.Trim([]byte(fmt.Sprintf(`{
		//		LicenseText:             ` + "`" + `%s` + "`" + `,
		//		StandardLicenseTemplate: ` + "`" + `%s` + "`" + `,
		//		Name:                    ` + "`" + `%s` + "`" + `,
		//		LicenseId:               ` + "`" + `%s` + "`" + `,
		//		Keywords:                []string{
		//			%s
		//		},
		//	},`,
		//		strings.Replace(f.LicenseText, "`", "` + \"`\" + `", -1),
		//		strings.Replace(f.StandardLicenseTemplate, "`", "` + \"`\" + `", -1),
		//		strings.Replace(f.Name, "`", "` + \"`\" + `", -1),
		//		strings.Replace(f.LicenseId, "`", "` + \"`\" + `", -1),
		//		key)), "\xef\xbb\xbf"))
		//}
		//
		//out.Write([]byte("}\n"))
		//out.Close()
	*/
}
