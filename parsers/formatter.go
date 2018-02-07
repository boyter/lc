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
		licenseConcluded, confidence := determineLicense(result)

		rootLicenseString := ""
		for _, v := range result.LicenseRoots {
			rootLicenseString += fmt.Sprintf("%s,", v.LicenseId)
		}
		rootLicenseString = strings.TrimRight(rootLicenseString, ", ")

		records = append(records, []string{
			result.Filename,
			result.Directory,
			licenseConcluded,
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

func joinLicenseList(licenseList []LicenseMatch, ignore []LicenseMatch, operator string) string {
	licenseDeclared := ""

	if len(licenseList) == 1 {
		if licenceListHasLicense(licenseList[0], ignore) == false {
			licenseDeclared = licenseList[0].LicenseId
		}
	} else if len(licenseList) >= 2 {
		licenseNames := []string{}
		for _, v := range licenseList {
			if licenceListHasLicense(v, ignore) == false {
				licenseNames = append(licenseNames, v.LicenseId)
			}
		}

		if len(licenseNames) == 1 {
			licenseDeclared = licenseNames[0]
		} else if len(licenseNames) != 0 {

			licenseDeclared = strings.Join(licenseNames, operator)

			if operator == " OR " {
				licenseDeclared = "(" + licenseDeclared + ")"
			}
		}
	}

	return licenseDeclared
}

func determineLicense(result FileResult) (string, string) {
	license := ""
	confidence := 100.00
	licenseMatches := []LicenseMatch{}

	if len(result.LicenseIdentified) != 0 {
		license = joinLicenseList(result.LicenseIdentified, result.LicenseRoots, " AND ")
		confidence = 100.00
	} else if len(result.LicenseGuesses) != 0 {
		license = result.LicenseGuesses[0].LicenseId
		confidence = result.LicenseGuesses[0].Percentage * 100
		licenseMatches = append(licenseMatches, result.LicenseGuesses[0])
	}

	rootLicenses := joinLicenseList(result.LicenseRoots, licenseMatches, " OR ")
	if rootLicenses != "" {
		if license == "" {
			license = rootLicenses
		} else {
			license = rootLicenses + " AND " + license
		}
	}

	if license == "" {
		license = "NOASSERTION"
	}

	return license, fmt.Sprintf("%.2f%%", confidence)
}

func toTabular(fileResults []FileResult) {
	output := []string{
		"Directory | File | License | Confidence | Root Licenses | Size",
	}

	for _, result := range fileResults {
		license, confidence := determineLicense(result)

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
		license = joinLicenseList(licenseIdentified, []LicenseMatch{}, " AND ")
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

	lines := []string{}

	// Determine the package licenses
	packageLicenseDeclared := "NONE"

	if len(fileResults) != 0 {
		packageLicenseDeclared = joinLicenseList(fileResults[0].LicenseRoots, []LicenseMatch{}, " OR ")
	}

	lines = append(lines, "SPDXVersion: SPDX-2.1")
	lines = append(lines, "DataLicense: CC0-1.0")
	lines = append(lines, "DocumentNamespace:http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82...") // TODO
	lines = append(lines, "DocumentName: DOCUMENTNAMEHEREFROMCLI")                                                         // TODO
	lines = append(lines, "SPDXID: SPDXRef-DOCUMENT")
	lines = append(lines, "Creator: Tool: "+ToolName+ToolVersion)
	lines = append(lines, "Created: "+time.Now().UTC().Format(time.RFC3339))
	lines = append(lines, "LicenseListVersion: 3.0")

	lines = append(lines, "")
	lines = append(lines, "PackageName: TODO") // TODO pass in from command line
	lines = append(lines, "SPDXID: SPDXRef-Package")
	lines = append(lines, "PackageDownloadLocation: NONE")
	lines = append(lines, "FilesAnalyzed: true")
	lines = append(lines, "PackageVerificationCode: 8b0600e4db514d62d9e2f10945f9c63488db9965") // TODO https://spdx.org/spdx-specification-21-web-version#h.2p2csry
	lines = append(lines, "PackageLicenseDeclared: "+packageLicenseDeclared)
	lines = append(lines, "PackageLicenseConcluded: "+packageLicenseDeclared)

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
			lines = append(lines, "PackageLicenseInfoFromFiles: "+license.LicenseId)
		}
	} else {
		lines = append(lines, "PackageLicenseInfoFromFiles: NONE")
	}

	lines = append(lines, "PackageCopyrightText: NOASSERTION")
	lines = append(lines, "")

	// Loop over all files and get a list of all unique licenses and print below
	// PackageLicenseInfoFromFiles: GPL-2.0

	for _, result := range fileResults {
		licenseConcluded, _ := determineLicense(result)

		filePath := filepath.Join(result.Directory, result.Filename)
		if strings.HasPrefix(filePath, "./") == false {
			filePath = "./" + filePath
		}

		lines = append(lines, "FileName: "+filePath)
		lines = append(lines, "SPDXID: SPDXRef-"+getSha1Hash([]byte(filePath)))
		lines = append(lines, "FileType: OTHER")
		lines = append(lines, "FileChecksum: SHA1: "+result.Sha1Hash)
		lines = append(lines, "FileChecksum: SHA256: "+result.Sha256Hash)
		lines = append(lines, "FileChecksum: MD5: "+result.Md5Hash)
		lines = append(lines, "LicenseConcluded: "+licenseConcluded)

		// FileComment: <text>The concluded license was taken from the package level that the file was included in.
		// This information was found in the COPYING.txt file in the xyz directory.</text>

		if len(result.LicenseIdentified) != 0 {
			for _, license := range result.LicenseIdentified {
				lines = append(lines, "LicenseInfoInFile:"+license.LicenseId)
			}
		} else {
			lines = append(lines, "LicenseInfoInFile: NONE")
		}

		// fmt.Println("FileComment: The concluded license was taken from the package level that the file was included in")

		lines = append(lines, "FileCopyrightText: NOASSERTION")
		lines = append(lines, "")
	}

	if FileOutput == "" {
		for _, line := range lines {
			fmt.Println(line)
		}
	} else {
		ioutil.WriteFile(FileOutput, []byte(strings.Join(lines, "\n")), 0600)
	}
}
