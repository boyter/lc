package parsers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	//vectorspace "github.com/boyter/golangvectorspace"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// Shared all over the place
var ToolName = "licensechecker"
var ToolVersion = "2.0.0"

// Set by user as command line arguments
var PossibleLicenceFiles = ""
var DirFilePaths []string
var PathBlacklist = ""
var Format = ""
var FileOutput = ""
var ExtentionBlacklist = ""
var MaxSize = 50000
var DocumentName = ""
var PackageName = ""
var DocumentNamespace = ""
var Debug = false
var Trace = false

var spdxLicenceRegex = regexp.MustCompile(`SPDX-License-Identifier:\s+(.*)[ |\n|\r\n]*?`)
var alphaNumericRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
var multipleSpacesRegex = regexp.MustCompile("\\s+")

// Identify licenses in the text which is using the SPDX indicator
// Can return multiple license matches
func identifierGuessLicence(content string, licenses []License) []LicenseMatch {
	var matchingLicenses []LicenseMatch
	matches := spdxLicenceRegex.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		for _, license := range licenses {
			if license.LicenseId == strings.TrimSpace(val[1]) {
				matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: 1})
			}
		}
	}

	return matchingLicenses
}

func findPossibleLicenseFiles(fileList []string) []string {
	var possibleList []string

	for _, filename := range fileList {
		possible := false

		for _, indicator := range strings.Split(PossibleLicenceFiles, ",") {
			if strings.Contains(strings.ToLower(filename), indicator) {
				possible = true
			}
		}

		for _, license := range Database {
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

// Caching the database load result reduces processing time by about 3x for this repository
var Database []License

func loadDatabase() []License {
	startTime := makeTimestampMilli()
	if len(Database) != 0 {
		return Database
	}

	var database []License
	data, _ := base64.StdEncoding.DecodeString(database_keywords)
	_ = json.Unmarshal(data, &database)

	//for i, v := range database {
	//	database[i].Concordance = vectorspace.BuildConcordance(strings.ToLower(v.LicenseText))
	//}

	Database = database

	if Trace {
		printTrace(fmt.Sprintf("milliseconds load database: %d", makeTimestampMilli()-startTime))
	}

	return database
}

func Process() {
	loadDatabase()

	if len(DirFilePaths) == 0 {
		DirFilePaths = append(DirFilePaths, ".")
	}

	fileListQueue := make(chan *File, 5000)
	fileResultQueue := make(chan *FileResult, 5000)

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		for _, fileDirectory := range DirFilePaths {
			if info, err := os.Stat(fileDirectory); err == nil && info.IsDir() {
				startTime := makeTimestampMilli()
				walkDirectoryFast(fileDirectory, [][]LicenseMatch{}, &fileListQueue)

				if Trace {
					printTrace(fmt.Sprintf("milliseconds walk file tree: %s: %d", fileDirectory, makeTimestampMilli()-startTime))
				}
			}
		}
		wg.Done()
		close(fileListQueue)
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			processFileFast(&fileListQueue, &fileResultQueue)
			wg.Done()
		}()
	}

	wg.Wait()
	close(fileResultQueue)

	var fileResults []FileResult

	// TODO this needs to run in goroutine to avoid deadlock
	for input := range fileResultQueue {
		fileResults = append(fileResults, *input)
	}

	sort.Slice(fileResults, func(i, j int) bool {
		return strings.Compare(fileResults[i].FullPath(), fileResults[j].FullPath()) < 0
	})

	switch strings.ToLower(Format) {
	case "csv":
		toCSV(fileResults)
	case "json":
		toJSON(fileResults)
	case "tabular":
		toTabular(fileResults)
	case "spdx21":
		toSPDX21(fileResults)
	case "spdx":
		toSPDX21(fileResults)
	default:
		fmt.Println("")
	}
}
