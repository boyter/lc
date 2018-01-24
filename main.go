package main

import (
	"encoding/json"
	"fmt"
	"github.com/boyter/golang-license-checker/parsers"
	vectorspace "github.com/boyter/golangvectorspace"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const dirPath = "/home/bboyter/Projects/hyperfine/"
const pathBlacklist = "/.git/,/.hg/,/.svn/"
const licenceFiles = "license,copying"
const extentionBlacklist = "woff,eot,cur,dm,xpm,emz,db,scc,idx,mpp,dot,pspimage,stl,dml,wmf,rvm,resources,tlb,docx,doc,xls,xlsx,ppt,pptx,msg,vsd,chm,fm,book,dgn,blines,cab,lib,obj,jar,pdb,dll,bin,out,elf,so,msi,nupkg,pyc,ttf,woff2,jpg,jpeg,png,gif,bmp,psd,tif,tiff,yuv,ico,xls,xlsx,pdb,pdf,apk,com,exe,bz2,7z,tgz,rar,gz,zip,zipx,tar,rpm,bin,dmg,iso,vcd,mp3,flac,wma,wav,mid,m4a,3gp,flv,mov,mp4,mpg,rm,wmv,avi,m4v,sqlite,class,rlib,ncb,suo,opt,o,os,pch,pbm,pnm,ppm,pyd,pyo,raw,uyv,uyvy,xlsm,swf"

func loadDatabase(filepath string) []parsers.License {
	jsonFile, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
		return []parsers.License{}
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var database []parsers.License
	err = json.Unmarshal(byteValue, &database)

	if err != nil {
		fmt.Println(err)
		return []parsers.License{}
	}

	for i, v := range database {
		database[i].Concordance = vectorspace.BuildConcordance(strings.ToLower(v.Text))
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
		return nil
	}

	app.Run(os.Args)

	// Everything after here needs to be refactored out to a subpackage
	licenses := loadDatabase("database_keywords.json")

	files, _ := ioutil.ReadDir(dirPath)
	for _, f := range files {
		if f.IsDir() {
			fmt.Println("DIR>", f.Name())
		} else {
			fmt.Println("FIL>", f.Name())
		}
	}

	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			run := true

			for _, ext := range strings.Split(extentionBlacklist, ",") {
				if strings.HasSuffix(path, ext) {
					// Needs to be smarter we should skip reading the contents but it should still be under the license in the root folders
					run = false
				}
			}

			for _, black := range strings.Split(pathBlacklist, ",") {
				if strings.Contains(path, black) {
					run = false
				}
			}

			if run == true {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Print(err)
				}
				content := string(b)

				licenseGuesses := parsers.GuessLicense(content, true, licenses)

				licenseString := ""
				for _, v := range licenseGuesses {
					licenseString += fmt.Sprintf(" %s (%.1f%%)", v.Shortname, (v.Percentage * 100))
				}

				fmt.Println(path, licenseString)
			}
		}

		return nil
	})
}
