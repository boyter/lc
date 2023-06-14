package parsers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func walkGitDirectory(directory string, files *[]string) error {
	var stdout bytes.Buffer

	// Check for a valid git command in the PATH
	_, err := exec.LookPath("git")
	if err != nil {
		// Fall back to flat directory walking
		return errors.New("Git command not found")
	}

	cmd := exec.Command("git", "ls-files")
	cmd.Stdout = &stdout

	// Make sure the command is executed from within the repository
	// directory in order to avoid out-of-repository errors from git
	cmd.Dir = directory

	err = cmd.Run()
	if err != nil {
		return errors.New("unable to obtain file list from git")
	}

	for {
		file, err := stdout.ReadString('\n')
		if err != nil {
			break
		}

		file = strings.TrimRight(file, "\r\n")
		*files = append(*files, file)
	}

	if err != nil && err != io.EOF {
		// Clear out any partially constructed file list
		files = nil

		return errors.New("failed to iterate file list from git")
	}

	return nil
}

func walkFlatDirectory(directory string, files *[]string, directories *[]string) error {
	all, err := ioutil.ReadDir(directory)

	if err != nil {
		return errors.New("unable to read or directory: " + directory)
	}

	for _, file := range all {
		if file.IsDir() {
			add := true

			for _, black := range strings.Split(PathBlacklist, ",") {
				if file.Name() == black {
					add = false
				}
			}

			if add == true {
				*directories = append(*directories, file.Name())
			}
		} else {
			*files = append(*files, file.Name())
		}
	}

	return nil
}

func walkDirectory(directory string, rootLicenses [][]LicenseMatch, output *chan *File) {
	var directories []string
	var files []string

	// In the case where the directory is a valid git repository, prefer
	// extracting the file list from git directly. This saves us from
	// needing to walk directories recursively, and has the added benefit
	// of excluding files not included in the distribution but which may
	// still be present in the working directory.
	if _, err := os.Stat(directory + "/.git"); !os.IsNotExist(err) {
		err = walkGitDirectory(directory, &files)
		if err != nil {
			if Debug {
				printDebug(err.Error())
			}

		}
	}

	// If earlier directory walking methods have failed, walk manually
	if len(files) == 0 {
		err := walkFlatDirectory(directory, &files, &directories)
		if err != nil {
			if Debug {
				printDebug(err.Error())
			}
		}
	}

	// Determine if any of the files might be a license
	possibleLicenseFiles := findPossibleLicenseFiles(files)
	var identifiedRootLicense []LicenseMatch

	// Determine the license for any of the possible files
	for _, file := range possibleLicenseFiles {
		bytes, err := os.ReadFile(filepath.Join(directory, file))

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
			Directory:      directory,
			File:           file,
			RootLicenses:   rootLicense,
			LicenseGuesses: identifiedRootLicense,
		}
	}

	var wg sync.WaitGroup
	for _, newDirectory := range directories {
		wg.Add(1)
		go func(directory string, newDirectory string, rootlicenses [][]LicenseMatch, output *chan *File) {
			walkDirectory(filepath.Join(directory, newDirectory), rootLicenses, output)
			wg.Done()
		}(directory, newDirectory, rootLicenses, output)
	}
	wg.Wait()
}

func processFileFast(input *chan *File, output *chan *FileResult) {
	for file := range *input {
		fileResult := processFile(file.Directory, file.File, file.LicenseGuesses, file.RootLicenses)
		*output <- &fileResult
	}
}

// GuessLicense will attempt to find any licenses in the content and return them as a sorted set of guesses with the highest
// score first (if more than one)
func GuessLicense(content []byte) []LicenseMatch {
	if len(content) > MaxSize {
		return []LicenseMatch{}
	}
	licenseGuesses := keywordGuessLicense(content, CommonDatabase)
	if len(licenseGuesses) == 0 {
		licenseGuesses = keywordGuessLicense(content, Database)
	}
	licenseGuesses = append(licenseGuesses, identifierGuessLicence(string(content), Database)...)
	if len(licenseGuesses) > 1 {
		sort.Slice(licenseGuesses, func(i, j int) bool {
			// For keywordMatch we want > but for distance we want <
			return licenseGuesses[i].Score < licenseGuesses[j].Score
		})
	}
	return licenseGuesses
}

func processFile(directory string, file string, guessed []LicenseMatch, rootLicenses []LicenseMatch) FileResult {
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

	if len(content) > MaxSize {
		process = false
	}

	// If we have something here it means it was identified as a license file
	// so we don't need to reprocess
	if len(guessed) != 0 {
		licenseGuesses = guessed
		process = false
	}

	if process == true {
		licenseGuesses = keywordGuessLicense(content, CommonDatabase)
		if len(licenseGuesses) == 0 {
			licenseGuesses = keywordGuessLicense(content, Database)
		}

		licenseIdentified = identifierGuessLicence(string(content), Database)
	}

	md5 := ""
	sha1 := ""
	sha256 := ""

	if strings.ToLower(Format) == "csv" || strings.ToLower(Format) == "xlsx" || strings.ToLower(Format) == "json" || strings.ToLower(Format) == "spdx21" || strings.ToLower(Format) == "spdx" {
		md5 = getMd5Hash(content)
		sha1 = getSha1Hash(content)
		sha256 = getSha256Hash(content)
	}

	fileResult := FileResult{
		Directory:         directory,
		Filename:          file,
		LicenseGuesses:    licenseGuesses,
		LicenseRoots:      rootLicenses,
		LicenseIdentified: licenseIdentified,
		Md5Hash:           md5,
		Sha1Hash:          sha1,
		Sha256Hash:        sha256,
		BytesHuman:        bytesToHuman(int64(len(content))),
		Bytes:             len(content)}

	return fileResult
}
