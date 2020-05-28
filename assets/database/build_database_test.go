package main

import (
	"testing"
)

func TestNgamOne(t *testing.T) {
	ngrams := findNgrams([]string{"1", "2", "3", "4"}, 1)

	if len(ngrams) != 4 {
		t.Errorf("should be 4")
	}
}

func TestNgamThree(t *testing.T) {
	ngrams := findNgrams([]string{"1", "2", "3", "4"}, 3)

	if len(ngrams) != 2 {
		t.Errorf("should be 2")
	}
}
