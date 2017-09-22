package euvv_test

import (
	"testing"

	"github.com/regorov/euvv"
)

func TestVatValidator_Split(t *testing.T) {

	table := map[string]bool{
		"":              false, // empty
		"CZ":            false, // less then 3 chars
		"US12345678":    false, // non EU prefix
		"IT06700351213": true,
	}

	validator := euvv.New(true, 3, false)

	for k, v := range table {

		_, _, err := validator.Split(k)
		if (err == nil) == (v == false) {
			t.Errorf("Table element '%s' failed", k)
		}
	}
}
