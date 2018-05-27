package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type License struct {
	LicenseText             string `json:"licenseText"`
	StandardLicenseTemplate string `json:"standardLicenseTemplate"`
	Name                    string `json:"name"`
	LicenseId               string `json:"licenseId"`
	Ngrams                  [][]string
}

type LicenseOutput struct {
	LicenseText             string   `json:"licenseText"`
	StandardLicenseTemplate string   `json:"standardLicenseTemplate"`
	Name                    string   `json:"name"`
	LicenseId               string   `json:"licenseId"`
	Keywords                []string `json:"keywords"`
}

var spdxLicenceRegex = regexp.MustCompile(`SPDX-License-Identifier:\s+(.*)[ |\n|\r\n]*?`)
var alphaNumericRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
var multipleSpacesRegex = regexp.MustCompile("\\s+")

func cleanText(content string) string {
	content = strings.ToLower(content)

	content = alphaNumericRegex.ReplaceAllString(content, " ")
	content = multipleSpacesRegex.ReplaceAllString(content, " ")

	return content
}

func findNgrams(list []string, size int) [][]string {
	var ngrams [][]string

	for i := 0; i < len(list); i++ {

		if i+size < len(list)+1 {
			ngram := list[i : i+size]
			ngrams = append(ngrams, ngram)
		}
	}

	return ngrams
}

func main() {
	files, _ := ioutil.ReadDir("./licenses/")

	var licenses []License

	// Load the licenses
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := ioutil.ReadFile(filepath.Join("./licenses/", f.Name()))

			var license License
			json.Unmarshal(bytes, &license)
			license.Ngrams = [][]string{}

			licenses = append(licenses, license)
		}
	}

	// Build ngrams for them
	for j := 0; j < len(licenses); j++ {
		split := strings.Split(cleanText(licenses[j].LicenseText), " ")

		for i := 2; i < 45; i++ { // 45 seems about right
			ngrams := findNgrams(split, i)
			licenses[j].Ngrams = append(licenses[j].Ngrams, ngrams...)
		}
	}

	//var licensesOutput []LicenseOutput

	var wg sync.WaitGroup
	// For each licence, check each ngram and see if it is unique
	for i := 0; i < len(licenses); i++ {
		wg.Add(1)
		go func(i int) {
			license := licenses[i]

			// for each licence that isn't this one
			// get all the ngrams and put it into a hash
			// then look each of our ngrams and check if it is contained
			contains := map[string]int{}
			for _, lic := range licenses {
				if lic.LicenseId != license.LicenseId {
					for _, ngram := range lic.Ngrams {
						contains[strings.Join(ngram, "")] = 1
					}
				}
			}

			var uniqueNgrams []string
			for _, ngram := range license.Ngrams {
				_, ok := contains[strings.Join(ngram, "")]

				if !ok {
					uniqueNgrams = append(uniqueNgrams, strings.Join(ngram, " "))
				}
			}

			fmt.Println(license.LicenseId, len(license.Ngrams), "Unique Ngrams", len(uniqueNgrams))
			wg.Done()
		}(i)

		// This is such an ugly way to control concurrency
		if i%4 == 0 { // This uses about 8 GB of RAM currently.
			time.Sleep(time.Second * 5) // Gotta wait for the goroutines to actually kick off first
			fmt.Println("Waiting for results...")
			wg.Wait()
		}
	}

	wg.Wait()
	out, _ := os.Create("database_keywords.json")

	data, _ := json.Marshal(licenses)
	out.Write(data)
	out.Close()

	//licenses = [{
	//	'licenseText': x['licenseText'],
	//	'name': x['name'],
	//	'licenseId': x['licenseId'],
	//	'keywords': x['keywords'][:50]
	//} for x in licenses]

}
