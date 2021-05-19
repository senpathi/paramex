package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/senpathi/paramex"
)

type headerParams struct {
	Name    string  `param:"name"`
	Age     int64   `param:"age"`
	Height  float64 `param:"height"`
	Married bool    `param:"married"`
}

func main() {
	req, err := http.NewRequest(`POST`, `https://nipuna.lk`, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set(`name`, `header_name`)
	req.Header.Set(`age`, `40`)
	req.Header.Set(`height`, `1.74`)
	req.Header.Set(`married`, `true`)

	headers := headerParams{}
	extractor := paramex.NewParamExtractor()

	err = extractor.ExtractHeaders(&headers, req)
	if err != nil {
		log.Fatalln(fmt.Errorf(`error extracting headers due to %v`, err))
	}

	fmt.Println(fmt.Sprintf(`request headers := %v`, headers))
	//Output : request headers := {header_name 40 1.74 true}
}
