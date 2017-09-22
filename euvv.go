// Package euvv implements EU VAT Number validation functions.
package euvv

import (
	"github.com/regorov/euvv/ecws"
)

// VatValidator provides high level functions for VAT number validation.
type VatValidator struct {

	// formatValidationRequired holds true if format pre-validation required.
	formatValidationRequired bool
	ws                       *ecws.CheckVatService
}

// New instantiates VatValidator.
func New(formatValidationRequired bool, timeout int, verboseMode bool) *VatValidator {
	return &VatValidator{
		formatValidationRequired: formatValidationRequired,
		ws: ecws.NewCheckVatService(timeout, verboseMode),
	}
}

// Validate returns true if number is valid.
func (vv *VatValidator) Validate(num string) (bool, error) {

	cc, vn, err := vv.Split(num)

	if err != nil {
		return false, err
	}

	resp, err := vv.ws.CheckVat(&ecws.CheckVat{CountryCode: cc, VatNumber: vn})
	if err != nil {
		return false, err
	}

	return resp.Valid, nil
}

// Split cuts input string "CZ28987373" to "CZ" and "28987373".
// Input string expected in upper case. Returns ErrUnknownCountryPrefix if
// country prefix unknown, returns ErrInvalidNumberFormat if number part validation
// falied. Returns ErrInvalidInputData in other cases.
func (vv *VatValidator) Split(s string) (country, num string, err error) {

	if len(s) < 3 {
		return "", "", ErrInvalidInputData
	}

	// currently formatValidationRequired is always false.
	if vv.formatValidationRequired == false {
		return s[:2], s[2:], nil
	}

	f, ok := vatFormatValidator[s[:2]]
	if ok == false {
		return "", "", ErrUnknownCountryPrefix
	}

	if f != nil {
		if f(s[2:]) == false {
			return "", "", ErrInvalidNumberFormat
		}
	}

	return s[:2], s[2:], nil
}
