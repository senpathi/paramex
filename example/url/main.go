package main

import (
	"fmt"
	"github.com/senpathi/paramex"
	"log"
	"net/http"
	"net/url"
)

type queryParams struct {
	Name    string  `param:"name"`
	Age     int32   `param:"age"`
	Height  float32 `param:"height"`
	Married bool    `param:"married"`
}

func main() {
	params :=
		"name=" + url.QueryEscape(`query_name`) + "&" +
			"age=" + url.QueryEscape(`20`) + "&" +
			"height=" + url.QueryEscape(`1.78`) + "&" +
			"married=" + url.QueryEscape(`false`)

	path := fmt.Sprintf("https://nipuna.lk?%s", params)

	req, err := http.NewRequest(`POST`, path, nil)
	if err != nil {
		log.Fatalln(err)
	}

	queries := queryParams{}
	extractor := paramex.NewParamExtractor()

	err = extractor.ExtractQueries(&queries, req)
	if err != nil {
		log.Fatalln(fmt.Errorf(`error extracting queries due to %v`, err))
	}

	fmt.Println(fmt.Sprintf(`request queries := %v`, queries))
	//Output : request queries := {query_name 20 1.78 false}
}
