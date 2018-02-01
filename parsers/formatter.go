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
	"time"
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

	// Determine the package licenses
	packageLicenseDeclared := "NONE"

	if len(fileResults) != 0 {
		if len(fileResults[0].LicenseRoots) == 1 {
			packageLicenseDeclared = fileResults[0].LicenseRoots[0].LicenseId
		} else if len(fileResults) >= 2 {
			rootLicenseNames := []string{}
			for _, v := range fileResults[0].LicenseRoots {
				rootLicenseNames = append(rootLicenseNames, v.LicenseId)
			}
			packageLicenseDeclared = "(" + strings.Join(rootLicenseNames, " AND ") + ")"
		}
	}

	fmt.Println("SPDXVersion: SPDX-2.1")
	fmt.Println("DataLicense: CC0-1.0")
	fmt.Println("SPDXID: SPDXRef-DOCUMENT")
	fmt.Println("DocumentName: DOCUMENTNAMEHEREFROMCLI")                                                         // TODO
	fmt.Println("DocumentNamespace:http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82...") // TODO
	fmt.Println("LicenseListVersion: 3.0")
	fmt.Println("Creator: Tool:", ToolName, ToolVersion)
	fmt.Println("Created:", time.Now().UTC().Format(time.RFC3339))

	fmt.Println("")
	fmt.Println("PackageName: TODO")             // TODO
	fmt.Println("SPDXID: SPDXRef-1")             // TODO
	fmt.Println("PackageDownloadLocation: NONE") // TODO pass in from CLI https://spdx.org/spdx-specification-21-web-version#h.49x2ik5
	fmt.Println("FilesAnalyzed: true")
	fmt.Println("PackageVerificationCode: TODO") // TODO https://spdx.org/spdx-specification-21-web-version#h.2p2csry
	fmt.Println("PackageLicenseDeclared:", packageLicenseDeclared)
	fmt.Println("")

	// Loop over all files and get a list of all unique licenses and print below
	// PackageLicenseInfoFromFiles: GPL-2.0

	for _, result := range fileResults {

		// TODO this needs to possibly be NOASSERTION if unsure
		licenseConcluded := "NONE"

		if len(result.LicenseGuesses) != 0 {
			licenseConcluded = result.LicenseGuesses[0].LicenseId
		}

		fmt.Println("FileName:", filepath.Join(result.Directory, result.Filename))
		fmt.Println("FileType: OTHER")
		fmt.Println("FileChecksum: SHA1:", result.Sha1Hash)
		fmt.Println("FileChecksum: SHA256:", result.Sha256Hash)
		fmt.Println("FileChecksum: MD5:", result.Md5Hash)
		fmt.Println("LicenseConcluded:", licenseConcluded)
		fmt.Println("FileCopyrightText: NOASSERTION")
		fmt.Println("FileSize:", result.BytesHuman, "("+strconv.Itoa(result.Bytes)+" bytes)")
		fmt.Println("")
	}
}
