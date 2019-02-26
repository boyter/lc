package parsers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
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
		csvFile, _ := os.OpenFile(FileOutput, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		defer csvFile.Close()

		w := csv.NewWriter(csvFile)
		w.WriteAll(records) // calls Flush internally

		if err := w.Error(); err != nil {
			log.Fatalln("error writing csv:", err)
		}

		fmt.Println("Results written to " + FileOutput)
	}
}

func toXLSX(fileResults []FileResult) {
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

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		log.Fatalln("error creating xlsx:", err)
	}

	for _, record := range records {
		row := sheet.AddRow()
		written := row.WriteSlice(&record, -1)
		if written < 0 {
			log.Fatalln("error writing row")
		}
	}

	// As this is a binary format, writing to stdout is primarily for the purpose of enabling redirection
	if FileOutput == "" {
		err = file.Write(os.Stdout)
		if err != nil {
			log.Fatalln("error writing xlsx:", err)
		}
	} else {
		file.Save(FileOutput)
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
		var licenseNames []string
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
	var licenseMatches []LicenseMatch

	if len(result.LicenseIdentified) != 0 {
		license = joinLicenseList(result.LicenseIdentified, result.LicenseRoots, " AND ")
		confidence = 100.00
	} else if len(result.LicenseGuesses) != 0 {
		license = result.LicenseGuesses[0].LicenseId
		confidence = result.LicenseGuesses[0].Score
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

var tabularShortBreak = "-------------------------------------------------------------------------------\n"
var tabularShortFormatHead = "%-71s %7s\n"

func toTabular(results []FileResult) {
	var str strings.Builder

	str.WriteString(tabularShortBreak)
	str.WriteString(fmt.Sprintf(tabularShortFormatHead, "File", "License"))
	str.WriteString(tabularShortBreak)

	wasTrimmed := false

	for _, result := range results {
		license, _ := determineLicense(result)
		tmp := filepath.Join(result.Directory, result.Filename)

		if len(tmp)+len(license) >= 79 {
			toTrim := 77 - len(license)

			if toTrim <= 0 {
				toTrim = 0
			}

			tmp = "~" + tmp[len(tmp)-toTrim:]

			if len(license) > 78 {
				license = license[:76] + "~"
			}
			wasTrimmed = true

		} else {
			toPad := 78 - len(tmp) - len(license)

			for i := 0; i < toPad; i++ {
				tmp = tmp + " "
			}
		}

		str.WriteString(fmt.Sprintf("%s %s\n", tmp, license))
	}

	str.WriteString(tabularShortBreak)

	if wasTrimmed {
		str.WriteString("    Some licenses were too long for display. Please use CSV or JSON output.\n")
		str.WriteString(tabularShortBreak)
	}

	fmt.Println(str.String())
}

func generatePackageVerificationCode(fileResults []FileResult) string {
	// Based on https://github.com/spdx/tools-python/blob/a48022e65a8897d0e4f2e93d8e53695d2c13ea23/spdx/package.py#L233
	var hashes []string

	for _, result := range fileResults {
		hashes = append(hashes, result.Sha1Hash)
	}

	sort.Strings(hashes)
	return getSha1Hash([]byte(strings.Join(hashes, "")))
}

func generateDocumentNamespace() string {
	if DocumentNamespace == "" {
		return "http://spdx.org/spdxdocs/" + PackageName + "-" + getSha1Hash([]byte(time.Now().UTC().Format(time.RFC3339)))
	}

	return DocumentNamespace
}

func toSPDX21(fileResults []FileResult) {
	var lines []string

	packageLicenseDeclared := "NONE"
	if len(fileResults) != 0 {
		packageLicenseDeclared = joinLicenseList(fileResults[0].LicenseRoots, []LicenseMatch{}, " OR ")
	}

	lines = append(lines, "SPDXVersion: SPDX-2.1")
	lines = append(lines, "DataLicense: CC0-1.0")
	lines = append(lines, "DocumentNamespace: "+generateDocumentNamespace())
	lines = append(lines, "DocumentName: "+DocumentName)
	lines = append(lines, "SPDXID: SPDXRef-DOCUMENT")
	lines = append(lines, "Creator: Tool: "+ToolName+" "+ToolVersion)
	lines = append(lines, "Created: "+time.Now().UTC().Format(time.RFC3339))
	lines = append(lines, "LicenseListVersion: 3.0")

	lines = append(lines, "")
	lines = append(lines, "PackageName: "+PackageName)
	lines = append(lines, "SPDXID: SPDXRef-Package")
	lines = append(lines, "PackageDownloadLocation: NONE")
	lines = append(lines, "FilesAnalyzed: true")
	lines = append(lines, "PackageVerificationCode: "+generatePackageVerificationCode(fileResults))
	lines = append(lines, "PackageLicenseDeclared: "+packageLicenseDeclared)
	lines = append(lines, "PackageLicenseConcluded: "+packageLicenseDeclared)

	var duplicateLicenseMatch []LicenseMatch
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

	for _, result := range fileResults {
		licenseConcluded, _ := determineLicense(result)

		filePath := filepath.Join(result.Directory, result.Filename)
		if strings.HasPrefix(filePath, "./") == false {
			filePath = "./" + filePath
		}

		lines = append(lines, "FileName: "+filePath)
		lines = append(lines, "SPDXID: SPDXRef-"+getSha1Hash([]byte(filePath)))
		lines = append(lines, "FileChecksum: SHA1: "+result.Sha1Hash)
		lines = append(lines, "FileChecksum: SHA256: "+result.Sha256Hash)
		lines = append(lines, "FileChecksum: MD5: "+result.Md5Hash)
		lines = append(lines, "LicenseConcluded: "+licenseConcluded)

		if len(result.LicenseIdentified) != 0 {
			for _, license := range result.LicenseIdentified {
				lines = append(lines, "LicenseInfoInFile:"+license.LicenseId)
			}
		} else {
			lines = append(lines, "LicenseInfoInFile: NONE")
		}

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
