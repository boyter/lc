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

	if FileOutput == "" {
		w := csv.NewWriter(os.Stdout)
		w.WriteAll(records) // calls Flush internally
		w.Flush()
	} else {
		csvfile, _ := os.OpenFile(FileOutput, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		defer csvfile.Close()

		w := csv.NewWriter(csvfile)
		w.WriteAll(records) // calls Flush internally

		if err := w.Error(); err != nil {
			log.Fatalln("error writing csv:", err)
		}

		fmt.Println("Results written to " + FileOutput)
	}
}

func toJSON(fileResults []FileResult) {
	t, _ := json.Marshal(fileResults)

	if FileOutput == "" {
		fmt.Println(string(t))
	} else {
		ioutil.WriteFile(FileOutput, t, 0600)
		fmt.Println("Results written to " + FileOutput)
	}
}

func joinLicenseList(licenseList []LicenseMatch, operator string) string {
	licenseDeclared := ""

	if len(licenseList) == 1 {
		licenseDeclared = licenseList[0].LicenseId
	} else if len(licenseList) >= 2 {
		licenseNames := []string{}
		for _, v := range licenseList {
			licenseNames = append(licenseNames, v.LicenseId)
		}
		licenseDeclared = "(" + strings.Join(licenseNames, operator) + ")"
	}

	return licenseDeclared
}

func toTabular(fileResults []FileResult) {
	output := []string{
		"Directory | File | License | Confidence | Root Licenses | Size",
	}

	for _, result := range fileResults {
		license := ""
		confidence := ""

		if len(result.LicenseIdentified) != 0 {
			license = joinLicenseList(result.LicenseIdentified, " AND ")
			confidence = fmt.Sprintf("%.2f%%", 100.00)
		} else if len(result.LicenseGuesses) != 0 {
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

func toProgress(directory string, file string, rootLicenses []LicenseMatch, licenseGuesses []LicenseMatch, licenseIdentified []LicenseMatch) {
	license := ""
	confidence := ""

	if len(licenseIdentified) != 0 {
		license = joinLicenseList(licenseIdentified, " AND ")
		confidence = fmt.Sprintf("%.2f%%", 100.00)
	} else if len(licenseGuesses) != 0 {
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
		packageLicenseDeclared = joinLicenseList(fileResults[0].LicenseRoots, " OR ")
	}

	fmt.Println("SPDXVersion: SPDX-2.1")
	fmt.Println("DataLicense: CC0-1.0")
	fmt.Println("DocumentNamespace:http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82...") // TODO
	fmt.Println("DocumentName: DOCUMENTNAMEHEREFROMCLI")                                                         // TODO
	fmt.Println("SPDXID: SPDXRef-DOCUMENT")
	fmt.Println("Creator: Tool:", ToolName, ToolVersion)
	fmt.Println("Created:", time.Now().UTC().Format(time.RFC3339))
	fmt.Println("LicenseListVersion: 3.0")

	fmt.Println("")
	fmt.Println("PackageName: TODO") // TODO pass in from command line
	fmt.Println("SPDXID: SPDXRef-Package")
	fmt.Println("PackageDownloadLocation: NONE")
	fmt.Println("FilesAnalyzed: true")
	fmt.Println("PackageVerificationCode: 8b0600e4db514d62d9e2f10945f9c63488db9965") // TODO https://spdx.org/spdx-specification-21-web-version#h.2p2csry
	fmt.Println("PackageLicenseDeclared:", packageLicenseDeclared)
	fmt.Println("PackageLicenseConcluded:", packageLicenseDeclared)

	duplicateLicenseMatch := []LicenseMatch{}
	for _, result := range fileResults {
		if len(result.LicenseIdentified) != 0 {
			for _, license := range result.LicenseIdentified {
				duplicateLicenseMatch = append(duplicateLicenseMatch, license)
			}
		}
	}
	if len(duplicateLicenseMatch) != 0 {
		for _, license := range uniqLicenseMatch(duplicateLicenseMatch) {
			fmt.Println("PackageLicenseInfoFromFiles:", license.LicenseId)
		}
	} else {
		fmt.Println("PackageLicenseInfoFromFiles: NONE")
	}

	fmt.Println("PackageCopyrightText: NOASSERTION")
	fmt.Println("")

	// Loop over all files and get a list of all unique licenses and print below
	// PackageLicenseInfoFromFiles: GPL-2.0

	for _, result := range fileResults {

		// TODO this needs to possibly be NOASSERTION if unsure
		licenseConcluded := "NONE"

		if len(result.LicenseIdentified) != 0 {
			licenseConcluded = joinLicenseList(result.LicenseIdentified, " AND ")
		} else if len(result.LicenseGuesses) != 0 {
			licenseConcluded = result.LicenseGuesses[0].LicenseId
		}

		filePath := filepath.Join(result.Directory, result.Filename)
		if strings.HasPrefix(filePath, "./") == false {
			filePath = "./" + filePath
		}

		fmt.Println("FileName:", filePath)
		fmt.Println("SPDXID: SPDXRef-" + getSha1Hash([]byte(filePath)))
		fmt.Println("FileType: OTHER")
		fmt.Println("FileChecksum: SHA1:", result.Sha1Hash)
		fmt.Println("FileChecksum: SHA256:", result.Sha256Hash)
		fmt.Println("FileChecksum: MD5:", result.Md5Hash)
		fmt.Println("LicenseConcluded:", licenseConcluded)

		// FileComment: <text>The concluded license was taken from the package level that the file was included in.
		// This information was found in the COPYING.txt file in the xyz directory.</text>

		if len(result.LicenseIdentified) != 0 {
			for _, license := range result.LicenseIdentified {
				fmt.Println("LicenseInfoInFile:", license.LicenseId)
			}
		} else {
			fmt.Println("LicenseInfoInFile: NONE")
		}

		// fmt.Println("LicenseComments: The concluded license was taken from the package level that the file was included in")

		fmt.Println("FileCopyrightText: NOASSERTION")
		fmt.Println("")
	}
}
