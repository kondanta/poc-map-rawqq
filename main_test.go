package main

import (
	"testing"
)

func TestSearchCount(t *testing.T) {
	res := Searchmanga("chiyu maho")

	if len(*res) != 2 {
		t.Errorf("Count of the manga is incorrect. got: %d, want: %d", len(*res), 2)
	}
}
