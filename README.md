licensechecker (lc)
-------------------

# NOTE - this is under heavy development, and as such master does not currently work, see a release for a working solution!

`lc` is a command line tool that recursively iterates over a supplied directory or file 
attempting to identify what software license each file is under using the list
of licenses supplied by the SPDX (Software Package Data Exchange) Project. It will pick up 
license files named appropriately or inline licenses such as the below in source files

`SPDX-License-Identifier: GPL-3.0-only`

In a nutshell this project is a reimplementation of http://www.boyter.org/2017/05/identify-software-licenses-python-vector-space-search-ngram-keywords/ using Go while I attempt to nut out the nuances of the language. 

It can produce report outputs as valid [SPDX](https://spdx.org/), CSV, XLSX, JSON and CLI formatted. It has been designed to work inside CI systems that capture either stdout or file artifacts.

[![Go](https://github.com/boyter/lc/actions/workflows/go.yml/badge.svg)](https://github.com/boyter/lc/actions/workflows/go.yml)
[![Scc Count Badge](https://sloc.xyz/github/boyter/lc/)](https://github.com/boyter/lc/)

Dual-licensed under MIT or the [UNLICENSE](http://unlicense.org).

### Why

In short taken from, http://ben.balter.com/licensee/

 * You've got an open source project. How do you know what you can and can't do with the software?
 * You've got a bunch of open source projects, how do you know what their licenses are?
 * You've got a project with a license file, but which license is it? Has it been modified?

Why should you care about what licenses your code runs under? See 

 * http://www.openlogic.com/resources/enterprise-blog/archive/use-spdx-for-open-source-license-compliance 
 * https://thenewstack.io/spdx-open-source-cheap-compliance-license-can-expensive/
 * https://www.infoworld.com/article/2839560/open-source-software/sticking-a-license-on-everything.html

### Installation

The binary name for `licencechecker` is `lc`.

For binary files see releases https://github.com/boyter/lc/releases To build from source you need to have Go setup with your GOPATH working and your go binary path exported like so,

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

```
$ lc --help
NAME:
   licensechecker - Check directory for licenses and list what license(s) a file is under

USAGE:
   lc [global options] [DIRECTORY|FILE] [DIRECTORY|FILE]

VERSION:
   1.3.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --format csv, -f csv                                Set output format, supports progress, tabular, json, spdx, summary, xlsx or csv (default: "tabular")
   --output FILE, -o FILE                              Set output file if not set will print to stdout FILE
   --confidence 0.95, -c 0.95                          Set required confidence level for licence matching between 0 and 1 E.G. 0.95 (default: "0.85")
   --deepguess true, --dg true                         Should attempt to deep guess the licence false or true true (default: "true")
   --filesize 50000, --fs 50000                        How large a file in bytes should be processed 50000 (default: "50000")
   --licensefiles copying,readme, --lf copying,readme  Possible license files to inspect for over-arching license as comma seperated list copying,readme (default: "license,licence,copying,readme")
   --pathblacklist .git,.hg,.svn, --pbl .git,.hg,.svn  Which directories should be ignored as comma seperated list .git,.hg,.svn (default: ".git,.hg,.svn")
   --extblacklist gif,jpg,png, --xbl gif,jpg,png       Which file extensions should be ignored for deep analysis as comma seperated list E.G. gif,jpg,png (default: "woff,eot,cur,dm,xpm,emz,db,scc,idx,
mpp,dot,pspimage,stl,dml,wmf,rvm,resources,tlb,docx,doc,xls,xlsx,ppt,pptx,msg,vsd,chm,fm,book,dgn,blines,cab,lib,obj,jar,pdb,dll,bin,out,elf,so,msi,nupkg,pyc,ttf,woff2,jpg,jpeg,png,gif,bmp,psd,tif,tif
f,yuv,ico,xls,xlsx,pdb,pdf,apk,com,exe,bz2,7z,tgz,rar,gz,zip,zipx,tar,rpm,bin,dmg,iso,vcd,mp3,flac,wma,wav,mid,m4a,3gp,flv,mov,mp4,mpg,rm,wmv,avi,m4v,sqlite,class,rlib,ncb,suo,opt,o,os,pch,pbm,pnm,ppm
,pyd,pyo,raw,uyv,uyvy,xlsm,swf")
   --documentname LicenseChecker, --dn LicenseChecker  SPDX only. Sets DocumentName E.G. LicenseChecker (default: "Unknown")
   --packagename LicenseChecker, --pn LicenseChecker   SPDX only. Sets PackageName E.G. LicenseChecker (default: "Unknown")
   --documentnamespace value, --dns value              SPDX only. Sets DocumentNamespace, if not set will default to http://spdx.org/spdxdocs/[packagename]-[HASH]
   --help, -h                                          show help
   --version, -v                                       print the version
```

More information about [what licensechecker looks at and how it works](what-we-look-at.md)

Probably the most useful functionality is the `-f` modifier which specifies the output format.
By default `licencechecker` will print out results in a tabular CLI format. However as it was designed
to run at the end of CI tasks you may want to change it. This can be done like so.

```
$ lc -f tabular .
$ lc -f progress .
$ lc -f spdx .
$ lc -f csv .
$ lc -f summary .
```

The above will process starting in the current directory and print out a formatted list of results to the CLI when finished.

Example output of `licencechecker` running against itself in tabular format while ignoring the .git, licenses and vendor directories

```
$ lc -pbl .git,vendor,licenses -f tabular .
-----------------------------------------------------------------------------------------------------------
Directory            File                    License                                      Confidence  Size
-----------------------------------------------------------------------------------------------------------
.                    .gitignore              (MIT OR Unlicense)                           100.00%     278B
.                    .travis.yml             (MIT OR Unlicense)                           100.00%     192B
.                    CODE_OF_CONDUCT.md      (MIT OR Unlicense)                           100.00%     3.1K
.                    CONTRIBUTING.md         (MIT OR Unlicense)                           100.00%     1.2K
.                    Gopkg.lock              (MIT OR Unlicense)                           100.00%     1.4K
.                    Gopkg.toml              (MIT OR Unlicense)                           100.00%     972B
.                    LICENSE                 Unlicense AND MIT                            94.83%      1.1K
.                    README.md               (MIT OR Unlicense)                           100.00%     10.6K
.                    UNLICENSE               MIT AND Unlicense                            95.16%      1.2K
.                    database_keywords.json  (MIT OR Unlicense)                           100.00%     3.6M
.                    licensechecker.spdx     (MIT OR Unlicense)                           100.00%     9.3K
.                    main.go                 (MIT OR Unlicense)                           100.00%     3.4K
.                    what-we-look-at.md      (MIT OR Unlicense)                           100.00%     3.7K
examples/identifier  LICENSE                 GPL-3.0+ AND MIT                             95.40%      1K
examples/identifier  LICENSE2                MIT AND GPL-3.0+                             99.65%      35K
examples/identifier  has_identifier.py       (MIT OR GPL-3.0+) AND GPL-2.0                100.00%     409B
parsers              constants.go            (MIT OR Unlicense)                           100.00%     4.8M
parsers              formatter.go            (MIT OR Unlicense)                           100.00%     8.5K
parsers              formatter_test.go       (MIT OR Unlicense)                           100.00%     1.3K
parsers              guesser.go              (MIT OR Unlicense)                           100.00%     9.8K
parsers              guesser_test.go         (MIT OR Unlicense) AND GPL-2.0 AND GPL-3.0+  100.00%     4.8K
parsers              helpers.go              (MIT OR Unlicense) AND Apache-2.0            100.00%     2.4K
parsers              helpers_test.go         (MIT OR Unlicense)                           100.00%     2.8K
parsers              structs.go              (MIT OR Unlicense)                           100.00%     679B
scripts              build_database.py       (MIT OR Unlicense)                           100.00%     4.6K
scripts              include.go              (MIT OR Unlicense)                           100.00%     951B
-----------------------------------------------------------------------------------------------------------
```

To write out the results to a CSV file

```
$ lc --format csv -output licences.csv --pathblacklist .git,licenses,vendor .
```

Or to a SPDX 2.1 file

```
$ lc -f spdx -o licensechecker.spdx --pbl .git,vendor,licenses -dn licensechecker -pn licensechecker .
```

You can specify multiple directories as additional arguments and all results will be merged into a single output

```
$ lc -f tabular ./examples/identifier ./scripts
```

You can also specify files and directories as additional arguments 

```
$ lc -f tabular README.md LICENSE ./examples/identifier
------------------------------------------------------------------------------------------
Directory              File               License                        Confidence  Size
------------------------------------------------------------------------------------------
                       README.md          NOASSERTION                    100.00%     11.3K
                       LICENSE            MIT                            94.83%      1.1K
./examples/identifier  LICENSE            GPL-3.0+ AND MIT               95.40%      1K
./examples/identifier  LICENSE2           MIT AND GPL-3.0+               99.65%      35K
./examples/identifier  has_identifier.py  (MIT OR GPL-3.0+) AND GPL-2.0  100.00%     409B
------------------------------------------------------------------------------------------
```

### SPDX

The ouput of SPDX is a valid SPDX 2.1 document. Validation was checked against the tools supplied by the SPDX group.
Running master against itself to produce a SPDX and the validating using the tools from https://github.com/spdx/tools

```
$ go run main.go  -f spdx -o spdx_example.spdx --pbl .git,vendor,licenses -dn licensechecker -pn licensechecker . && java -jar ./spdx-tools-2.1.12-SNAPSHOT-jar-with-dependencies.jar Verify ./spdx_example.spdx
ERROR StatusLogger No log4j2 configuration file found. Using default configuration: logging only errors to the console. Set system property 'log4j2.debug' to show Log4j2 internal initialization logging.
03:49:29.479 [main] ERROR org.apache.jena.rdf.model.impl.RDFReaderFImpl - Rewired RDFReaderFImpl - configuration changes have no effect on reading
03:49:29.482 [main] ERROR org.apache.jena.rdf.model.impl.RDFReaderFImpl - Rewired RDFReaderFImpl - configuration changes have no effect on reading
This SPDX Document is valid.
```

### Package

Run go build for windows and linux then the following in linux, keep in mind need to update the version

```
zip -r9 lc-1.0.0-x86_64-pc-windows.zip lc.exe && zip -r9 lc-1.0.0-x86_64-unknown-linux.zip lc

GOOS=darwin GOARCH=amd64 go build && zip -r9 lc-1.0.0-x86_64-apple-darwin.zip lc
GOOS=windows GOARCH=amd64 go build && zip -r9 lc-1.0.0-x86_64-pc-windows.zip lc.exe
GOOS=linux GOARCH=amd64 go build && zip -r9 lc-1.0.0-x86_64-unknown-linux.zip lc
```

### Most Common Software Licences

Source https://www.blackducksoftware.com/top-open-source-licenses

Source https://blog.sourced.tech/post/gld/pga-licenses.csv

```
Rank 	Open Source License 	                            %
1.      MIT License 	                                    38%
2.      GNU General Public License (GPL 2.0) 	            14%
3.      Apache License 2.0                                  13%
4.      ISC License 	                                    10%
5.      GNU General Public License (GNU) 3.0 	            6%
6.      BSD License 2.0 (3-clause, New or Revised) License  5%
7.      Artistic License (Perl)                             3%
8.      GNU Lesser General Public License (LGPL) 2.1 	    3%
9.      GNU Lesser General Public License (LGPL) 3.0 	    1%
10. 	Eclipse Public License (EPL) 	                    1%
11. 	Microsoft Public License                            1%
12. 	Simplified BSD License (BSD) 	                    1%
13. 	Code Project Open License 1.02 	                    < 1%
14. 	Mozilla Public License (MPL) 1.1                    < 1%
15. 	GNU Affero General Public License v3 or later 	    < 1%
16. 	Common Development and Distribution License (CDDL)  < 1%
17. 	DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE 	    < 1%
18. 	Microsoft Reciprocal License 	                    < 1%
19. 	Sun GPL with Classpath Exception v2.0 	            < 1%
20. 	zlib/libpng License 	                            < 1%
```

### TODO

* Add error handling for all the file operations and just in general. Most are currently ignored
* Add logic to guess the file type for SPDX value FileType
* Add addtional unit and integration tests
* Investigate using "github.com/gosuri/uitable" for formatting https://github.com/gosuri/uitable
* https://web.archive.org/web/20180822173147/https://blog.sourced.tech/post/gld/
* https://github.com/boyter/boyter.org/blob/01601a2cafc2b2788b29b6943ad45ad40316d9a8/content/posts/improving-lc-performance.md
* https://reuse.software/
