package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/regorov/euvv"
)

// Version holds a value assignable from outside by `go build`.
var Version string = "0.1"

func main() {
	var (
		// num receives VAT number.
		num string

		// help receives true if usage info requested.
		help bool

		// timeout receives integer value in seconds.
		timeout int

		// verbose receives true if verbose mode requested.
		verbose bool
	)

	// command line parsing.
	flag.StringVar(&num, "num", "", "VAT number")
	flag.BoolVar(&help, "h", false, "Print this screen")
	flag.IntVar(&timeout, "t", 3, "HTTP request timeout in seconds. Infinite if 0(zero)")
	flag.BoolVar(&verbose, "v", false, "Verbose mode on")

	flag.Parse()

	if help == true {
		fmt.Printf("euvv - EU VAT number validation tool. Version: %s\n", Version)
		flag.PrintDefaults()
		osExit(2, "")
	}

	// clean spaces before/after if $ euvv -num " "
	num = strings.TrimSpace(num)

	if len(num) == 0 {
		osExit(2, "Error: No VAT number specified! Use -h to get usage info.")
	}

	num = strings.ToUpper(num)

	// processing single VAT number.
	// formatValidationRequired = false.
	validator := euvv.New(false, timeout, verbose)
	ok, err := validator.Validate(num)
	if err != nil {
		osExit(3, err.Error())
	}

	fmt.Printf("VAT %s is %s.\n", num, map[bool]string{true: "valid", false: "invalid"}[ok])

	if ok {
		osExit(0, "")
	}

	osExit(1, "")
}

func osExit(code int, msg string) {
	if len(msg) > 0 {
		fmt.Println(msg)
	}
	os.Exit(code)
}
