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

// License contains details from loading off the SPDX list
type License struct {
	LicenseText             string `json:"licenseText"`
	StandardLicenseTemplate string `json:"standardLicenseTemplate"`
	StandardLicenseHeader   string `json:"standardLicenseHeader"`
	Name                    string `json:"name"`
	LicenseId               string `json:"licenseId"`

	// the below are not standard but used internally when processing
	Ngrams           []string
	Duplicates       []string `json:"duplicates"`
	ExtraLicenseText []string
}

// LicenseOutput is the output format that we save to disk and import into lc
type LicenseOutput struct {
	LicenseText             string   `json:"licenseText"`
	StandardLicenseTemplate string   `json:"standardLicenseTemplate"`
	Name                    string   `json:"name"`
	LicenseId               string   `json:"licenseId"`
	Keywords                []string `json:"keywords"`
	Duplicates              []string `json:"duplicates"`
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

	same := map[string]int{}
	var licenses []License
	// Load all of the licenses from disk
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

			licenses = append(licenses, license)

			// track where the license text is the same
			same[license.LicenseText] = same[license.LicenseText] + 1
		}
	}

	fmt.Println("the following are duplicates... ie have the same license text...")
	for k, v := range same {
		if v != 1 {
			dupes := []string{}
			for _, lic := range licenses {
				if lic.LicenseText == k {
					dupes = append(dupes, lic.LicenseId)
				}
			}

			// update any license with this text to tell it about all the duplicates
			for i := 0; i < len(licenses); i++ {
				if licenses[i].LicenseText == k {
					licenses[i].Duplicates = dupes
					fmt.Println(fmt.Sprintf("	%v %v duplicates %v", licenses[i].LicenseId, len(dupes), dupes))
				}
			}
		}
	}

	fmt.Println("building ngrams for each license")
	// Build ngrams for each license
	for j := 0; j < len(licenses); j++ {
		split := strings.Fields(processor.LcCleanText(licenses[j].LicenseText))

		for i := startNgrams; i < endNgrams; i++ {
			ngrams := findNgrams(split, i)
			licenses[j].Ngrams = append(licenses[j].Ngrams, ngrams...)
		}

		// first get all ngrams for each text
		// store them all seperately
		// TODO store some from each one
		// the problem we have is that we have multiple texts...
		// but ngrams from the different texts might overlap, which
		// isnt a problem because they still refer to a single licence
		// so for each licence we need to get
		// we also need to ensue we keep keywords from each license we have
		// so that we can match, because we limit ourselves to some amount,
		// so we must iterate each one and mix them together to give the best chance

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
