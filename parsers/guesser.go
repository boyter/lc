package parsers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	vectorspace "github.com/boyter/golangvectorspace"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Shared all over the place
var ToolName = "licensechecker"
var ToolVersion = "2.0.0"

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
var Debug = false
var Trace = false

var spdxLicenceRegex = regexp.MustCompile(`SPDX-License-Identifier:\s+(.*)[ |\n|\r\n]*?`)
var alphaNumericRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
var multipleSpacesRegex = regexp.MustCompile("\\s+")

func cleanText(content string) string {
	content = strings.ToLower(content)

	content = alphaNumericRegex.ReplaceAllString(content, " ")
	content = multipleSpacesRegex.ReplaceAllString(content, " ")

	return content
}

// Identify licenses in the text which is using the SPDX indicator
// Can return multiple license matches
func identifierGuessLicence(content string, licenses []License) []LicenseMatch {
	matchingLicenses := []LicenseMatch{}
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

// Parses the supplied file content against the list of licences and
// returns the matching licences with the shortname and the percentage of match.
// If fast lookup methods fail it will try deep matching which is slower.
func guessLicense(content string, deepguess bool, licenses []License) []LicenseMatch {
	matchingLicenses := []LicenseMatch{}

	for _, license := range keywordGuessLicense([]byte(content), licenses) {
		matchingLicense := License{}

		for _, l := range licenses {
			if l.LicenseId == license.LicenseId {
				matchingLicense = l
				break
			}
		}

		runecontent := []rune(content)
		trimto := utf8.RuneCountInString(matchingLicense.LicenseText)

		if trimto > len(runecontent) {
			trimto = len(runecontent)
		}

		contentConcordance := vectorspace.BuildConcordance(string(runecontent[:trimto]))
		relation := vectorspace.Relation(matchingLicense.Concordance, contentConcordance)

		if relation >= confidence {
			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: relation})
		}
	}

	if len(matchingLicenses) == 0 && deepguess == true {
		for _, license := range licenses {
			runecontent := []rune(content)
			trimto := utf8.RuneCountInString(license.LicenseText)

			if trimto > len(runecontent) {
				trimto = len(runecontent)
			}

			contentConcordance := vectorspace.BuildConcordance(string(runecontent[:trimto]))
			relation := vectorspace.Relation(license.Concordance, contentConcordance)

			if relation >= confidence {
				matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: relation})
			}
		}
	}

	sort.Slice(matchingLicenses, func(i, j int) bool {
		return matchingLicenses[i].Percentage > matchingLicenses[j].Percentage
	})

	// Special cases such as MIT and JSON here
	if len(matchingLicenses) > 2 && ((matchingLicenses[0].LicenseId == "JSON" && matchingLicenses[1].LicenseId == "MIT") ||
		(matchingLicenses[0].LicenseId == "MIT" && matchingLicenses[1].LicenseId == "JSON")) {
		if strings.Contains(strings.ToLower(content), "not evil") {
			// Its JSON
			matchingLicenses = []LicenseMatch{}
			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: "JSON", Percentage: 1})
		} else {
			// Its MIT
			matchingLicenses = []LicenseMatch{}
			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: "MIT", Percentage: 1})
		}
	}

	return matchingLicenses
}

// Shamelessly stolen from https://github.com/src-d/go-license-detector
// https://github.com/src-d/go-license-detector#L63
// SPDX-License-Identifier: Apache-2.0
//var (
//	licenseFileNames = []string{
//		"li[cs]en[cs]e(s?)",
//		"legal",
//		"copy(left|right|ing)",
//		"unlicense",
//		"l?gpl([-_ v]?)(\\d\\.?\\d)?",
//		"bsd",
//		"mit",
//		"apache",
//		"readme",
//	}
//	licenseFileRe = regexp.MustCompile(
//		fmt.Sprintf("^(|.*[-_. ])(%s)(|[-_. ].*)$",
//			strings.Join(licenseFileNames, "|")))
//)
//
//func findPossibleLicenseFiles(fileList []string) []string {
//	var possibleList []string
//
//	for _, filename := range fileList {
//		if licenseFileRe.MatchString(strings.ToLower(filename)) {
//			possibleList = append(possibleList, filename)
//		}
//	}
//
//	return possibleList
//}

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

	for i, v := range database {
		database[i].Concordance = vectorspace.BuildConcordance(strings.ToLower(v.LicenseText))
	}

	Database = database

	if Trace {
		printTrace(fmt.Sprintf("milliseconds load database: %d", makeTimestampMilli()-startTime))
	}

	return database
}

func walkDirectory(directory string, rootLicenses [][]LicenseMatch) []FileResult {
	startTime := makeTimestampMilli()
	var fileResults []FileResult
	all, _ := ioutil.ReadDir(directory)

	var directories []string
	var files []string

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
	var identifiedRootLicense []LicenseMatch
	for _, possibleLicense := range possibleLicenses {
		content := string(readFile(filepath.Join(directory, possibleLicense)))
		guessLicenses := guessLicense(content, deepGuess, loadDatabase())

		if len(guessLicenses) != 0 {
			identifiedRootLicense = append(identifiedRootLicense, guessLicenses[0])
		}
	}

	if len(identifiedRootLicense) != 0 {
		rootLicenses = append(rootLicenses, identifiedRootLicense)
	}

	for _, file := range files {

		var rootLicense []LicenseMatch
		if len(rootLicenses) != 0 {
			rootLicense = rootLicenses[len(rootLicenses)-1]
		}

		fileResult := processFile(directory, file, rootLicense)
		fileResults = append(fileResults, fileResult)

		if strings.ToLower(Format) == "progress" {
			toProgress(directory, file, rootLicense, fileResult.LicenseGuesses, fileResult.LicenseIdentified)
		}
	}

	for _, newdirectory := range directories {
		results := walkDirectory(filepath.Join(directory, newdirectory), rootLicenses)
		fileResults = append(fileResults, results...)
	}

	if Trace {
		printTrace(fmt.Sprintf("milliseconds walk file tree: %s: %d", directory, makeTimestampMilli()-startTime))
	}

	return fileResults
}

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
	loadDatabase()

	var fileResults []FileResult

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

			if strings.ToLower(Format) == "progress" {
				toProgress(directory, file, []LicenseMatch{}, fileResult.LicenseGuesses, fileResult.LicenseIdentified)
			}
		}
	}

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
		fmt.Println("")
	}
}

func ProcessFast() {
	processArguments()
	loadDatabase()

	if len(DirFilePaths) == 0 {
		DirFilePaths = append(DirFilePaths, ".")
	}

	fileListQueue := make(chan *File, 5000)
	fileResultQueue := make(chan *FileResult, 5000)
	go processFileFast(&fileListQueue, &fileResultQueue)

	for _, fileDirectory := range DirFilePaths {
		if info, err := os.Stat(fileDirectory); err == nil && info.IsDir() {
			startTime := makeTimestampMilli()
			walkDirectoryFast(fileDirectory, [][]LicenseMatch{}, &fileListQueue)

			if Trace {
				printTrace(fmt.Sprintf("milliseconds walk file tree: %s: %d", fileDirectory, makeTimestampMilli()-startTime))
			}
		}
	}
}
