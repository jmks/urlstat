package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/fatih/color"
)

func main() {
	onlyList := flag.Bool("list", false, "only list URIs found in files (i.e. no status check)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of URIstat:")
		flag.PrintDefaults()
	}

	flag.Parse()

	filepaths := getFilepaths()

	if len(filepaths) == 0 {
		fmt.Println("No files to scan")
		os.Exit(1)
	}

	existingSrc := filepathProducer(filepaths)
	uriSrc := uriProducer(existingSrc)

	if *onlyList {
		for uri := range uriSrc {
			fmt.Println(uri)
		}
	} else {
		printStatuses(uriSrc)
	}
}

func getFilepaths() (filepaths []string) {
	if flag.Parsed() {
		filepaths = flag.Args()
	} else {
		filepaths = os.Args[1:]
	}

	if len(filepaths) == 0 {
		inStat, _ := os.Stdin.Stat()

		if (inStat.Mode() & os.ModeCharDevice) == 0 {
			stdinScanner := bufio.NewScanner(os.Stdin)

			for stdinScanner.Scan() {
				filepaths = append(filepaths, stdinScanner.Text())
			}
		}
	}

	return filepaths
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
				color.Red("%v error: %v\n", uri, err)
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
