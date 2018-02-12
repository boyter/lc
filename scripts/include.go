// SPDX-License-Identifier: MIT

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func readFile(filepath string) []byte {
	// TODO only read as deep into the file as we need
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Print(err)
	}

	return bytes
}

// Reads all .json files in the current folder
// and encodes them as strings literals in constants.go
func main() {
	files, _ := ioutil.ReadDir(".")
	out, _ := os.Create("./parsers/constants.go")

	// Open constants
	out.Write([]byte("package parsers \n\nconst (\n"))

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			// The constant variable name
			out.Write([]byte(strings.TrimSuffix(f.Name(), ".json") + " = `"))

			contents := readFile(f.Name())
			str := base64.StdEncoding.EncodeToString(contents)

			out.Write([]byte(str))
			out.Write([]byte("`\n"))
		}
	}

	// Close out constants
	out.Write([]byte(")\n"))
}
