package parsers

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Given a list of files use the possible license files list to find any that might be a match
func findPossibleLicenseFiles(fileList []string) []string {
	possibleList := []string{}

	for _, filename := range fileList {
		possible := false

		// Check agains the list that can be controlled by the user
		for _, indicator := range strings.Split(PossibleLicenceFiles, ",") {
			if strings.Contains(strings.ToLower(filename), indicator) {
				possible = true
			}
		}

		// Check against any that are called MIT, BSD etc....
		for _, license := range loadDatabase() {
			if strings.Split(filename, ".")[0] == strings.ToLower(license.LicenseId) {
				possible = true
			}
		}

		if possible == true {
			possibleList = append(possibleList, filename)
		}
	}

	return possibleList
}

func walkDirectory(output *chan FileJob, directory string, rootLicenses [][]LicenseMatch) []FileResult {
	fileResults := []FileResult{}
	all, _ := ioutil.ReadDir(directory)

	directories := []string{}
	files := []string{}

	// Work out which directories and files we want to investigate
	for _, f := range all {
		if f.IsDir() {
			add := true

			for _, black := range strings.Split(PathBlacklist, ",") {
				if f.Name() == black {
					add = false
					break
				}
			}

			if add == true {
				directories = append(directories, f.Name())
			}
		} else {
			files = append(files, f.Name())
		}
	}

	possibleLicenseFiles := findPossibleLicenseFiles(files)
	for _, file := range files {
		fileProcess(possibleLicenseFiles, directory, file, rootLicenses)
		*output <- FileJob{
			PossibleLicenses: possibleLicenseFiles,
			Directory:        directory,
			File:             file,
			RootLicenses:     rootLicenses,
		}
	}

	for _, newdirectory := range directories {
		walkDirectory(output, filepath.Join(directory, newdirectory), rootLicenses)
	}

	return fileResults
}
