package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
)

func main() {
	files := os.Args[1:]

	existingSrc := existingProducer(files)
	uriSrc := uriProducer(existingSrc)
	printStatuses(uriSrc)
}

func existingProducer(filepaths []string) <-chan string {
	c := make(chan string, 100)

	go func() {
		for _, filepath := range filepaths {
			if len(filepath) == 0 {
				continue
			}

			if _, err := os.Stat(filepath); !os.IsNotExist(err) {
				c <- filepath
			}
		}
		close(c)
	}()

	return c
}

func uriProducer(filepathSrc <-chan string) <-chan string {
	dest := make(chan string, 100)

	go func() {
		wg := sync.WaitGroup{}

		for filepath := range filepathSrc {
			wg.Add(1)

			go func(path string) {
				defer wg.Done()

				if uris := urisIn(path); len(uris) > 0 {
					for _, uri := range uris {
						dest <- uri
					}
				}
			}(filepath)
		}

		wg.Wait()
		close(dest)
	}()

	return dest
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

func printStatuses(uriSrc <-chan string) {
	wg := sync.WaitGroup{}

	for uri := range uriSrc {
		wg.Add(1)

		go func(uri string) {
			if resp, err := http.Head(uri); err != nil {
				fmt.Printf("%v error: %v\n", uri, err)
			} else {
				defer resp.Body.Close()

				fmt.Printf("%v : %v\n", resp.Status, uri)
			}

			wg.Done()
		}(uri)
	}

	wg.Wait()
}
