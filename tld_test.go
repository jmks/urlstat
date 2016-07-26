package main

import (
	"testing"

	"github.com/jmks/urlstat/tld"
)

func TestIsValid(t *testing.T) {
	examples := map[string]bool{
		"com":         true,
		"ca":          true,
		"嘉里":          true,
		"notarealtld": false,
	}

	for tldStr, expected := range examples {
		actual := tld.IsValid(tldStr)
		if actual != expected {
			t.Errorf("Expected %v to be %v, but got %v", tldStr, expected, actual)
		}
	}
}
