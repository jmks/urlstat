package options

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// Options parsed at the command line
type Options struct {
	list      *bool
	ok        *bool
	notOk     *bool
	Filepaths []string
}

// Parse arguments and flags and returns options configuration struct
func Parse() (opts Options) {
	opts.list = flag.Bool("list", false, "only list URIs found in files (i.e. no status check)")
	opts.ok = flag.Bool("ok", false, "only list URIs with HTTP status code 200 OK")
	opts.notOk = flag.Bool("no-ok", false, "list URIs with HTTP status code other than 200 OK (overrides --ok)")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of URLstat: urlstat [options] files...")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	opts.populateFilepaths()

	return opts
}

func (opts *Options) populateFilepaths() {
	if flag.Parsed() {
		opts.Filepaths = flag.Args()
	} else {
		opts.Filepaths = os.Args[1:]
	}

	if len(opts.Filepaths) == 0 {
		inStat, _ := os.Stdin.Stat()

		if (inStat.Mode() & os.ModeCharDevice) == 0 {
			stdinScanner := bufio.NewScanner(os.Stdin)

			for stdinScanner.Scan() {
				opts.Filepaths = append(opts.Filepaths, stdinScanner.Text())
			}
		}
	}
}

// IsValid returns whether Options is valid
func (opts Options) IsValid() bool {
	return len(opts.Filepaths) > 0
}

// PrintError prints reason Options was invalid and usage info to stderr
func (opts Options) PrintError() {
	if len(opts.Filepaths) == 0 {
		fmt.Fprintln(os.Stderr, "No files to scan")
	}

	fmt.Fprintln(os.Stderr, "")
	flag.Usage()
}

// ListOnly returns bool indicating if the found URIs should be printed
func (opts Options) ListOnly() bool {
	return *opts.list
}

// IsOkListable returns true when OK responses should be printed
func (opts Options) IsOkListable() bool {
	if *opts.ok || *opts.notOk {
		return !*opts.notOk
	}

	return true
}

// IsNotOkListable returns true when non-OK responses should be printed
func (opts Options) IsNotOkListable() bool {
	if *opts.ok || *opts.notOk {
		return *opts.notOk
	}

	return true
}
