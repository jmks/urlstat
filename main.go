package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/fatih/color"
	"github.com/jmks/uristat/options"
)

// https://github.com/matthewrudy/regexpert/blob/master/lib/regexpert.rb
var uriRegex = regexp.MustCompile(`(^(http|https):\/\/[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(([0-9]{1,5})?\/.*)?$)`)

func main() {
	opts := options.Parse()

	if !opts.IsValid() {
		opts.PrintError()
		os.Exit(1)
	}

	filepathSrc := filepathProducer(opts.Filepaths)
	uriSrc := uriProducer(filepathSrc)
	uniqURIs := uniqAccumulator(uriSrc)

	if opts.ListOnly() {
		for _, uri := range uniqURIs {
			fmt.Println(uri)
		}
	} else {
		printStatuses(uniqURIs)
	}
}

func filepathProducer(filepaths []string) <-chan string {
	dest := make(chan string, 100)

	go func() {
		for _, filepath := range filepaths {
			if len(filepath) == 0 {
				continue
			}

			if _, err := os.Stat(filepath); !os.IsNotExist(err) {
				dest <- filepath
			}
		}
		close(dest)
	}()

	return dest
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

func uniqAccumulator(src <-chan string) []string {
	uniq := make(map[string]bool)

	for str := range src {
		uniq[str] = true
	}

	uris := make([]string, len(uniq))
	i := 0
	for uri := range uniq {
		uris[i] = uri
		i++
	}

	return uris
}

func urisIn(filepath string) (uris []string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error '%v': %v\n", filepath, err)
		return uris
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := uriRegex.FindAllString(line, -1); matches != nil {
			uris = append(uris, matches...)
		}
	}

	return uris
}

func printStatuses(uris []string) {
	wg := sync.WaitGroup{}

	for _, uri := range uris {
		wg.Add(1)

		go func(uri string) {
			if resp, err := http.Head(uri); err != nil {
				redden := color.New(color.FgRed).SprintFunc()
				fmt.Printf("%v : %v\n", redden("HTTP ERROR"), uri)
			} else {
				defer resp.Body.Close()
				colorize := statusCodePrinterFunc(resp.StatusCode)
				fmt.Printf("%v : %v\n", colorize(resp.Status), uri)
			}

			wg.Done()
		}(uri)
	}

	wg.Wait()
}

func statusCodePrinterFunc(code int) func(...interface{}) string {
	switch code {
	case 200:
		return color.New(color.FgGreen).SprintFunc()
	case 404:
		return color.New(color.FgYellow).SprintFunc()
	default:
		return color.New(color.FgRed).SprintFunc()
	}
}
