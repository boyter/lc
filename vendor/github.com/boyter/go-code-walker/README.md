# go-code-walker

[![Go Report Card](https://goreportcard.com/badge/github.com/boyter/go-code-walker)](https://goreportcard.com/report/github.com/boyter/go-code-walker)
[![Str Count Badge](https://sloc.xyz/github/boyter/go-code-walker/)](https://github.com/boyter/go-code-walker/)

Library to help with walking of code directories in Go

Package provides file operations specific to code repositories such as walking the file tree obeying .ignore and .gitignore files
or looking for the root directory assuming already in a git project.

Note that it currently has dependancy on go-gitignore which is pulled in here to avoid external dependencies. This needs to be rewritten
as there are some bugs in that implementation.

All code is dual-licenced as either MIT or Unlicence.
Note that as an Australian I cannot put this into the public domain, hence the choice most liberal licences I can find.
