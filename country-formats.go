package euvv

import (
	"errors"
)

var (
	ErrUnknownCountryPrefix = errors.New("unknown country prefix")
	ErrInvalidNumberFormat  = errors.New("invalid number format")
	ErrInvalidInputData     = errors.New("invalid input data")
)

// vatFormatValidator holds references to format validation
// functions per EU country. Countries have different VAT formats due to
// https://en.wikipedia.org/wiki/VAT_identification_number and some times
// cannot be fixed by regexp.
var vatFormatValidator = map[string]func(string) bool{
	"CZ": cz,
	"DE": de,
	"SE": se,
	"NL": nl,
	"GB": gb,
	"PL": pl,
	"IT": it,
	// TODO - extend list of supported prefixes.
}

func cz(num string) bool {
	// Add numeric part validation if required
	// pre-validation before sending to the web service.
	return true
}

func de(num string) bool {
	return true
}

func se(num string) bool {
	return true
}

func nl(num string) bool {
	return true
}

func gb(num string) bool {
	return true
}

func pl(num string) bool {
	return true
}

func it(num string) bool {
	return true
}
