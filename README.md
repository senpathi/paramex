# paramex

![paramex_logo](./docs/images/logo.png)

[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/senpathi/paramex)
[![build](https://github.com/senpathi/paramex/workflows/build/badge.svg)](https://github.com/senpathi/paramex/actions)
[![Coverage](https://codecov.io/gh/senpathi/paramex/branch/master/graph/badge.svg)](https://codecov.io/gh/senpathi/paramex)
[![Releases](https://img.shields.io/github/release/senpathi/paramex/all.svg?style=flat-square)](https://github.com/senpathi/paramex/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/senpathi/paramex)](https://goreportcard.com/report/github.com/senpathi/paramex)
[![LICENSE](https://img.shields.io/github/license/senpathi/paramex.svg?style=flat-square)](https://github.com/senpathi/paramex/blob/master/LICENSE)

Paramex is a library that binds `http request parameters` to a Go struct annotated with `param`.

## Description

To extract http parameters `(headers, url query values, form values)`, multiple code lines need to be written in
`http Handlers`. 

But, Using **Paramex** `http headers, url query values or form values` can be extracted by calling a single function.

Sample code example code to extract request form values using `paramex` is shown below.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/senpathi/paramex"
)

type formParams struct {
	Name    string  `param:"name"`
	Age     int     `param:"age"`
	Height  float64 `param:"height"`
	Married bool    `param:"married"`
}

func main() {
	reqForm := url.Values{}
	reqForm.Set(`name`, `form_name`)
	reqForm.Set(`age`, `50`)
	reqForm.Set(`height`, `1.72`)
	reqForm.Set(`married`, `true`)

	req, err := http.NewRequest(`POST`, `https://nipuna.lk`, strings.NewReader(reqForm.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	forms := formParams{}
	extractor := paramex.NewParamExtractor()

	err = extractor.ExtractForms(&forms, req)
	if err != nil {
		log.Fatalln(fmt.Errorf(`error extracting forms due to %v`, err))
	}

	fmt.Println(fmt.Sprintf(`request forms := %v`, forms))  
	//Output : request forms := {form_name 50 1.72 true}
}
```

Examples codes to extract http headers, url query values and form values are implemented in 
[example](https://github.com/senpathi/paramex/tree/master/example) directory.

### Supported parameter types

 - string
 - bool
 - int32
 - int
 - int64
 - float32
 - float64