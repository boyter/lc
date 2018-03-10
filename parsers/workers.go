package parsers

import (
	// "fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
)

func fileProcess(possibleLicenses []string, directory string, file string, rootLicenses [][]LicenseMatch) FileResult {
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

	rootLicense := []LicenseMatch{}
	if len(rootLicenses) != 0 {
		rootLicense = rootLicenses[len(rootLicenses)-1]
	}

	fileResult := processFile(directory, file, rootLicense)
	return fileResult
}

// Does the actual processing of stats and is the hot path
func fileProcessorWorker(input *chan FileJob, output *chan FileResult) {
	var wg sync.WaitGroup
	for res := range *input {
		wg.Add(1)
		go func(res FileJob) {
			result := fileProcess(res.PossibleLicenses, res.Directory, res.File, res.RootLicenses)
			*output <- result
			wg.Done()
		}(res)
	}

	go func() {
		wg.Wait()
		close(*output)
	}()
}
