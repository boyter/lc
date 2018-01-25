package parsers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func toCSV(fileResults []FileResult) {
	records := [][]string{
		{"filename", "directory", "license", "confidence"},
	}

	for _, result := range fileResults {

		license := ""
		confidence := ""

		if len(result.LicenseGuesses) != 0 {
			license = result.LicenseGuesses[0].Shortname
			confidence = fmt.Sprintf("%.2f", result.LicenseGuesses[0].Percentage*100)
		}

		records = append(records, []string{result.Filename, result.Directory, license, confidence})
	}

	w := csv.NewWriter(os.Stdout)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
