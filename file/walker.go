package file

import (
	"github.com/monochromegane/go-gitignore"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Location string
	Filename string
}

type FileWalker struct {
	directory              string
	fileListQueue          chan *File // holds all of the files it has foundÒ
	LocationExcludePattern []string // case sensitive patterns which exclude files
	PathExclude            []string // paths to always ignore such as .git,.svn and .hg
	IgnoreIgnoreFile       bool     // should .ignore files be respected?
	IgnoreGitIgnore        bool     // should .gitignore files be respected?
	IncludeHidden          bool     // should hidden files and directories be included/walked
	InstanceId             int
	AllowListExtensions    []string // which extensions should be allowed
}

func (f *FileWalker) walkDirectoryRecursive(directory string, ignores []gitignore.IgnoreMatcher) error {


	// walk the directory
	// identify any potential license files
	// if we have one try to identify it
	// then process any remaining files
	// the process future directries passing in licences

	fileInfos, err := ioutil.ReadDir(directory)

	if err != nil {
		return err
	}

	files := []os.FileInfo{}
	dirs := []os.FileInfo{}

	// We want to break apart the files and directories from the
	// return as we loop over them differently and this avoids some
	// nested if logic at the expense of a "redundant" loop
	for _, file := range fileInfos {
		if file.IsDir() {
			dirs = append(dirs, file)
		} else {
			files = append(files, file)
		}
	}

	// Pull out all of the ignore and gitignore files and add them
	// to out collection of ignores to be applied for this pass
	// and any subdirectories
	for _, file := range files {
		if !f.IgnoreGitIgnore {
			if file.Name() == ".gitignore" {
				ignore, err := gitignore.NewGitIgnore(filepath.Join(directory, file.Name()))
				if err == nil {
					ignores = append(ignores, ignore)
				}
			}
		}

		if !f.IgnoreIgnoreFile {
			if file.Name() == ".ignore" {
				ignore, err := gitignore.NewGitIgnore(filepath.Join(directory, file.Name()))
				if err == nil {
					ignores = append(ignores, ignore)
				}
			}
		}
	}

	// identify any files that could be a license file
	for _, file := range files {

	}

	// Process files first to start feeding whatever process is consuming
	// the output before traversing into directories for more files
	for _, file := range files {
		shouldIgnore := false

		// Check against the ignore files we have if the file we are looking at
		// should be ignored
		// It is safe to always call this because the ignores will not be added
		// in previous steps
		for _, ignore := range ignores {
			if ignore.Match(filepath.Join(directory, file.Name()), file.IsDir()) {
				shouldIgnore = true
			}
		}

		// Ignore hidden files
		if !f.IncludeHidden {
			s, err := IsHidden(file, directory)
			if s {
				shouldIgnore = true
			}
			if err != nil {
				return err
			}
		}

		// Check against extensions
		if len(f.AllowListExtensions) != 0 {
			ext := GetExtension(file.Name())
			a := false
			for _, v := range f.AllowListExtensions {
				if v == ext {
					a = true
				}
			}

			if !a {
				shouldIgnore = true
			}
		}

		if !shouldIgnore {
			for _, p := range f.LocationExcludePattern {
				if strings.Contains(filepath.Join(directory, file.Name()), p) {
					shouldIgnore = true
				}
			}

			if !shouldIgnore {
				f.fileListQueue <- &File{
					Location: filepath.Join(directory, file.Name()),
					Filename: file.Name(),
				}
			}
		}
	}

	// Now we process the directories after hopefully giving the
	// channel some files to process
	for _, dir := range dirs {
		shouldIgnore := false

		// Check against the ignore files we have if the file we are looking at
		// should be ignored
		// It is safe to always call this because the ignores will not be added
		// in previous steps
		for _, ignore := range ignores {
			if ignore.Match(filepath.Join(directory, dir.Name()), dir.IsDir()) {
				shouldIgnore = true
			}
		}

		// Confirm if there are any files in the path deny list which usually includes
		// things like .git .hg and .svn
		for _, deny := range f.PathExclude {
			if strings.HasSuffix(dir.Name(), deny) {
				shouldIgnore = true
			}
		}

		// Ignore hidden directories
		if !f.IncludeHidden {
			s, err := IsHidden(dir, directory)
			if s {
				shouldIgnore = true
			}
			if err != nil {
				return err
			}
		}

		if !shouldIgnore {
			for _, p := range f.LocationExcludePattern {
				if strings.Contains(filepath.Join(directory, dir.Name()), p) {
					shouldIgnore = true
				}
			}

			err = f.walkDirectoryRecursive(filepath.Join(directory, dir.Name()), ignores)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
