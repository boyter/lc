package main

import (
	"encoding/json"
	"fmt"
	vectorspace "github.com/boyter/golangvectorspace"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode/utf8"
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

func loadDatabase() []License {
	jsonFile, err := os.Open("database_keywords.json")

	if err != nil {
		fmt.Println(err)
		return []License{}
	}

	fmt.Println("Successfully Opened database_keywords.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var database []License
	err = json.Unmarshal(byteValue, &database)

	if err != nil {
		fmt.Println(err)
		return []License{}
	}

	for _, v := range database {
		println(v.Shortname, utf8.RuneCountInString(v.Text))
		v.Concordance = vectorspace.BuildConcordance(v.Text)
	}

	return database
}

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

	loadDatabase()

	// println(fileList)

	// for _, v := range fileList {
	// 	println(v)
	// }

	// filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
	// 	if !info.IsDir() {
	// 		println(info.Name())
	// 		println(path)

	// 		b, err := ioutil.ReadFile(path) // just pass the file name
	// 		if err != nil {
	// 			fmt.Print(err)
	// 		}
	// 		str := string(b) // convert content to a 'string'

	// 		var concordance = vectorspace.BuildConcordance(str)
	// 		println(concordance)
	// 	}

	// 	return nil
	// })
}
