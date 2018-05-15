package parsers

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func walkDirectoryFast(directory string, rootLicenses [][]LicenseMatch, output *chan *File) {
	all, err := ioutil.ReadDir(directory)

	if err != nil {
		if Debug {
			printDebug(fmt.Sprintf("unable to read or directory: %s", directory))
		}
		return
	}

	var directories []string
	var files []string

	for _, file := range all {
		if file.IsDir() {
			add := true

			for _, black := range strings.Split(PathBlacklist, ",") {
				if file.Name() == black {
					add = false
				}
			}

			if add == true {
				directories = append(directories, file.Name())
			}
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
			guessLicenses := keywordGuessLicense(bytes, Database)

			if len(guessLicenses) != 0 {
				identifiedRootLicense = append(identifiedRootLicense, guessLicenses[0])
			}
		} else {
			if Debug {
				printDebug(fmt.Sprintf("unable to read or process file: %s", filepath.Join(directory, file)))
			}
		}
	}

	var rootLicense []LicenseMatch

	if len(identifiedRootLicense) != 0 {
		rootLicense = identifiedRootLicense
		rootLicenses = append(rootLicenses, identifiedRootLicense)
	} else if len(rootLicenses) != 0 {
		rootLicense = rootLicenses[len(rootLicenses)-1]
	}

	// Given the possible license files pass those and this file into channel for processing
	for _, file := range files {
		*output <- &File{
			Directory:    directory,
			File:         file,
			RootLicenses: rootLicense,
		}
	}

	//var wg sync.WaitGroup
	for _, newDirectory := range directories {
		//wg.Add(1)
		//go func(directory string, newDirectory string, rootlicenses [][]LicenseMatch, output *chan *File) {
			walkDirectoryFast(filepath.Join(directory, newDirectory), rootLicenses, output)
			//wg.Done()
		//}(directory, newDirectory, rootLicenses, output)
	}
	//wg.Wait()
}

func processFileFast(input *chan *File, output *chan *FileResult) {
	for file := range *input {
		fileResult := processFile(file.Directory, file.File, file.RootLicenses)
		*output <- &fileResult
	}
}

func processFile(directory string, file string, rootLicenses []LicenseMatch) FileResult {
	process := true

	for _, ext := range strings.Split(ExtentionBlacklist, ",") {
		if strings.HasSuffix(file, "."+ext) {
			// Needs to be smarter we should skip reading the contents but it should still be under the license in the root folders
			process = false
		}
	}

	content, err := ioutil.ReadFile(filepath.Join(directory, file))

	if err != nil {
		if Debug {
			printDebug(fmt.Sprintf("unable to read or process file: %s", filepath.Join(directory, file)))
		}
		process = false
	}

	var licenseGuesses []LicenseMatch
	var licenseIdentified []LicenseMatch

	if len(content) > maxSize {
		process = false
	}

	if process == true {
		licenseGuesses = keywordGuessLicense(content, Database)
		licenseIdentified = identifierGuessLicence(string(content), Database)
	}

	fileResult := FileResult{
		Directory:         directory,
		Filename:          file,
		LicenseGuesses:    licenseGuesses,
		LicenseRoots:      rootLicenses,
		LicenseIdentified: licenseIdentified,
		Md5Hash:           getMd5Hash(content),
		Sha1Hash:          getSha1Hash(content),
		Sha256Hash:        getSha256Hash(content),
		BytesHuman:        bytesToHuman(int64(len(content))),
		Bytes:             len(content)}

	return fileResult
}
