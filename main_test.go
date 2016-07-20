package main

import (
	"reflect"
	"testing"
)

func TestURIRegexMatches(t *testing.T) {
	data := map[string]bool{
		"withoutsubdomain.tld":                                             false,
		"subdomain.host.tld":                                               true,
		"subdomain.host.tld/path/to/resource":                              true,
		"subdomain.host.tld/hypenated-paths-are-cool-too":                  true,
		"subdomain.host.tld/path?key=value":                                true,
		"subdomain.host.tld#fragment":                                      true,
		"subdomain.host.tld/path#fragment":                                 true,
		"subdomain.host.tld/another/path?key1=value1&key2=value2#fragment": true,
		"https://subdomain.host.tld/path?timey=wimey&wibbly=wobbly#drwho":  true,
		// localhost is OK
		"localhost:1337":                       true,
		"https://localhost:1701/wibbly/wobbly": true,
		"http://localhost:1701#timeywimey":     true,
		"https://localhost:1701/doctor#who":    true,
		// ignore other schemes
		"ftp://some-ftp-host":   false,
		"mailto:eg@example.com": false,
	}

	for url, match := range data {
		if match != uriRegex.MatchString(url) {
			if match {
				t.Errorf("Expected to match %v but did not\n", url)
			} else {
				t.Errorf("Matched %v when it should not have\n", url)
			}
		}
	}
}

func TestURIRegexExtract(t *testing.T) {
	data := map[string][]string{
		"blah blah www.jmks.ca is awesome":                                                        []string{"www.jmks.ca"},
		"Multiple URLs like mail.google.com and http://www.calmingmanatee.com/1 are both matched": []string{"mail.google.com", "http://www.calmingmanatee.com/1"},
	}

	for haystack, expected := range data {
		actual := uriRegex.FindAllString(haystack, -1)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expected %v from %v, but found %v", expected, haystack, actual)
		}
	}
}
