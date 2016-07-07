package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	files := os.Args[1:]
	existing := existingPaths(files)

	fmt.Println(existing)

	fmt.Printf("Scanning %v files...\n", len(existing))

	if len(existing) > 1 {
		for _, filepath := range existing {
			if uris := urisIn(filepath); uris != nil {
				for _, uri := range uris {
					fmt.Println(uri)
				}
			}
		}
	}
}

func existingPaths(paths []string) (existing []string) {
	existing = make([]string, 0, len(paths))

	for _, path := range paths {
		if len(path) == 0 {
			continue
		}

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			existing = append(existing, path)
		}
	}

	return existing
}

func urisIn(filepath string) (uris []string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error '%v': %v\n", filepath, err)
		return uris
	}

	// https://github.com/matthewrudy/regexpert/blob/master/lib/regexpert.rb
	uriRe := regexp.MustCompile(`(^(http|https):\/\/[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(([0-9]{1,5})?\/.*)?$)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := uriRe.FindAllString(line, -1); matches != nil {
			uris = append(uris, matches...)
		}
	}

	return uris
}
