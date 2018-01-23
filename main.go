package main

import (
	"encoding/json"
	"fmt"
	vectorspace "github.com/boyter/golangvectorspace"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const dirPath = "/home/bboyter/Projects/hyperfine/"

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

func loadDatabase(filepath string) []License {
	jsonFile, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
		return []License{}
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var database []License
	err = json.Unmarshal(byteValue, &database)

	if err != nil {
		fmt.Println(err)
		return []License{}
	}

	for _, v := range database {
		v.Concordance = vectorspace.BuildConcordance(v.Text)
	}

	return database
}

func keywordGuessLicense(content string, licenses []License) []LicenseMatch {
	content = strings.ToLower(content)
	var matchingLicenses = []LicenseMatch{}

	for _, license := range licenses {
		var keywordmatch = 0
		var contains = false

		for _, keyword := range license.Keywords {
			contains = strings.Contains(content, keyword)
			if contains {
				keywordmatch++
			}
		}

		if keywordmatch > 0 {
			var percentage = (float64(keywordmatch) / float64(len(license.Keywords))) * 100
			matchingLicenses = append(matchingLicenses, LicenseMatch{Shortname: license.Shortname, Percentage: percentage})
		}
	}

	return matchingLicenses
}

func guessLicense(content string, licenses []License) {
	var matchingLicenses = keywordGuessLicense(content, licenses)

	for _, license := range matchingLicenses {
		vectorspace.BuildConcordance(license.Shortname)
		fmt.Println(license.Shortname)
	}

}

// def guess_license(check_license, licenses):
//     matching = _keyword_guess(check_license, licenses)

//     matches = []
//     vector_compare = VectorCompare()
//     for match in matching:
//         for license in [x for x in licenses if x['shortname'] in [y['shortname'] for y in matching]]:
//             licence_concordance = vector_compare.concordance(license['clean'])

//             check_license_concordance = vector_compare.concordance(check_license[:len(license['clean'])])

//             relation = vector_compare.relation(license['concordance'], check_license_concordance)

//             if relation >= 0.85:
//                 matches.append({
//                     'relation': relation,
//                     'license': license
//                 })

//     matches.sort(reverse=True)

//     return matches

func main() {
	// walk all files in directory

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "golang-license-checker"
	app.Version = "1.0"
	app.Usage = "Check directory for licenses and list what license(s) a file is under"
	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}

	app.Run(os.Args)

	fileList := []string{}
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	licenses := loadDatabase("database_keywords.json")

	// println(fileList)

	// for _, v := range fileList {
	// 	println(v)
	// }

	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			b, err := ioutil.ReadFile(path) // just pass the file name
			if err != nil {
				fmt.Print(err)
			}
			str := string(b) // convert content to a 'string'

			var guesses = keywordGuessLicense(str, licenses)
			guessLicense(str, licenses)

			for _, v := range guesses {
				fmt.Println(path, v.Shortname, v.Percentage)
			}

			// fmt.Println(path)
		}

		return nil
	})
}
