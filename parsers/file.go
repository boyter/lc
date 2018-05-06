package parsers

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func walkDirectoryFast(directory string, rootLicenses [][]LicenseMatch, fileListQueue *chan *File) {
	all, err := ioutil.ReadDir(directory)

	if err != nil {
		// TODO log error
		return
	}

	var directories []string
	var files []string

	for _, file := range all {
		if file.IsDir() {
			directories = append(directories, file.Name())
		} else {
			files = append(files, file.Name())
		}
	}

	// Determine if any of the files might be a license
	possibleLicenseFiles := findPossibleLicenseFiles(files)
	var identifiedRootLicense []LicenseMatch

	// Determine the license for any of the possible files
	for _, file := range possibleLicenseFiles {
		bytes, err := ioutil.ReadFile(filepath.Join(directory, file))

		if err == nil {
			guessLicenses := keywordGuessLicenseFast(bytes, Database)

			fmt.Println(filepath.Join(directory, file), guessLicenses)

			if len(guessLicenses) != 0 {
				identifiedRootLicense = append(identifiedRootLicense, guessLicenses[0])
			}
		} else {
			// TODO log error
		}
	}

	var rootLicense []LicenseMatch
	if len(identifiedRootLicense) != 0 {
		rootLicense = identifiedRootLicense
	} else if len(rootLicenses) != 0 {
		rootLicense = rootLicenses[len(rootLicenses)-1]
	}

	// Given the possible license files pass those and this file into channel for processing
	for _, file := range files {
		*fileListQueue <- &File {
			Directory: directory,
			File: file,
			RootLicenses: rootLicense,
		}
	}

	for _, newdirectory := range directories {
		walkDirectoryFast(filepath.Join(directory, newdirectory), rootLicenses, fileListQueue)
	}
}

func processFileFast(directory string, file string, rootLicenses []LicenseMatch) {

}
