package parsers

import (
	"code.cloudfoundry.org/bytefmt"
	"encoding/base64"
	"encoding/json"
	"fmt"
	vectorspace "github.com/boyter/golangvectorspace"
	"github.com/briandowns/spinner"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var confidence = 0.85
var Confidence = "0.85"
var PossibleLicenceFiles = "license,copying,readme"
var DirPath = "/home/bboyter/Projects/hyperfine/"
var PathBlacklist = ".git,.hg,.svn"
var deepGuess = true
var DeepGuess = "true"
var Format = "cli"
var FileOutput = ""

// Will not attempt tp process but will still list under
var ExtentionBlacklist = "woff,eot,cur,dm,xpm,emz,db,scc,idx,mpp,dot,pspimage,stl,dml,wmf,rvm,resources,tlb,docx,doc,xls,xlsx,ppt,pptx,msg,vsd,chm,fm,book,dgn,blines,cab,lib,obj,jar,pdb,dll,bin,out,elf,so,msi,nupkg,pyc,ttf,woff2,jpg,jpeg,png,gif,bmp,psd,tif,tiff,yuv,ico,xls,xlsx,pdb,pdf,apk,com,exe,bz2,7z,tgz,rar,gz,zip,zipx,tar,rpm,bin,dmg,iso,vcd,mp3,flac,wma,wav,mid,m4a,3gp,flv,mov,mp4,mpg,rm,wmv,avi,m4v,sqlite,class,rlib,ncb,suo,opt,o,os,pch,pbm,pnm,ppm,pyd,pyo,raw,uyv,uyvy,xlsm,swf"

func cleanText(content string) string {
	content = strings.ToLower(content)

	alphaNumeric := regexp.MustCompile("[^a-zA-Z0-9 ]")
	multipleSpaces := regexp.MustCompile("\\s+")

	content = alphaNumeric.ReplaceAllString(content, " ")
	content = multipleSpaces.ReplaceAllString(content, " ")

	return content
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

		if possible == true {
			possibleList = append(possibleList, filename)
		}
	}

	return possibleList
}

func loadDatabase() []License {
	var database []License
	data, _ := base64.StdEncoding.DecodeString(database_keywords)
	_ = json.Unmarshal(data, &database)

	for i, v := range database {
		database[i].Concordance = vectorspace.BuildConcordance(strings.ToLower(v.LicenseText))
	}

	return database
}

func walkDirectory(directory string, rootLicenses []LicenseMatch) []FileResult {
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
	for _, possibleLicense := range possibleLicenses {
		content := string(readFile(filepath.Join(directory, possibleLicense)))
		guessLicenses := guessLicense(content, deepGuess, loadDatabase())

		if len(guessLicenses) != 0 {
			rootLicenses = append(rootLicenses, guessLicenses[0])
		}
	}

	for _, file := range files {
		process := true

		for _, ext := range strings.Split(ExtentionBlacklist, ",") {
			if strings.HasSuffix(file, ext) {
				// Needs to be smarter we should skip reading the contents but it should still be under the license in the root folders
				process = false
			}
		}

		content := readFile(filepath.Join(directory, file))
		licenseGuesses := []LicenseMatch{}
		if process == true {
			licenseGuesses = guessLicense(string(content), deepGuess, loadDatabase())
		}

		fileResult := FileResult{
			Directory:      directory,
			Filename:       file,
			LicenseGuesses: licenseGuesses,
			LicenseRoots:   rootLicenses,
			Md5Hash:        getMd5Hash(content),
			Sha1Hash:       getSha1Hash(content),
			Sha256Hash:     getSha256Hash(content),
			BytesHuman:     bytefmt.ByteSize(uint64(len(content))),
			Bytes:          len(content)}

		fileResults = append(fileResults, fileResult)

		if strings.ToLower(Format) == "progress" {
			license := ""
			confidence := ""

			if len(licenseGuesses) != 0 {
				license = licenseGuesses[0].LicenseId
				confidence = fmt.Sprintf("%.2f%%", licenseGuesses[0].Percentage*100)
			}

			rootLicenseString := ""
			for _, v := range rootLicenses {
				rootLicenseString += fmt.Sprintf("%s,", v.LicenseId)
			}
			rootLicenseString = strings.TrimRight(rootLicenseString, ", ")

			fmt.Println(directory, file, license, confidence, rootLicenseString, bytefmt.ByteSize(uint64(len(content))))
		}

	}

	for _, newdirectory := range directories {
		results := walkDirectory(filepath.Join(directory, newdirectory), rootLicenses)
		fileResults = append(fileResults, results...)
	}

	return fileResults
}

func Process() {
	conf, err := strconv.ParseFloat(Confidence, 64)

	if err == nil {
		confidence = conf
	} else {
		fmt.Println("Using default confidence value")
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Suffix = " processing"

	if strings.ToLower(Format) != "progress" {
		s.Start()
	}

	fileResults := []FileResult{}
	if info, err := os.Stat(DirPath); err == nil && info.IsDir() {
		fileResults = walkDirectory(DirPath, []LicenseMatch{})
	} else {
		content := string(readFile(DirPath))
		guessLicenses := guessLicense(content, deepGuess, loadDatabase())
		fmt.Println(guessLicenses)
	}

	s.Stop()

	switch strings.ToLower(Format) {
	case "csv":
		toCSV(fileResults)
	case "json":
		toJSON(fileResults)
	case "cli":
		toCli(fileResults)
	default:
		fmt.Println("")
	}
}
