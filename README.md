# README #

A reimplementation of http://www.boyter.org/2017/05/identify-software-licenses-python-vector-space-search-ngram-keywords/ using Go while I attempt to nut out the nuances of the language. Will this turn into a full blown SPDX parser of some description in the future? That depends on how sucessful I am writing Go and how much I enjoy it :)

Would be nice to have the following output styles,

	formatted IE CLI viewable
	JSON
	CSV

command line options include, deep guess, confidence, path, database file, short licence name to classify everything under by default

Probably want to have something that allows you to specify which extensions to look for explicitly to speed things up