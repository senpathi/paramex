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
