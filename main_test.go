package main

import (
	"strings"
	"testing"
)

func TestExtractURLs(t *testing.T) {
	examples := map[string][]string{
		"string without any apparent url":                            []string{},
		"http://example.com":                                         []string{"http://example.com"},
		"a url https://example.com in a sentance":                    []string{"https://example.com"},
		"http://xkcd.com twice in one sentence? http://www.xkcd.com": []string{"http://xkcd.com", "http://www.xkcd.com"},
		"ftp://valid.uri.tld will not be extracted":                  []string{},
		"http:// blank hosts are also not extracted":                 []string{},
		// TODO: want to extract these
		"a URN like xkcd.com":                                []string{},
		`embedded URLs like <a href="http://xkcd.com/974/">`: []string{},
	}

	for source, expected := range examples {
		actual := extractURLs(strings.NewReader(source))

		if !stringSlicesEqual(actual, expected) {
			t.Errorf("Expected %v from '%v', but actually got %v", expected, source, actual)
		}
	}
}

func stringSlicesEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && len(b) == 0) || (b == nil && len(a) == 0) {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
