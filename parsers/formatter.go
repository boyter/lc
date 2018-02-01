// SPDX-License-Identifier: GPL-3.0-only

package parsers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ryanuber/columnize"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func toCSV(fileResults []FileResult) {
	records := [][]string{{
		"filename",
		"directory",
		"license",
		"confidence",
		"rootlicenses",
		"md5",
		"sha1",
		"sha256",
		"bytes",
		"byteshuman"},
	}

	for _, result := range fileResults {

		license := ""
		confidence := ""

		if len(result.LicenseGuesses) != 0 {
			license = result.LicenseGuesses[0].LicenseId
			confidence = fmt.Sprintf("%.3f", result.LicenseGuesses[0].Percentage*100)
		}

		rootLicenseString := ""
		for _, v := range result.LicenseRoots {
			rootLicenseString += fmt.Sprintf("%s,", v.LicenseId)
		}
		rootLicenseString = strings.TrimRight(rootLicenseString, ", ")

		records = append(records, []string{
			result.Filename,
			result.Directory,
			license,
			confidence,
			rootLicenseString,
			result.Md5Hash,
			result.Sha1Hash,
			result.Sha256Hash,
			strconv.Itoa(result.Bytes),
			result.BytesHuman})
	}

	csvfile, _ := os.OpenFile(FileOutput, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	defer csvfile.Close()

	w := csv.NewWriter(csvfile)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	fmt.Println("Results written to " + FileOutput)
}

func toJSON(fileResults []FileResult) {
	t, _ := json.Marshal(fileResults)
	ioutil.WriteFile(FileOutput, t, 0600)

	fmt.Println("Results written to " + FileOutput)
}

func toCli(fileResults []FileResult) {

	output := []string{
		"Directory | File | License | Confidence | Root Licenses | Size",
	}

	for _, result := range fileResults {
		license := ""
		confidence := ""

		if len(result.LicenseGuesses) != 0 {
			license = result.LicenseGuesses[0].LicenseId
			confidence = fmt.Sprintf("%.2f%%", result.LicenseGuesses[0].Percentage*100)
		}

		rootLicenseString := ""
		for _, v := range result.LicenseRoots {
			rootLicenseString += fmt.Sprintf("%s,", v.LicenseId)
		}
		rootLicenseString = strings.TrimRight(rootLicenseString, ", ")

		output = append(output, fmt.Sprintf("%s | %s | %s | %s | %s | %s", result.Directory, result.Filename, license, confidence, rootLicenseString, result.BytesHuman))
	}

	result := columnize.SimpleFormat(output)

	fmt.Println(result)
}

func toProgress(directory string, file string, rootLicenses []LicenseMatch, licenseGuesses []LicenseMatch) {
	license := ""
	confidence := ""

	if len(licenseGuesses) != 0 {
		license = licenseGuesses[0].LicenseId
		confidence = fmt.Sprintf("%.2f%%", licenseGuesses[0].Percentage*100)
	}

	rootLicenseString := ""
	for _, v := range rootLicenses {
		rootLicenseString += fmt.Sprintf("%s, ", v.LicenseId)
	}
	rootLicenseString = strings.TrimRight(rootLicenseString, ", ")

	fmt.Println("Filename:", file)
	fmt.Println("Directory:", directory)
	fmt.Println("License:", license, confidence)
	fmt.Println("Root License(s):", rootLicenseString)
	fmt.Println("----------------------------")
}

func toSPDX21(fileResults []FileResult) {

	for _, result := range fileResults {

		fmt.Println("")
		fmt.Println("FileName:", filepath.Join(result.Directory, result.Filename))
		fmt.Println("FileType: OTHER")
		fmt.Println("FileChecksum: SHA1:", result.Sha1Hash)
		fmt.Println("FileChecksum: SHA256:", result.Sha256Hash)
		fmt.Println("FileChecksum: MD5:", result.Md5Hash)
		fmt.Println("FileSize:", result.BytesHuman, "("+strconv.Itoa(result.Bytes)+" bytes)")

		// FileName: ./setproctitle.xs
		// FileType: OTHER
		// FileChecksum: SHA1: cc2d0d110e6a621f110a8bfb2fcf37499f99c2f3
		// FileChecksum: SHA256: f5c6e27f69ec93ffe803df83f18337aa341f56f388328444c022bad7c13ecb7c
		// FileChecksum: MD5: c0210487bec6a2997243e92694d77cee
		// FileChecksum: SSDEEP: 24:RS35k3ZJZutfhaVrEs7ahBrrrGsThM30IAzj8sRkm00csj2MdlwsWzGS3L:k8J6hqrr7krr6se30z9zdjh7oB3L
		// FileSize: 1 Kb (1178 bytes)
	}
}
