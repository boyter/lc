package parsers

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"strings"
)

func walkDirectoryFast(directory string, rootLicenses [][]LicenseMatch, output *chan *File) {
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
		*output <- &File {
			Directory: directory,
			File: file,
			RootLicenses: rootLicense,
		}
	}

	for _, newDirectory := range directories {
		walkDirectoryFast(filepath.Join(directory, newDirectory), rootLicenses, output)
	}
}

func processFileFast(input *chan *File, output *chan *FileResult) {
	var wg sync.WaitGroup

	for i := range *input {
		wg.Add(1)
		go func(file *File) {
			fileResult := processFile2(file.Directory, file.File, file.RootLicenses)
			*output <- &fileResult
			wg.Done()
		}(i)
	}

	wg.Wait()
	close(*input)
}

func processFile2(directory string, file string, rootLicenses []LicenseMatch) FileResult {
	process := true

	for _, ext := range strings.Split(ExtentionBlacklist, ",") {
		if strings.HasSuffix(file, "."+ext) {
			// Needs to be smarter we should skip reading the contents but it should still be under the license in the root folders
			process = false
		}
	}

	content := readFile(filepath.Join(directory, file))

	var licenseGuesses []LicenseMatch
	var licenseIdentified []LicenseMatch

	if len(content) > maxSize {
		process = false
	}

	if process == true {
		licenseGuesses = keywordGuessLicenseFast(content, Database)
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