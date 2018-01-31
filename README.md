# README #

A reimplementation of http://www.boyter.org/2017/05/identify-software-licenses-python-vector-space-search-ngram-keywords/ using Go while I attempt to nut out the nuances of the language. Will this turn into a full blown SPDX parser of some description in the future? That depends on how sucessful I am writing Go and how much I enjoy it :)

Would be nice to have the following output styles,

	formatted IE CLI viewable
	JSON
	CSV
	SPDX

command line options include, deep guess, confidence, path

Probably want to have something that allows you to specify which extensions to look for explicitly to speed things up
Having the SPDX version would be nice although currently just targetting 2.1

Example output running against itself ignoring the examples directory

	$ go run main.go --pbl .git,examples -f cli .
	Directory  File                      License       Confidence  Root Licenses  Size
	.          .gitignore                                          GPL-3.0-only   275B
	.          LICENSE                   GPL-3.0-only  99.68%      GPL-3.0-only   34.3K
	.          README.md                                           GPL-3.0-only   2.2K
	.          classifier_database.json                            GPL-3.0-only   158.5K
	.          database_keywords.json                              GPL-3.0-only   3.6M
	.          main.go                   GPL-3.0-only  100.00%     GPL-3.0-only   2.8K
	parsers    constants.go                                        GPL-3.0-only   5M
	parsers    formatter.go              GPL-3.0-only  100.00%     GPL-3.0-only   3.1K
	parsers    guesser.go                GPL-3.0-only  100.00%     GPL-3.0-only   9.1K
	parsers    helpers.go                GPL-3.0-only  100.00%     GPL-3.0-only   786B
	parsers    structs.go                GPL-3.0-only  100.00%     GPL-3.0-only   692B
	scripts    build_database.py         GPL-3.0-only  100.00%     GPL-3.0-only   4.7K
	scripts    include.go                GPL-3.0-only  100.00%     GPL-3.0-only   1008B