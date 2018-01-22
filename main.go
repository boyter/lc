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

const dirPath = "/home/bboyter/Projects/python-license-checker/"

type License struct {
	Keywords    []string `json:"keywords"`
	Text        string   `json:"text"`
	Fullname    string   `json:"fullname"`
	Shortname   string   `json:"shortname"`
	Header      string   `json:"header"`
	Concordance vectorspace.Concordance
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

func keywordGuessLicense(content string, licenses []License) string {

	var matchingName = ""
	var contains = false

	for _, license := range licenses {
		for _, keyword := range license.Keywords {
			contains = strings.Contains(content, keyword)
			if contains {
				matchingName = license.Shortname
				return matchingName
			}
		}
	}

	return matchingName
}

// def _keyword_guess(check_license, licenses):
//     matching = []

//     for license in licenses:
//         keywordmatch = 0
//         for keyword in license['keywords']:
//             if keyword in check_license:
//                 keywordmatch = keywordmatch + 1

//         if len(license['keywords']):
//             if keywordmatch >= 1:
//                 matching.append({
//                     'shortname': license['shortname'],
//                     'percentage': (float(keywordmatch) / float(len(license['keywords'])) * 100)
//                 })

//     return matching

func guessLicense() {

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

			fmt.Println(info.Name(), path, keywordGuessLicense(str, licenses))
		}

		return nil
	})
}
