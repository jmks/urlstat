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
		"a blank host like http:// is ignored":                       []string{},
		"URNs like xkcd.com or mail.google.com":                      []string{"xkcd.com", "mail.google.com"},
		"local addresses like http://localhost:9000":                 []string{"http://localhost:9000"},
		"local addresses must have a scheme localhost:9000":          []string{},
		// TODO: want to extract these
		`embedded URLs like <a href="http://xkcd.com/974/"></a>`: []string{},
		// Don't extract URI-looking things without a known TLD
		"Class.new.method": []string{},
		"host.fakeTLD":     []string{},
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
