// SPDX-License-Identifier: MIT OR Unlicense
package processor

import (
	"fmt"
	"os"

	"github.com/boyter/gocodewalker"
)

var Version = "2.0.0 alpha"

// Set by user as command line arguments
var PossibleLicenceFiles = ""
var DirFilePaths []string
var PathBlacklist = ""

var FileOutput = ""
var ExtentionBlacklist = ""
var MaxSize = 50_000
var DocumentName = ""
var PackageName = ""
var DocumentNamespace = ""
var Debug = false
var Trace = false

var IncludeBinaryFiles = false
var IgnoreIgnoreFile = false
var IgnoreGitIgnore = false
var IncludeHidden = false
var AllowListExtensions []string
var Format = ""

type Process struct {
	Directory string // What directory are we searching
	FindRoot  bool
}

func NewProcess(directory string) Process {
	return Process{
		Directory: directory,
	}
}

// StartProcess is the main entry point of the command line output it sets everything up and starts running
func (process *Process) StartProcess() {
	lg := NewLicenceGuesser(true, true)
	lg.UseFullDatabase = false
	flg := NewLicenceGuesser(true, true)
	flg.UseFullDatabase = true

	fileListQueue := make(chan *gocodewalker.File, 1000)

	fileWalker := gocodewalker.NewFileWalker(".", fileListQueue)
	fileWalker.IgnoreGitIgnore = false
	fileWalker.IgnoreIgnoreFile = false
	//fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, "go")

	go func() { _ = fileWalker.Start() }()

	for f := range fileListQueue {
		data, err := os.ReadFile(f.Location)
		if err == nil {

			// TODO should be configurable and be in the read file to avoid doing it at all
			// TODO actually it shouldn't even read it unless the filename is useful
			if len(data) > 100_000 {
				data = data[:100_000]
			}

			if process.isBinary(data) {
				continue
			}

			licenceFile := IsLicenceFile(f.Filename)

			if licenceFile {
				fmt.Println(f.Location)
				// should we should boost the guesses here because we are fairly sure there is a licence in there?
				license := lg.GuessLicense(data)

				// if we didn't find anything try using everything
				if len(license) == 0 {
					license = flg.GuessLicense(data)
				}

				if len(license) == 0 {
					fmt.Println(" possible licence file but unable to identify")
				}
				for _, x := range license {
					fmt.Println(fmt.Sprintf(" Licence: %s (%.1f%%)", x.LicenseId, x.ScorePercentage))
				}
			} else {
				// look for SPDX markers only as its not a licence file
				license := lg.SpdxIdentify(string(data))
				if len(license) != 0 {
					fmt.Println(f.Location)
				}
				for _, x := range license {
					fmt.Println(fmt.Sprintf(" Licence: %s (100%%)", x))
				}
			}
		}

	}
}

// Helper function that looks through supplied bytes looking for null which indicates
// it is a binary file and returns true/false
func (process *Process) isBinary(data []byte) bool {
	// Check if this content is binary by checking for null bytes and if found assume it is binary
	// this is how GNU Grep, git and ripgrep check for binary files
	for _, b := range data {
		if b == 0 {
			return true
		}
	}

	return false
}
