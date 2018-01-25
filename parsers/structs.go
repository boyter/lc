package parsers

import (
	vectorspace "github.com/boyter/golangvectorspace"
)

type License struct {
	Keywords    []string `json:"keywords"`
	Text        string   `json:"text"`
	Fullname    string   `json:"fullname"`
	Shortname   string   `json:"shortname"`
	Header      string   `json:"header"`
	Concordance vectorspace.Concordance
}

type LicenseMatch struct {
	Shortname  string
	Percentage float64
}

type FileResult struct {
	Directory      string
	Filename       string
	LicenseGuesses []LicenseMatch
	LicenseRoots   []LicenseMatch
	Md5Hash        string
	Sha1Hash       string
	Sha256Hash     string
	BytesHuman     string
	Bytes          int
}
