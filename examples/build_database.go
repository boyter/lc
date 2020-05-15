package main

import (
	"encoding/json"
	"fmt"
	"github.com/boyter/lc/processor"
	"io/ioutil"
	"os"
	"path/filepath"
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
		license := licenses[i]

		// what we should do is get every ngram into a huge map EXCEPT for those from this license...
		// then for each one check if its in the map if it isnt its unique... more ram but SOOOO much faster
		ngramMap := map[string]bool{}
		for _, lic := range licenses {
			if license.LicenseId != lic.LicenseId {
				for _, ng := range lic.Ngrams {
					ngramMap[ng] = true
				}
			}
		}

		// go through every ngram for this license and check that it does not occur anywhere else
		var uniqueNgrams []string

		for _, ngram := range license.Ngrams {
			_, ok := ngramMap[ngram]

			if !ok {
				uniqueNgrams = append(uniqueNgrams, ngram)
			}
		}

		fmt.Println(license.LicenseId, "Ngrams", len(license.Ngrams), "Unique Ngrams", len(uniqueNgrams))

		outputLicenses = append(outputLicenses, LicenseOutput{
			LicenseId:               license.LicenseId,
			Name:                    license.Name,
			LicenseText:             license.LicenseText,
			StandardLicenseTemplate: license.StandardLicenseTemplate,
			Keywords:                uniqueNgrams,
		})
	}


	out, _ := os.Create("database_keywords.json")

	data, _ := json.Marshal(outputLicenses)
	_, _ = out.Write(data)
	_ = out.Close()
}
