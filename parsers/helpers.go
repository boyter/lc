// SPDX-License-Identifier: GPL-3.0-only

package parsers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func getMd5Hash(content []byte) string {
	hasher := md5.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

func getSha1Hash(content []byte) string {
	hasher := sha1.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

func getSha256Hash(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

func readFile(filepath string) []byte {
	// TODO only read as deep into the file as we need
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Print(err)
	}

	return bytes
}
