// SPDX-License-Identifier: MIT OR Unlicense

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func readFile(filepath string) []byte {
	// TODO only read as deep into the file as we need
	bytes, err := os.ReadFile(filepath)

	if err != nil {
		fmt.Print(err)
	}

	return bytes
}

// Reads all .json files in the current folder
// and encodes them as strings literals in constants.go
func main() {
	fmt.Println("running...")
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println(err.Error())
	}
	out, err := os.Create("./processor/constants.go")
	if err != nil {
		fmt.Println(err.Error())
	}

	// Open constants
	out.Write([]byte("package processor \n\nconst (\n"))

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
	_, err = out.Write([]byte(")\n"))
	if err != nil {
		fmt.Println(err.Error())
	}
	err = out.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
}
