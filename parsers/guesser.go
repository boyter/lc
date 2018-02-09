package parsers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	vectorspace "github.com/boyter/golangvectorspace"
	"github.com/briandowns/spinner"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Shared all over the place
var ToolName = "licensechecker"
var ToolVersion = "1.0.0"

// Set by user as command line arguments
var confidence = 0.0
var Confidence = ""
var PossibleLicenceFiles = ""
var DirPath = "."
var PathBlacklist = ""
var deepGuess = true
var DeepGuess = ""
var Format = ""
var FileOutput = ""
var ExtentionBlacklist = ""
var maxSize = 0
var MaxSize = ""
var DocumentName = ""
var PackageName = ""
var DocumentNamespace = ""

var spdxLicenceRegex = regexp.MustCompile(`SPDX-License-Identifier:\s+(.*)[ |\n|\r\n]`)
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

// Fast method of checking if supplied content contains a licence using
// matching keyword ngrams to find if the licence is a match or not
// returns the maching licences with shortname and the percentage of match.
func keywordGuessLicense(content string, licenses []License) []LicenseMatch {
	content = cleanText(content)

	matchingLicenses := []LicenseMatch{}

	for _, license := range licenses {
		keywordmatch := 0
		contains := false

		for _, keyword := range license.Keywords {
			contains = strings.Contains(content, strings.ToLower(keyword))

			if contains == true {
				keywordmatch++
			}
		}

		if keywordmatch > 0 {
			percentage := (float64(keywordmatch) / float64(len(license.Keywords))) * 100
			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: percentage})
		}
	}

	return matchingLicenses
}

// Parses the supplied file content against the list of licences and
// returns the matching licences with the shortname and the percentage of match.
// If fast lookup methods fail it will try deep matching which is slower.
func guessLicense(content string, deepguess bool, licenses []License) []LicenseMatch {
	matchingLicenses := []LicenseMatch{}

	for _, license := range keywordGuessLicense(content, licenses) {
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

	return matchingLicenses
}

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

// Caching the database load result reduces processing time by about 3x for this repository
var Database = []License{}

func loadDatabase() []License {
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

	return database
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
		content := string(readFile(filepath.Join(directory, possibleLicense)))
		guessLicenses := guessLicense(content, deepGuess, loadDatabase())

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

		if strings.ToLower(Format) == "progress" {
			toProgress(directory, file, rootLicense, fileResult.LicenseGuesses, fileResult.LicenseIdentified)
		}
	}

	for _, newdirectory := range directories {
		results := walkDirectory(filepath.Join(directory, newdirectory), rootLicenses)
		fileResults = append(fileResults, results...)
	}

	return fileResults
}

func processFile(directory string, file string, rootLicenses []LicenseMatch) FileResult {
	process := true

	for _, ext := range strings.Split(ExtentionBlacklist, ",") {
		if strings.HasSuffix(file, "."+ext) {
			// Needs to be smarter we should skip reading the contents but it should still be under the license in the root folders
			process = false
		}
	}

	content := readFile(filepath.Join(directory, file))
	licenseGuesses := []LicenseMatch{}
	licenseIdentified := []LicenseMatch{}

	if len(content) > maxSize {
		process = false
	}

	if process == true {
		licenseGuesses = guessLicense(string(content), deepGuess, loadDatabase())
		licenseIdentified = identifierGuessLicence(string(content), loadDatabase())
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
		BytesHuman:        bytesToHuman(len(content)),
		Bytes:             len(content)}

	return fileResult
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
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Prefix = "Processing... "

	if strings.ToLower(Format) != "progress" && runtime.GOOS != "windows" {
		s.Start()
	}

	fileResults := []FileResult{}

	if DirPath == "" {
		DirPath = "."
	}

	if info, err := os.Stat(DirPath); err == nil && info.IsDir() {
		fileResults = walkDirectory(DirPath, [][]LicenseMatch{})
	} else {
		content := string(readFile(DirPath))
		guessLicenses := guessLicense(content, deepGuess, loadDatabase())
	}

	s.Stop()

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
