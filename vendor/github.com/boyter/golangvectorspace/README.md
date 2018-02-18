# README #

An implementation of the VectorSpace model in Golang. Pass in two strings and get back a number indicating how similar they are.

[![Build Status](https://travis-ci.org/boyter/golangvectorspace.svg?branch=master)](https://travis-ci.org/boyter/golangvectorspace)

Usage like so,

```
var concordance1 = Concordance("Go has a lightweight test framework composed of the go test command and the testing package.")
var concordance2 = Concordance("Package testing provides support for automated testing of Go packages. It is intended to be used in concert with the go test command, which automates execution of any function of the form.")

// value of got will be 0.48211825989991874   
got := Relation(concordance1, concordance2)
```


See tests for other examples.

To benchmark,

```
go test -bench .
```

On a late 2013 Macbook Pro 2.6 GHz Intel Core i5

```
BenchmarkRelation-4   	 1000000	      2364 ns/op
BenchmarkRelation-4   	 1000000	      2396 ns/op
BenchmarkRelation-4   	 1000000	      2318 ns/op
```

Coverage

```
go test -cover .
ok  	github.com/boyter/golangvectorspace	0.006s	coverage: 95.7% of statements
```