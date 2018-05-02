package parsers

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func walkDirectoryFast(directory string) {
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

	// Determine the license for any of the possible files
	for _, file := range possibleLicenseFiles {
		bytes, err := ioutil.ReadFile(filepath.Join(directory, file))

		if err == nil {
			startTime := makeTimestampMilli()
			guessLicenses := keywordGuessLicenseFast(bytes, Database)

			if Trace {
				printTrace(fmt.Sprintf("milliseconds to process file: %s: %d", filepath.Join(directory, file), makeTimestampMilli()-startTime))
			}

			fmt.Println(filepath.Join(directory, file), guessLicenses)

			//if len(guessLicenses) != 0 {
			//	identifiedRootLicense = append(identifiedRootLicense, guessLicenses[0])
			//}
		} else {
			// TODO log error
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
