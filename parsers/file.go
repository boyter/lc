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
	// NB this doubles the process time appox
	possibleLicenses := findPossibleLicenseFiles(files)

	// Determine the license for any of the possible files
	for _, possibleLicense := range possibleLicenses {

		bytes, err := ioutil.ReadFile(filepath.Join(directory, possibleLicense))

		if err == nil {
			guessLicenses := guessLicense(string(bytes), deepGuess, loadDatabase())

			fmt.Println(guessLicenses)
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