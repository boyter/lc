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

Example output running against itself currently

	$ licensechecker -f cli licensechecker
	Directory                File                      License       Confidence  Root Licenses  Size
	licensechecker           .gitignore                                          GPL-3.0-only   275B
	licensechecker           LICENSE                   GPL-3.0-only  99.68%      GPL-3.0-only   34.3K
	licensechecker           README.md                                           GPL-3.0-only   718B
	licensechecker           classifier_database.json                            GPL-3.0-only   158.5K
	licensechecker           database_keywords.json                              GPL-3.0-only   3.6M
	licensechecker           main.go                   GPL-3.0       100.00%     GPL-3.0-only   2.8K
	licensechecker/examples  has_identifier.py         GPL-2.0       100.00%     GPL-3.0-only   422B
	licensechecker/parsers   constants.go                                        GPL-3.0-only   5M
	licensechecker/parsers   formatter.go              GPL-3.0       100.00%     GPL-3.0-only   3K
	licensechecker/parsers   guesser.go                GPL-3.0       100.00%     GPL-3.0-only   9.1K
	licensechecker/parsers   helpers.go                GPL-3.0       100.00%     GPL-3.0-only   781B
	licensechecker/parsers   structs.go                GPL-3.0       100.00%     GPL-3.0-only   687B
	licensechecker/scripts   build_database.py         GPL-3.0       100.00%     GPL-3.0-only   4.7K
	licensechecker/scripts   include.go                GPL-3.0       100.00%     GPL-3.0-only   1003B