# euvv

EU VAT Number Validation Golang Package and Command Line Tool. Verify the validity of a VAT number issued by any Member State.

[![GoDoc](https://godoc.org/github.com/regorov/euvv?status.svg)](https://godoc.org/github.com/regorov/euvv)

## Usage

### Command Line Tool
```
$ euvv -num CZ28987373
VAT CZ28987373 is valid.

$ euvv -num CZ289873731
VAT CZ28987373 is invalid.

```

**Exit codes:**
* 0 - VAT number is valid.
* 1 - VAT number is invalid.
* 2 - Execution broken due to invalid command line parameters.
* 3 - Execution failed. See text message.

### From Go Application
```Go
	import (
		"github.com/regorov/euvv"
	)

	validator := euvv.New(false, 3, false)
	ok, err := validator.Validate("CZ28987373")
	if err != nil {
		fmt.Println(err)
		return
	}

	if ok {
		fmt.Println("Valid")
	} else {
		fmt.Println("Invalid")
	}
```

## Install and Run
* Intall Go www.golang.org
* $ go get github.com/regorov/euvv
* $ cd $GOPATH/src/github.com/regorov/euvv/cmd/euvv
* $ go build
* $ ./euvv -num CZ28987373

## Credits
Partially generated via https://github.com/hooklift/gowsdl

## License
MIT
