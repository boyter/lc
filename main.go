package main

import (
	"encoding/json"
	"fmt"
	vectorspace "github.com/boyter/golangvectorspace"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

const dirPath = "/home/bboyter/Projects/hyperfine/"
const pathBlacklist = "/.git/,/.hg/,/.svn/"
const licenceFiles = "license,copying"
const extentionBlacklist = "woff,eot,cur,dm,xpm,emz,db,scc,idx,mpp,dot,pspimage,stl,dml,wmf,rvm,resources,tlb,docx,doc,xls,xlsx,ppt,pptx,msg,vsd,chm,fm,book,dgn,blines,cab,lib,obj,jar,pdb,dll,bin,out,elf,so,msi,nupkg,pyc,ttf,woff2,jpg,jpeg,png,gif,bmp,psd,tif,tiff,yuv,ico,xls,xlsx,pdb,pdf,apk,com,exe,bz2,7z,tgz,rar,gz,zip,zipx,tar,rpm,bin,dmg,iso,vcd,mp3,flac,wma,wav,mid,m4a,3gp,flv,mov,mp4,mpg,rm,wmv,avi,m4v,sqlite,class,rlib,ncb,suo,opt,o,os,pch,pbm,pnm,ppm,pyd,pyo,raw,uyv,uyvy,xlsm,swf"
const confidence = 0.85

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

	for i, v := range database {
		database[i].Concordance = vectorspace.BuildConcordance(strings.ToLower(v.Text))
	}

	return database
}

func keywordGuessLicense(content string, licenses []License) []LicenseMatch {
	content = strings.ToLower(content)
	matchingLicenses := []LicenseMatch{}

	for _, license := range licenses {
		keywordmatch := 0
		contains := false

		for _, keyword := range license.Keywords {
			contains = strings.Contains(content, keyword)
			if contains {
				keywordmatch++
			}
		}

		if keywordmatch > 0 {
			percentage := (float64(keywordmatch) / float64(len(license.Keywords))) * 100
			matchingLicenses = append(matchingLicenses, LicenseMatch{Shortname: license.Shortname, Percentage: percentage})
		}
	}

	return matchingLicenses
}

func guessLicense(content string, licenses []License) []LicenseMatch {
	matchingLicenses := []LicenseMatch{}

	for _, license := range keywordGuessLicense(content, licenses) {

		matchingLicense := License{}

		for _, l := range licenses {
			if l.Shortname == license.Shortname {
				matchingLicense = l
				break
			}
		}

		runecontent := []rune(content)
		trimto := utf8.RuneCountInString(matchingLicense.Text)

		if trimto > len(runecontent) {
			trimto = len(runecontent)
		}

		contentConcordance := vectorspace.BuildConcordance(string(runecontent[:trimto]))
		relation := vectorspace.Relation(matchingLicense.Concordance, contentConcordance)

		if relation >= confidence {
			matchingLicenses = append(matchingLicenses, LicenseMatch{Shortname: license.Shortname, Percentage: relation})
		}
	}

	if len(matchingLicenses) == 0 {
		for _, license := range licenses {
			runecontent := []rune(content)
			trimto := utf8.RuneCountInString(license.Text)

			if trimto > len(runecontent) {
				trimto = len(runecontent)
			}

			contentConcordance := vectorspace.BuildConcordance(string(runecontent[:trimto]))
			relation := vectorspace.Relation(license.Concordance, contentConcordance)

			if relation >= confidence {
				matchingLicenses = append(matchingLicenses, LicenseMatch{Shortname: license.Shortname, Percentage: relation})
			}
		}
	}

	sort.Slice(matchingLicenses, func(i, j int) bool {
		return matchingLicenses[i].Percentage > matchingLicenses[j].Percentage
	})

	return matchingLicenses
}

func main() {
	// walk all files in directory

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "golang-license-checker"
	app.Version = "1.0"
	app.Usage = "Check directory for licenses and list what license(s) a file is under"
	app.Action = func(c *cli.Context) error {
		return nil
	}

	app.Run(os.Args)

	// Everything after here needs to be refactored out to a subpackage
	licenses := loadDatabase("database_keywords.json")

	extentionBlacklistStrings := strings.Split(extentionBlacklist, ",")

	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			run := true

			for _, ext := range extentionBlacklistStrings {
				if strings.HasSuffix(path, ext) {
					// Needs to be smarter we should skip reading the contents but it should still be under the license in the root folders
					run = false
				}
			}

			if strings.Contains(path, "/.git/") {
				run = false
			}

			if run == true {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Print(err)
				}
				content := string(b)

				licenseGuesses := guessLicense(content, licenses)

				licenseString := ""
				for _, v := range licenseGuesses {
					licenseString += fmt.Sprintf(" %s (%.2f)", v.Shortname, v.Percentage)
				}

				fmt.Println(path, licenseString)
			}
		}

		return nil
	})
}
