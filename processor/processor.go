// SPDX-License-Identifier: MIT OR Unlicense
package processor

import (
	"fmt"
	file "github.com/boyter/go-code-walker"
	"io/ioutil"
)

var Version = "2.0.0 alpha"

// Set by user as command line arguments
var PossibleLicenceFiles = ""
var DirFilePaths []string
var PathBlacklist = ""

var FileOutput = ""
var ExtentionBlacklist = ""
var MaxSize = 50000
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

// Process is the main entry point of the command line output it sets everything up and starts running
func (process *Process) StartProcess() {
	lg := NewLicenceGuesser(true, true)

	fileListQueue := make(chan *file.File, 1000)

	fileWalker := file.NewFileWalker(".", fileListQueue)
	fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, "go")

	go fileWalker.Start()

	for f := range fileListQueue {
		data, err := ioutil.ReadFile(f.Location)
		if err == nil {
			fmt.Println()
			fmt.Println(f.Location)
			for _, x := range lg.SpdxIdentify(string(data)) {
				fmt.Println(x.LicenseId)
			}
		}

	}
}
