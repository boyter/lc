licensechecker (lc)
-------------------
`lc` is a command line tool that recursively iterates over a supplied directory
attempting to identify what software license each file is under using the list
of licenses supplied by the SPDX (Software Package Data Exchange) Project. It will pick up 
license files named appropiateld or inline licenses such as the below in source files

`SPDX-License-Identifier: GPL-3.0-only`

In a nutshell this project is a reimplementation of http://www.boyter.org/2017/05/identify-software-licenses-python-vector-space-search-ngram-keywords/ using Go while I attempt to nut out the nuances of the language. 

It can produce report outputs as valid SDPX, CSV, and CLI formatted. It has been designed to work inside CI systems that capture either stdout or file artifacts.

### Installation

The binary name for `licencechecker` is `lc`.

Binary files will be distributed at some point in the future probably. Currently to install you need to have Go setup with your GOPATH working and your go binary path exported like so,

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
Directory            File                                                  License                        Confidence  Root Licenses  Size
.                    .gitignore                                                                                       GPL-3.0-only   275B
.                    Gopkg.lock                                                                                       GPL-3.0-only   1.4K
.                    Gopkg.toml                                                                                       GPL-3.0-only   972B
.                    LICENSE                                               GPL-3.0-only                   99.68%      GPL-3.0-only   34.3K
.                    README.md                                                                                        GPL-3.0-only   7.1K
.                    classifier_database.json                                                                         GPL-3.0-only   158.5K
.                    database_keywords.json                                                                           GPL-3.0-only   3.6M
.                    main.go                                               GPL-3.0-only                   100.00%     GPL-3.0-only   2.8K
examples/identifier  LICENSE                                               MIT                            95.40%      MIT            1K
examples/identifier  has_identifier.py                                     (GPL-2.0 AND GPL-3.0+)         100.00%     MIT            428B
parsers              constants.go                                                                                     GPL-3.0-only   5M
parsers              formatter.go                                          GPL-3.0-only                   100.00%     GPL-3.0-only   7.5K
parsers              guesser.go                                            GPL-3.0-only                   100.00%     GPL-3.0-only   9.5K
parsers              helpers.go                                            (GPL-3.0-only AND Apache-2.0)  100.00%     GPL-3.0-only   2.2K
parsers              structs.go                                            GPL-3.0-only                   100.00%     GPL-3.0-only   754B
scripts              build_database.py                                     GPL-3.0-only                   100.00%     GPL-3.0-only   4.7K
scripts              include.go                                            GPL-3.0-only                   100.00%     GPL-3.0-only   1008B
```

Or to write out the results to a CSV file

```
$ lc --format csv -output licences.csv --pathblacklist .git,licenses,vendor .
```


### SPDX

Running against itself to produce a SPDX file.

```
$ go run main.go --pbl .git,examples,vendor -f spdx . > ./spdx_example.spdx && java -jar ./spdx-tools-2.1.12-SNAPSHOT-jar-with-dependencies.jar Verify ./spdx_example.spdx
ERROR StatusLogger No log4j2 configuration file found. Using default configuration: logging only errors to the console. Set system property 'log4j2.debug' to show Log4j2 internal initialization logging.
03:49:29.479 [main] ERROR org.apache.jena.rdf.model.impl.RDFReaderFImpl - Rewired RDFReaderFImpl - configuration changes have no effect on reading
03:49:29.482 [main] ERROR org.apache.jena.rdf.model.impl.RDFReaderFImpl - Rewired RDFReaderFImpl - configuration changes have no effect on reading
This SPDX Document is valid.
```

### TODO

Add error handling for all the file operations and just in general. Most are currently ignored

Add unit and integration tests

Investigate using "github.com/gosuri/uitable" for formatting https://github.com/gosuri/uitable

Investigate using zlib compression for databases

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

For checking the spdx validity use the following built from https://github.com/spdx/tools

```
java -jar ./spdx-tools-2.1.12-SNAPSHOT-jar-with-dependencies.jar Verify spdx_example.spdx
java -jar ./spdx-tools-2.1.12-SNAPSHOT-jar-with-dependencies.jar TagToRDF spdx_example.spdx TagToRDF.rdf

lc --pbl .git,examples,vendor -f spdx . > ~/Projects/tools/target/spdx_example.spdx
java -jar ./spdx-tools-2.1.12-SNAPSHOT-jar-with-dependencies.jar Verify ./spdx_example.spdx

go run main.go --pbl .git,examples,vendor -f spdx . > ./spdx_example.spdx && java -jar ./spdx-tools-2.1.12-SNAPSHOT-jar-with-dependencies.jar Verify ./spdx_example.spdx
```

