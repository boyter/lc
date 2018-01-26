package parsers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ryanuber/columnize"
	"io/ioutil"
	"log"
	"os"
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
			license = result.LicenseGuesses[0].Shortname
			confidence = fmt.Sprintf("%.3f", result.LicenseGuesses[0].Percentage*100)
		}

		rootLicenseString := ""
		for _, v := range result.LicenseRoots {
			rootLicenseString += fmt.Sprintf("%s,", v.Shortname)
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
			license = result.LicenseGuesses[0].Shortname
			confidence = fmt.Sprintf("%.2f%%", result.LicenseGuesses[0].Percentage*100)
		}

		rootLicenseString := ""
		for _, v := range result.LicenseRoots {
			rootLicenseString += fmt.Sprintf("%s,", v.Shortname)
		}
		rootLicenseString = strings.TrimRight(rootLicenseString, ", ")

		output = append(output, fmt.Sprintf("%s | %s | %s | %s | %s | %s", result.Directory, result.Filename, license, confidence, rootLicenseString, result.BytesHuman))
	}

	result := columnize.SimpleFormat(output)

	fmt.Println(result)
}
