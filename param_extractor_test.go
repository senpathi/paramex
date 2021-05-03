package paramex

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestParamExtractor_ExtractQuery(t *testing.T) {
	params :=
		"name=" + url.QueryEscape(`query_name`) + "&" +
			"age=" + url.QueryEscape(`20`) + "&" +
			"height=" + url.QueryEscape(`1.78`) + "&" +
			"married=" + url.QueryEscape(`false`)
	path := fmt.Sprintf("https://httpbin.org/get?%s", params)

	reqForm := url.Values{}
	reqForm.Set(`name`, `form_name`)
	reqForm.Set(`age`, `50`)
	reqForm.Set(`height`, `1.72`)
	reqForm.Set(`married`, `true`)

	req, err := http.NewRequest(`POST`, path, strings.NewReader(reqForm.Encode()))
	if err != nil {
		t.Error(`error creating request`, err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(reqForm.Encode())))

	req.Header.Set(`name`, `header_name`)
	req.Header.Set(`age`, `40`)
	req.Header.Set(`height`, `1.74`)
	req.Header.Set(`married`, `true`)

	headers := headerParams{}
	forms := formParams{}
	queries := queryParams{}

	extractor := NewParamExtractor()
	err = extractor.ExtractHeaders(&headers, req)
	if err != nil {
		t.Errorf(`error extracting headers due to %v`, err)
	}

	err = extractor.ExtractForms(&forms, req)
	if err != nil {
		t.Errorf(`error extracting forms due to %v`, err)
	}

	err = extractor.ExtractQueries(&queries, req)
	if err != nil {
		t.Errorf(`error extracting queries due to %v`, err)
	}

	fmt.Println(`extracted results ...`)
	fmt.Println(fmt.Sprintf(`request headers := %v`, headers))
	fmt.Println(fmt.Sprintf(`request forms := %v`, forms))
	fmt.Println(fmt.Sprintf(`request queries := %v`, queries))
}

type headerParams struct {
	Name    string  `param:"name"`
	Age     int64   `param:"age"`
	Height  float64 `param:"height"`
	Married bool    `param:"married"`
}

type formParams struct {
	Name    string  `param:"name"`
	Age     int     `param:"age"`
	Height  float64 `param:"height"`
	Married bool    `param:"married"`
}

type queryParams struct {
	Name    string  `param:"name"`
	Age     int32   `param:"age"`
	Height  float32 `param:"height"`
	Married bool    `param:"married"`
}
