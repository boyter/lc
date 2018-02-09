licensechecker (lc)
-------------------
`lc` is a command line tool that recursively iterates over a supplied directory
attempting to identify what software license each file is under using the list
of licenses supplied by the SPDX (Software Package Data Exchange) Project. It will pick up 
license files named appropiateld or inline licenses such as the below in source files

`SPDX-License-Identifier: GPL-3.0-only`

In a nutshell this project is a reimplementation of http://www.boyter.org/2017/05/identify-software-licenses-python-vector-space-search-ngram-keywords/ using Go while I attempt to nut out the nuances of the language. 

It can produce report outputs as valid [SPDX](https://spdx.org/), CSV, JSON and CLI formatted. It has been designed to work inside CI systems that capture either stdout or file artifacts.

[![Build Status](https://travis-ci.org/boyter/lc.svg?branch=master)](https://travis-ci.org/boyter/lc)

### Why

Why should you care about what licenses your code runs under? See http://www.openlogic.com/resources/enterprise-blog/archive/use-spdx-for-open-source-license-compliance https://thenewstack.io/spdx-open-source-cheap-compliance-license-can-expensive/ https://www.infoworld.com/article/2839560/open-source-software/sticking-a-license-on-everything.html

### Installation

The binary name for `licencechecker` is `lc`.

For binary files see releases. To build from source you need to have Go setup with your GOPATH working and your go binary path exported like so,

```
export PATH=$PATH:$(go env GOPATH)/bin
```

then to install

```
$ go install
```


### Usage

Command line usage of `licensechecker` is designed to be as simple as possible.
Full details can be found in `lc --help`.

Probably the most useful functionality is the `-f` modifier which specifies the output format.
By default `licencechecker` will print out results as it processes files. However as it was designed
to run at the end of CI tasks you may want to get a nicer output which can be done like so.

```
$ lc -f tabular .
```

The above will process starting in the current directory and print out a formatted list of results when finished.

To view all command line options

```
$ lc --help
```

Example output of `licencechecker` running against itself in tabular format while ignoring the .git, licenses and vendor directories

```
$ lc -pbl .git,vendor,licenses -f tabular .
Directory            File                                                  License                        Confidence  Size
.                    .gitignore                                            GPL-3.0-only                   100.00%     275B
.                    .travis.yml                                           GPL-3.0-only                   100.00%     26B
.                    Gopkg.lock                                            GPL-3.0-only                   100.00%     1.4K
.                    Gopkg.toml                                            GPL-3.0-only                   100.00%     972B
.                    LICENSE                                               GPL-3.0-only                   99.68%      34.3K
.                    README.md                                             GPL-3.0-only                   100.00%     6.6K
.                    classifier_database.json                              GPL-3.0-only                   100.00%     158.5K
.                    database_keywords.json                                GPL-3.0-only                   100.00%     3.6M
.                    lc.exe                                                GPL-3.0-only                   100.00%     9M
.                    main.go                                               GPL-3.0-only                   100.00%     3.4K
.                    spdx-tools-2.1.12-SNAPSHOT-jar-with-dependencies.jar  GPL-3.0-only                   100.00%     43.1M
.                    spdx_example.spdx                                     GPL-3.0-only                   100.00%     9.3K
examples/identifier  LICENSE                                               GPL-3.0+ AND MIT               95.40%      1K
examples/identifier  LICENSE2                                              MIT AND GPL-3.0+               99.65%      35K
examples/identifier  has_identifier.py                                     (MIT OR GPL-3.0+) AND GPL-2.0  100.00%     428B
parsers              constants.go                                          GPL-3.0-only                   100.00%     5M
parsers              formatter.go                                          GPL-3.0-only                   100.00%     8.5K
parsers              formatter_test.go                                     GPL-3.0-only                   100.00%     976B
parsers              guesser.go                                            GPL-3.0-only                   100.00%     9.7K
parsers              guesser_test.go                                       GPL-3.0-only                   100.00%     610B
parsers              helpers.go                                            GPL-3.0-only AND Apache-2.0    100.00%     2.5K
parsers              structs.go                                            GPL-3.0-only                   100.00%     863B
scripts              build_database.py                                     GPL-3.0-only                   100.00%     4.7K
scripts              include.go                                            GPL-3.0-only                   100.00%     1008B
```

Or to write out the results to a CSV file

```
$ lc --format csv -output licences.csv --pathblacklist .git,licenses,vendor .
```

Or to a SPDX 2.1 file

```
$lc -f spdx -o spdx_example.spdx --pbl .git,vendor,licenses -dn licensechecker -pn licensechecker .
```


### SPDX

Running against itself to produce a SPDX file using tools from https://github.com/spdx/tools

```
$ go run main.go  -f spdx -o spdx_example.spdx --pbl .git,vendor,licenses -dn licensechecker -pn licensechecker . && java -jar ./spdx-tools-2.1.12-SNAPSHOT-jar-with-dependencies.jar Verify ./spdx_example.spdx
ERROR StatusLogger No log4j2 configuration file found. Using default configuration: logging only errors to the console. Set system property 'log4j2.debug' to show Log4j2 internal initialization logging.
03:49:29.479 [main] ERROR org.apache.jena.rdf.model.impl.RDFReaderFImpl - Rewired RDFReaderFImpl - configuration changes have no effect on reading
03:49:29.482 [main] ERROR org.apache.jena.rdf.model.impl.RDFReaderFImpl - Rewired RDFReaderFImpl - configuration changes have no effect on reading
This SPDX Document is valid.
```

### Package

Run go build for windows and linux then the following in linux

```
zip -r9 lc-1.0.0-x86_64-pc-windows.zip lc.exe && zip -r9 lc-1.0.0-x86_64-unknown-linux.zip lc
```


### TODO

Add error handling for all the file operations and just in general. Most are currently ignored

Add logic to guess the file type for SPDX value FileType

Add addtional unit and integration tests

Investigate using "github.com/gosuri/uitable" for formatting https://github.com/gosuri/uitable

Investigate using zlib compression for databases as per the below

```
package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
)

func readFile(filepath string) []byte {
	// TODO only read as deep into the file as we need
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Print(err)
	}

	return bytes
}

func main() {

	contents := readFile("database_keywords.json")
	fmt.Println(len(contents))

	var in bytes.Buffer
	b := []byte(contents)
	w := zlib.NewWriter(&in)
	w.Write(b)
	w.Close()

	fmt.Println(len(in.String()))

	var out bytes.Buffer
	r, _ := zlib.NewReader(&in)
	io.Copy(&out, r)
	fmt.Println(len(out.String()))
	// fmt.Println(len(out.String()))
}
```
