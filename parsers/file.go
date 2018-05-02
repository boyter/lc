package parsers

import (
	"io/ioutil"
	"path/filepath"
	"fmt"
)

func walkDirectoryFast(directory string) {
	all, _ := ioutil.ReadDir(directory)

	var directories []string
	var files []string

	// Work out which directories and files
	for _, f := range all {
		if f.IsDir() {
			directories = append(directories, f.Name())
		} else {
			files = append(files, f.Name())
		}
	}

	// Determine if any of the files might be a possible license deceleration
	possibleLicenseFiles := findPossibleLicenseFiles(files)

	// Determine the license for any of the possible files
	for _, file := range possibleLicenseFiles {

		bytes, err := ioutil.ReadFile(filepath.Join(directory, file))

		if err == nil {
			guessLicenses := keywordGuessLicenseFast(bytes, loadDatabase())
			fmt.Println(file, guessLicenses)

			//if len(guessLicenses) != 0 {
			//	identifiedRootLicense = append(identifiedRootLicense, guessLicenses[0])
			//}
		} else {
			// LOG ERROR
		}
	}

	//// Given the possible license files pass those and this file into channel for processing
	//for _, file := range files {
	//	fmt.Println(file, possibleLicenses)
	//}

	for _, newdirectory := range directories {
		walkDirectoryFast(filepath.Join(directory, newdirectory))
	}
}