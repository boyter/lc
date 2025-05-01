package processor

import (
	"fmt"
	"github.com/boyter/gocodewalker"
	"os"
	"path/filepath"
	"runtime"
)

// PathDenyList sets the paths that should be skipped
var PathDenyList = []string{}

// DirFilePaths is not set via flags but by arguments following the flags for file or directory to process
var DirFilePaths = []string{}

func Process() {
	// Clean up any invalid arguments before setting everything up
	if len(DirFilePaths) == 0 {
		DirFilePaths = append(DirFilePaths, ".")
	}

	filePaths := []string{}
	dirPaths := []string{}

	// Check if the paths or files added exist and exit if not
	for _, f := range DirFilePaths {
		fpath := filepath.Clean(f)

		s, err := os.Stat(fpath)
		if err != nil {
			fmt.Println("file or directory could not be read: " + fpath)
			os.Exit(1)
		}

		if s.IsDir() {
			dirPaths = append(dirPaths, fpath)
		} else {
			filePaths = append(filePaths, fpath)
		}
	}

	potentialFilesQueue := make(chan *gocodewalker.File, runtime.NumCPU()) // files that pass the .gitignore checks

	fileWalker := gocodewalker.NewParallelFileWalker(dirPaths, potentialFilesQueue)
	fileWalker.SetErrorHandler(func(e error) bool {
		fmt.Println(e.Error())
		return true
	})
	//fileWalker.IgnoreGitIgnore = GitIgnore
	//fileWalker.IgnoreIgnoreFile = Ignore
	//fileWalker.IgnoreGitModules = GitModuleIgnore
	fileWalker.IncludeHidden = true
	fileWalker.ExcludeDirectory = PathDenyList

	go func() {
		err := fileWalker.Start()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	for fi := range potentialFilesQueue {
		fileInfo, err := os.Lstat(fi.Location)
		if err != nil {
			continue
		}

		if !fileInfo.IsDir() {
			fmt.Println(fi.Location, fi.Filename, fileInfo)
		}
	}

}
