package parsers

import (
	vectorspace "github.com/boyter/golangvectorspace"
	"path/filepath"
)

type License struct {
	Keywords    []string `json:"keywords"`
	LicenseText string   `json:"licenseText"`
	Name        string   `json:"name"`
	LicenseId   string   `json:"licenseId"`
	Concordance vectorspace.Concordance
}

type LicenseMatch struct {
	LicenseId  string
	Percentage float64
}

type FileResult struct {
	Directory         string
	Filename          string
	LicenseGuesses    []LicenseMatch
	LicenseRoots      []LicenseMatch
	LicenseIdentified []LicenseMatch
	Md5Hash           string
	Sha1Hash          string
	Sha256Hash        string
	BytesHuman        string
	Bytes             int
}

type File struct {
	Directory      string
	File           string
	RootLicenses   []LicenseMatch
	LicenseGuesses []LicenseMatch
}

func (c *FileResult) FullPath() string {
	return filepath.Join(c.Directory, c.Filename)
}
