package parsers

import (
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Shared all over the place
var ToolName = "licensechecker"
var ToolVersion = "1.3.1"

// Set by user as command line arguments
var confidence = 0.85
var Confidence = ""
var PossibleLicenceFiles = ""
var DirFilePaths = []string{}
var PathBlacklist = ""
var deepGuess = true
var DeepGuess = ""
var Format = ""
var FileOutput = ""
var ExtentionBlacklist = ""
var maxSize = 50000
var MaxSize = ""
var DocumentName = ""
var PackageName = ""
var DocumentNamespace = ""

func processArguments() {
	conf, err := strconv.ParseFloat(Confidence, 64)
	if err == nil {
		confidence = conf
	} else {
		fmt.Println("Using default confidence value")
	}

	size, err := strconv.ParseInt(MaxSize, 10, 32)
	if err == nil {
		maxSize = int(size)
	} else {
		fmt.Println("Using default filesize value")
	}

	deep, err := strconv.ParseBool(DeepGuess)
	if err == nil {
		deepGuess = deep
	} else {
		fmt.Println("Using default deepguess value")
	}
}

func Process() {
	processArguments()
	loadDatabase() // To trigger the caching of it and allow for goroutines
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Prefix = "Processing... "

	if runtime.GOOS != "windows" {
		s.Start()
	}

	fileResults := []FileResult{}

	if len(DirFilePaths) == 0 {
		DirFilePaths = append(DirFilePaths, ".")
	}

	for _, fileDirectory := range DirFilePaths {
		if info, err := os.Stat(fileDirectory); err == nil && info.IsDir() {
			fileResults = append(fileResults, walkDirectory(fileDirectory, [][]LicenseMatch{})...)
		} else {
			directory, file := filepath.Split(fileDirectory)
			fileResult := processFile(directory, file, []LicenseMatch{})
			fileResults = append(fileResults, fileResult)
		}
	}
	s.Stop()

	switch strings.ToLower(Format) {
	case "csv":
		toCSV(fileResults)
	case "json":
		toJSON(fileResults)
	case "tabular":
		toTabular(fileResults)
	case "summary":
		toSummary(fileResults)
	case "spdx21":
		toSPDX21(fileResults)
	case "spdx":
		toSPDX21(fileResults)
	default:
		toTabular(fileResults)
	}
}
