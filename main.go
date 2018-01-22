package main

import (
	"fmt"
	vectorspace "github.com/boyter/golangvectorspace"
	"io/ioutil"
	"os"
	"path/filepath"
)

const dirPath = "/home/bboyter/Projects/python-license-checker/"

func main() {
	// walk all files in directory
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			println(info.Name())
			println(path)

			b, err := ioutil.ReadFile(path) // just pass the file name
			if err != nil {
				fmt.Print(err)
			}
			str := string(b) // convert content to a 'string'

			var concordance = vectorspace.BuildConcordance(str)
			println(concordance)

		}
		return nil
	})
}
