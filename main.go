package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/jmks/urlstat/options"
	"github.com/jmks/urlstat/tld"
)

func main() {
	opts := options.Parse()

	if !opts.IsValid() {
		opts.PrintError()
		os.Exit(1)
	}

	filepathSrc := filepathProducer(opts.Filepaths)
	urlSrc := urlProducer(filepathSrc)
	uniqURLs := uniqAccumulator(urlSrc)

	if opts.ListOnly() {
		for _, url := range uniqURLs {
			fmt.Println(url)
		}
	} else {
		printStatuses(uniqURLs, opts)
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

func urlProducer(filepathSrc <-chan string) <-chan string {
	dest := make(chan string, 100)

	go func() {
		wg := sync.WaitGroup{}

		for filepath := range filepathSrc {
			wg.Add(1)

			go func(path string) {
				defer wg.Done()

				file, err := os.Open(path)
				if err != nil {
					fmt.Printf("Error '%v'\n", err)
					return
				}

				for _, url := range extractURLs(file) {
					dest <- url
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

	urls := make([]string, len(uniq))
	i := 0
	for url := range uniq {
		urls[i] = url
		i++
	}

	return urls
}

func extractURLs(source io.Reader) []string {
	var urls []string

	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		for _, word := range strings.Fields(scanner.Text()) {
			u, err := url.Parse(word)
			if err != nil {
				continue
			}

			// skip all non-http[s]? schemes
			if len(u.Scheme) > 0 && !strings.HasPrefix(u.Scheme, "http") {
				continue
			}

			if len(u.Host) > 0 && !tld.HasKnownTLD(u.Host) {
				continue
			}

			if len(u.Host) == 0 && len(u.Path) > 0 && !(looksLikeURN(u.Path) && tld.HasKnownTLD(u.Path)) {
				continue
			}

			if len(u.Host) == 0 && len(u.Path) == 0 {
				continue
			}

			urls = append(urls, u.String())
		}
	}

	return urls
}

var urnPattern = regexp.MustCompile(`^(?P<host>(?:\w+\.)+)(?P<tld>\w+).*`)

func looksLikeURN(s string) bool {
	return urnPattern.MatchString(s)
}

func printStatuses(urls []string, opts options.Options) {
	wg := sync.WaitGroup{}

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			if resp, err := http.Head(url); err != nil {
				redden := color.New(color.FgRed).SprintFunc()
				fmt.Printf("%v : %v\n", redden("HTTP ERROR"), url)
			} else {
				defer resp.Body.Close()

				if isStatusPrintable(resp.StatusCode, opts) {
					colorize := statusCodePrinterFunc(resp.StatusCode)
					fmt.Printf("%v : %v\n", colorize(resp.Status), url)
				}
			}

			wg.Done()
		}(url)
	}

	wg.Wait()
}

func isStatusPrintable(status int, opts options.Options) bool {
	return (status == 200 && opts.IsOkListable()) || (status != 200 && opts.IsNotOkListable())
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
