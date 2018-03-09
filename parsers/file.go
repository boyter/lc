package parsers

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func findPossibleLicenseFiles(fileList []string) []string {
	possibleList := []string{}

	for _, filename := range fileList {
		possible := false

		for _, indicator := range strings.Split(PossibleLicenceFiles, ",") {
			if strings.Contains(strings.ToLower(filename), indicator) {
				possible = true
			}
		}

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

func walkDirectory(directory string, rootLicenses [][]LicenseMatch) []FileResult {
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
				}
			}

			if add == true {
				directories = append(directories, f.Name())
			}
		} else {
			files = append(files, f.Name())
		}
	}

	// Determine any possible licence files which would classify everything else
	possibleLicenses := findPossibleLicenseFiles(files)
	identifiedRootLicense := []LicenseMatch{}
	for _, possibleLicense := range possibleLicenses {
		contents, _ := ioutil.ReadFile(filepath.Join(directory, possibleLicense))
		guessLicenses := guessLicense(string(contents), deepGuess, loadDatabase())

		if len(guessLicenses) != 0 {
			identifiedRootLicense = append(identifiedRootLicense, guessLicenses[0])
		}
	}

	if len(identifiedRootLicense) != 0 {
		rootLicenses = append(rootLicenses, identifiedRootLicense)
	}

	// TODO fan this out to many GoRoutines and process in parallel
	for _, file := range files {

		rootLicense := []LicenseMatch{}
		if len(rootLicenses) != 0 {
			rootLicense = rootLicenses[len(rootLicenses)-1]
		}

		fileResult := processFile(directory, file, rootLicense)
		fileResults = append(fileResults, fileResult)
	}

	for _, newdirectory := range directories {
		results := walkDirectory(filepath.Join(directory, newdirectory), rootLicenses)
		fileResults = append(fileResults, results...)
	}

	return fileResults
}
