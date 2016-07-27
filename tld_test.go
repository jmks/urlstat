package main

import (
	"testing"

	"github.com/jmks/urlstat/tld"
)

func TestHasKnownTLD(t *testing.T) {
	examples := map[string]bool{
		"subdomain.example.com":       true,
		"example.badtld":              false,
		"a.b.c.嘉里":                    true,
		"host.org/path/to/resource":   true,
		"host.net#fragment":           true,
		"example.badtld#withfragment": false,
		"example.com?test=true":       true,
	}

	for host, expected := range examples {
		actual := tld.HasKnownTLD(host)

		if actual != expected {
			t.Errorf("Expected %v to be %v, but got %v", host, expected, actual)
		}
	}
}
