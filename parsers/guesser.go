package parsers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/texttheater/golang-levenshtein/levenshtein"
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
				matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Score: 1})
			}
		}
	}

	return matchingLicenses
}

// Given a list of files attempts to determine if they might contain
// a software licence and if so returns those that match
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
var CommonDatabase []License

// LoadDatabase will initialize the database values and should only be called once such as in an init
func LoadDatabase() []License {
	startTime := makeTimestampMilli()
	if len(Database) != 0 {
		return Database
	}

	var database []License
	data, _ := base64.StdEncoding.DecodeString(database_keywords)
	_ = json.Unmarshal(data, &database)

	Database = database

	// Load smaller faster database for checking most common licenses
	common := []string{
		"MIT",
		"Apache-2.0",
		"GPL-3.0",
		"AGPL-3.0",
		"BSD-3-Clause",
		"GPL-2.0",
		"BSD-2-Clause",
		"CC0-1.0",
		"LGPL-3.0",
		"LGPL-2.1",
		"ISC",
		"0BSD",
		"LGPL-2.0",
		"Unlicense",
		"BSD-3-Clause-No-Nuclear-License-2014",
		"MPL-2.0",
		"EPL-1.0",
		"MPL-2.0-no-copyleft-exception",
		"AGPL-1.0",
		"CC-BY-4.0",
		"IPL-1.0",
		"CPL-1.0",
		"CC-BY-3.0",
		"CC-BY-SA-4.0",
		"WTFPL",
		"Zlib",
		"CC-BY-SA-3.0",
		"Cube",
		"JSON",
		"BitTorrent-1.0",
	}

	for _, license := range database {
		for _, com := range common {
			if license.LicenseId == com {
				CommonDatabase = append(CommonDatabase, license)
			}
		}
	}

	if Trace {
		printTrace(fmt.Sprintf("milliseconds load database: %d", makeTimestampMilli()-startTime))
	}

	return database
}

// fast method of checking if supplied content contains a licence using
// matching keyword ngrams to find if the licence is a match or not
// returns the matching licences with shortname and the percentage of match.
func keywordGuessLicense(content []byte, licenses []License) []LicenseMatch {
	content = cleanText(content)

	var wg sync.WaitGroup
	output := make(chan LicenseMatch, len(licenses))

	for _, license := range licenses {
		wg.Add(1)
		go func(license License) {
			keywordMatch := 0

			for _, keyword := range license.Keywords {
				if bytes.Contains(content, []byte(keyword)) {
					keywordMatch++
				}
			}

			if keywordMatch > 100 {
				distance := levenshtein.DistanceForStrings([]rune(string(content)), []rune(string(cleanText([]byte(license.LicenseText)))), levenshtein.DefaultOptions)
				output <- LicenseMatch{LicenseId: license.LicenseId, Score: float64(distance)}
			}
			wg.Done()
		}(license)
	}

	wg.Wait()
	close(output)

	var matchingLicenses []LicenseMatch
	for license := range output {
		matchingLicenses = append(matchingLicenses, license)
	}

	sort.Slice(matchingLicenses, func(i, j int) bool {
		// For keywordMatch we want > but for distance we want <
		return matchingLicenses[i].Score < matchingLicenses[j].Score
	})

	matchingLicenses = specialCases(content, matchingLicenses)

	return matchingLicenses
}

func cleanText(content []byte) []byte {
	content = bytes.ToLower(content)

	tmp := alphaNumericRegex.ReplaceAllString(string(content), " ")
	tmp = multipleSpacesRegex.ReplaceAllString(tmp, " ")
	tmp = strings.TrimSpace(tmp)

	return []byte(tmp)
}

func specialCases(content []byte, matchingLicenses []LicenseMatch) []LicenseMatch {
	// Quite often JSON and MIT are confused
	if len(matchingLicenses) > 2 && ((matchingLicenses[0].LicenseId == "JSON" && matchingLicenses[1].LicenseId == "MIT") ||
		(matchingLicenses[0].LicenseId == "MIT" && matchingLicenses[1].LicenseId == "JSON")) {
		if bytes.Contains(content, []byte("not evil")) {
			matchingLicenses = []LicenseMatch{{LicenseId: "JSON", Score: 1}}
		} else {
			matchingLicenses = []LicenseMatch{{LicenseId: "MIT", Score: 1}}
		}
	}

	// Another one is MIT-feh and MIT
	if len(matchingLicenses) > 2 && matchingLicenses[0].LicenseId == "MIT" {
		if bytes.HasPrefix(content, []byte("mit license")) || bytes.HasPrefix(content, []byte("the mit license")) {
			matchingLicenses = []LicenseMatch{{LicenseId: "MIT", Score: 100}}
		} else {
			matchingLicenses = []LicenseMatch{{LicenseId: "MIT-feh", Score: 100}}
		}
	}

	return matchingLicenses
}

func Process() {
	LoadDatabase()

	if len(DirFilePaths) == 0 {
		DirFilePaths = append(DirFilePaths, ".")
	}

	fileListQueue := make(chan *File, 5000)
	fileResultQueue := make(chan *FileResult, 5000)

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		for _, fileDirectory := range DirFilePaths {
			if info, err := os.Stat(fileDirectory); err == nil {
				if info.IsDir() {
					startTime := makeTimestampMilli()
					walkDirectory(fileDirectory, [][]LicenseMatch{}, &fileListQueue)

					if Trace {
						printTrace(fmt.Sprintf("milliseconds walk file tree: %s: %d", fileDirectory, makeTimestampMilli()-startTime))
					}
				} else {
					fileListQueue <- &File{
						Directory:      filepath.Dir(fileDirectory),
						File:           info.Name(),
						RootLicenses:   []LicenseMatch{},
						LicenseGuesses: []LicenseMatch{},
					}
				}
			}
		}
		wg.Done()
		close(fileListQueue)
	}()

	go func() {
		wg.Add(1)
		processFileFast(&fileListQueue, &fileResultQueue)
		wg.Done()
		close(fileResultQueue)
	}()

	var fileResults []FileResult

	// TODO this needs to run in goroutine to avoid deadlock
	for input := range fileResultQueue {
		fileResults = append(fileResults, *input)
	}
	wg.Wait()

	sort.Slice(fileResults, func(i, j int) bool {
		return strings.Compare(fileResults[i].FullPath(), fileResults[j].FullPath()) < 0
	})

	switch strings.ToLower(Format) {
	case "csv":
		toCSV(fileResults)
	case "xlsx":
		toXLSX(fileResults)
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
