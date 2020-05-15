// SPDX-License-Identifier: MIT
// SPDX-License-Identifier: Unlicense
package file

import (
	"os"
	"path/filepath"
	"strings"
)

// Walk the supplied directory backwards looking for .git or .hg
// directories indicating we should start our search from that
// location as its the root.
// Returns the first directory below supplied with .git or .hg in it
// otherwise the supplied directory
func FindRepositoryRoot(startDirectory string) string {
	// Firstly try to determine our real location
	curdir, err := os.Getwd()
	if err != nil {
		return startDirectory
	}

	// Check if we have .git or .hg where we are and if
	// so just return because we are already there
	if checkForGitOrMercurial(curdir) {
		return startDirectory
	}

	// We did not find something, so now we need to walk the file tree
	// backwards in a cross platform way and if we find
	// a match we return that
	lastIndex := strings.LastIndex(curdir, string(os.PathSeparator))
	for lastIndex != -1 {
		curdir = curdir[:lastIndex]

		if checkForGitOrMercurial(curdir) {
			return curdir
		}

		lastIndex = strings.LastIndex(curdir, string(os.PathSeparator))
	}

	// If we didn't find a good match return the supplied directory
	// so that we start the search from where we started at least
	// rather than the root
	return startDirectory
}

// Check if there is a .git or .hg folder in the supplied directory
func checkForGitOrMercurial(curdir string) bool {
	if stat, err := os.Stat(filepath.Join(curdir, ".git")); err == nil && stat.IsDir() {
		return true
	}

	if stat, err := os.Stat(filepath.Join(curdir, ".hg")); err == nil && stat.IsDir() {
		return true
	}

	return false
}

// A custom version of extracting extensions for a file
// which deals with extensions specific to code such as
// .travis.yml and the like
func GetExtension(name string) string {
	name = strings.ToLower(name)
	ext := filepath.Ext(name)

	if ext == "" || strings.LastIndex(name, ".") == 0 {
		ext = name
	} else {
		// Handling multiple dots or multiple extensions only needs to delete the last extension
		// and then call filepath.Ext.
		// If there are multiple extensions, it is the value of subExt,
		// otherwise subExt is an empty string.
		subExt := filepath.Ext(strings.TrimSuffix(name, ext))
		ext = strings.TrimPrefix(subExt+ext, ".")
	}

	return ext
}
