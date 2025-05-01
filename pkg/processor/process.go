package processor

import (
	"fmt"
	"os"
	"path/filepath"
)

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

	//fileWalker := gocodewalker.NewParallelFileWalker(dirPaths, potentialFilesQueue)
	//fileWalker.SetErrorHandler(func(e error) bool {
	//	printError(e.Error())
	//	return true
	//})
	//fileWalker.IgnoreGitIgnore = GitIgnore
	//fileWalker.IgnoreIgnoreFile = Ignore
	//fileWalker.IgnoreGitModules = GitModuleIgnore
	//fileWalker.IncludeHidden = true
	//fileWalker.ExcludeDirectory = PathDenyList
	//fileWalker.SetConcurrency(DirectoryWalkerJobWorkers)

}
