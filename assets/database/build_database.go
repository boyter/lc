package main

import (
	"encoding/json"
	"fmt"
	"github.com/boyter/lc/processor"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type License struct {
	LicenseText             string `json:"licenseText"`
	StandardLicenseTemplate string `json:"standardLicenseTemplate"`
	StandardLicenseHeader   string `json:"standardLicenseHeader"`
	Name                    string `json:"name"`
	LicenseId               string `json:"licenseId"`
	Ngrams                  []string
}

type LicenseOutput struct {
	LicenseText             string   `json:"licenseText"`
	StandardLicenseTemplate string   `json:"standardLicenseTemplate"`
	Name                    string   `json:"name"`
	LicenseId               string   `json:"licenseId"`
	Keywords                []string `json:"keywords"`
}

type NgramUnique struct {
	Ngram string
	Count int
}

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

func main() {
	files, _ := ioutil.ReadDir("./licenses/")

	fmt.Println("loading licenses")
	var licenses []License
	// Load the licenses
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := ioutil.ReadFile(filepath.Join("./licenses/", f.Name()))

			var license License
			json.Unmarshal(bytes, &license)
			license.Ngrams = []string{}

			licenses = append(licenses, license)
		}
	}

	var wg sync.WaitGroup
	fmt.Println("building ngrams for each license")
	// Build ngrams for each license
	for j := 0; j < len(licenses); j++ {
		wg.Add(1)
		go func(k int) {
			split := strings.Split(
				processor.LcCleanText(licenses[k].StandardLicenseHeader)+" "+
						processor.LcCleanText(licenses[k].LicenseText), " ")

			for i := 3; i < 7; i++ {
				ngrams := findNgrams(split, i)
				licenses[k].Ngrams = append(licenses[k].Ngrams, ngrams...)
			}
			wg.Done()
		}(j)
	}
	wg.Wait()

	fmt.Println("finding unique ngrams")

	outputLicenses := []LicenseOutput{}
	// For each licence, check each ngram and see if it is unique
	for i := 0; i < len(licenses); i++ {
		currentLicense := licenses[i]

		// what we should do is get every ngram into a huge map EXCEPT for those from this currentLicense...
		// then for each one check if its in the map if it isnt its unique... more ram but SOOOO much faster
		ngramMap := map[string]bool{}
		for _, lic := range licenses {
			if currentLicense.LicenseId != lic.LicenseId {
				for _, ng := range lic.Ngrams {
					ngramMap[ng] = true
				}
			}
		}

		// go through every ngram for this currentLicense and check that it does not occur anywhere else
		var uniqueNgrams []string

		for _, ngram := range currentLicense.Ngrams {
			_, ok := ngramMap[ngram]

			if !ok {
				uniqueNgrams = append(uniqueNgrams, ngram)
			}
		}

		// If we don't have anything we should try and find the MOST unique ones for this
		if len(uniqueNgrams) == 0 {
			mostlyUniqueNgrams := []NgramUnique{}
			var mostlyUniqueNgramsMutex sync.Mutex
			fmt.Println("finding mostly unique ngrams for", currentLicense.LicenseId)

			var wg sync.WaitGroup

			for i, ngram := range currentLicense.Ngrams {
				wg.Add(1)
				go func(i int, ngram string) {
					ngramCount := 0

					// go through every other ngram
					for _, lic := range licenses {
						if currentLicense.LicenseId != lic.LicenseId {
							for _, ng := range lic.Ngrams {
								if ng == ngram {
									ngramCount++
								}
							}
						}
					}

					mostlyUniqueNgramsMutex.Lock()
					mostlyUniqueNgrams = append(mostlyUniqueNgrams, NgramUnique{
						Ngram: ngram,
						Count: ngramCount,
					})
					mostlyUniqueNgramsMutex.Unlock()

					if i % 1000 == 0 {
						fmt.Println(i, "done for", len(currentLicense.Ngrams))
					}
					wg.Done()
				}(i, ngram)
			}
			wg.Wait()

			sort.Slice(mostlyUniqueNgrams, func(i, j int) bool {
				return mostlyUniqueNgrams[i].Count > mostlyUniqueNgrams[j].Count
			})

			for _, m := range mostlyUniqueNgrams {
				if m.Count < 10 {
					uniqueNgrams = append(uniqueNgrams, m.Ngram)
				}
			}
		}


		fmt.Println(currentLicense.LicenseId, "Ngrams", len(currentLicense.Ngrams), "Unique Ngrams", len(uniqueNgrams))

		if len(uniqueNgrams) > 100 {
			uniqueNgrams = uniqueNgrams[:100]
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
}
