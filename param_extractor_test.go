package paramex

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func makeRequest() (*http.Request, error) {
	params :=
		"name=" + url.QueryEscape(`query_name`) + "&" +
			"age=" + url.QueryEscape(`20`) + "&" +
			"height=" + url.QueryEscape(`1.78`) + "&" +
			"married=" + url.QueryEscape(`false`)
	path := fmt.Sprintf("https://nipuna.lk?%s", params)

	reqForm := url.Values{}
	reqForm.Set(`name`, `form_name`)
	reqForm.Set(`age`, `50`)
	reqForm.Set(`height`, `1.72`)
	reqForm.Set(`married`, `true`)

	req, err := http.NewRequest(`POST`, path, strings.NewReader(reqForm.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(reqForm.Encode())))

	req.Header.Set(`name`, `header_name`)
	req.Header.Set(`age`, `40`)
	req.Header.Set(`height`, `1.74`)
	req.Header.Set(`married`, `true`)

	return req, nil
}

func TestParamExtractor_ExtractQuery(t *testing.T) {
	req, err := makeRequest()
	if err != nil {
		t.Error(`error creating request`, err)
	}

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

func TestExtractor_ExtractHeaders_EmptyFieldTag(t *testing.T) {
	req, err := makeRequest()
	if err != nil {
		t.Error(`error creating request`, err)
	}

	empty := emptyTag{}
	extractor := NewParamExtractor()
	err = extractor.ExtractHeaders(&empty, req)
	if err != nil {
		t.Errorf(`error extracting headers due to %v`, err)
	}
	if empty.Name != `` {
		t.Errorf(`expected "", but received %v`, empty.Name)
	}
	if empty.Age != 0 {
		t.Errorf(`expected 0, but received %v`, empty.Age)
	}
	if empty.Height != 0 {
		t.Errorf(`expected 0, but received %v`, empty.Height)
	}
}

func TestExtractor_ExtractForms_EmptyFieldTag(t *testing.T) {
	req, err := makeRequest()
	if err != nil {
		t.Error(`error creating request`, err)
	}

	empty := emptyTag{}
	extractor := NewParamExtractor()
	err = extractor.ExtractForms(&empty, req)
	if err != nil {
		t.Errorf(`error extracting forms due to %v`, err)
	}
	if empty.Name != `` {
		t.Errorf(`expected "", but received %v`, empty.Name)
	}
	if empty.Age != 0 {
		t.Errorf(`expected 0, but received %v`, empty.Age)
	}
	if empty.Height != 0 {
		t.Errorf(`expected 0, but received %v`, empty.Height)
	}
}

func TestExtractor_ExtractQueries_EmptyFieldTag(t *testing.T) {
	req, err := makeRequest()
	if err != nil {
		t.Error(`error creating request`, err)
	}

	empty := emptyTag{}
	extractor := NewParamExtractor()
	err = extractor.ExtractQueries(&empty, req)
	if err != nil {
		t.Errorf(`error extracting quaries due to %v`, err)
	}
	if empty.Name != `` {
		t.Errorf(`expected "", but received %v`, empty.Name)
	}
	if empty.Age != 0 {
		t.Errorf(`expected 0, but received %v`, empty.Age)
	}
	if empty.Height != 0 {
		t.Errorf(`expected 0, but received %v`, empty.Height)
	}
}

func TestExtractor_Errors(t *testing.T) {
	req, err := makeRequest()
	if err != nil {
		t.Error(`error creating request`, err)
	}
	extractor := NewParamExtractor()

	t.Run(`test not assignable error`, func(t *testing.T) {
		header := headerParams{}
		err := extractor.ExtractHeaders(header, req)
		_, ok := err.(ErrorNotAssignable)
		if !ok {
			t.Errorf(`expected "ErrorNotAssignable", but received %v`, reflect.TypeOf(err))
			t.Fail()
		}
	})

	t.Run(`test non struct type error`, func(t *testing.T) {
		i := int64(0)
		err := extractor.ExtractHeaders(&i, req)
		_, ok := err.(ErrorUnSupportedType)
		if !ok {
			t.Errorf(`expected "ErrorUnSupportedType", but received %v`, reflect.TypeOf(err))
		}
	})

	t.Run(`test unsupported param type error`, func(t *testing.T) {
		obj := unSupportedTypeError{}
		err := extractor.ExtractHeaders(&obj, req)
		_, ok := err.(ErrorUnSupportedParamType)
		if !ok {
			t.Errorf(`expected "ErrorUnSupportedParamType", but received %v`, reflect.TypeOf(err))
			t.Fail()
		}
	})

	t.Run(`test unmarshal type error "bool"`, func(t *testing.T) {
		obj := boolError{}
		err := extractor.ExtractHeaders(&obj, req)
		if err == nil {
			t.Errorf(`expected "ErrorUnmarshalType", but received "%v"`, nil)
			return
		}
		_, ok := err.(ErrorUnmarshalType)
		if !ok {
			t.Errorf(`expected "ErrorUnmarshalType", but received "%v"`, reflect.TypeOf(err))
		}
		exErr := `error unmarshalling [header_name] into [bool] due to strconv.ParseBool: parsing "header_name": invalid syntax`
		if err.Error() != exErr {
			t.Errorf(`expexted [%v], but received [%v]`, exErr, err.Error())
		}
	})

	t.Run(`test unmarshal type error "int32"`, func(t *testing.T) {
		obj := int32Error{}
		err := extractor.ExtractHeaders(&obj, req)
		if err == nil {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, nil)
			return
		}
		_, ok := err.(ErrorUnmarshalType)
		if !ok {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, reflect.TypeOf(err))
		}
		exErr := `error unmarshalling [header_name] into [int32] due to strconv.Atoi: parsing "header_name": invalid syntax`
		if err.Error() != exErr {
			t.Errorf(`expexted [%v], but received [%v]`, exErr, err.Error())
		}
	})

	t.Run(`test unmarshal type error "int"`, func(t *testing.T) {
		obj := intError{}
		err := extractor.ExtractHeaders(&obj, req)
		if err == nil {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, nil)
			return
		}
		_, ok := err.(ErrorUnmarshalType)
		if !ok {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, reflect.TypeOf(err))
		}
		exErr := `error unmarshalling [header_name] into [int] due to strconv.Atoi: parsing "header_name": invalid syntax`
		if err.Error() != exErr {
			t.Errorf(`expexted [%v], but received [%v]`, exErr, err.Error())
		}
	})

	t.Run(`test unmarshal type error "int64"`, func(t *testing.T) {
		obj := int64Error{}
		err := extractor.ExtractHeaders(&obj, req)
		if err == nil {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, nil)
			return
		}
		_, ok := err.(ErrorUnmarshalType)
		if !ok {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, reflect.TypeOf(err))
		}
		exErr := `error unmarshalling [header_name] into [int64] due to strconv.ParseInt: parsing "header_name": invalid syntax`
		if err.Error() != exErr {
			t.Errorf(`expexted [%v], but received [%v]`, exErr, err.Error())
		}
	})

	t.Run(`test unmarshal type error "float32"`, func(t *testing.T) {
		obj := float32Error{}
		err := extractor.ExtractHeaders(&obj, req)
		if err == nil {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, nil)
			return
		}
		_, ok := err.(ErrorUnmarshalType)
		if !ok {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, reflect.TypeOf(err))
		}
		exErr := `error unmarshalling [header_name] into [float32] due to strconv.ParseFloat: parsing "header_name": invalid syntax`
		if err.Error() != exErr {
			t.Errorf(`expexted [%v], but received [%v]`, exErr, err.Error())
		}
	})

	t.Run(`test unmarshal type error "float64"`, func(t *testing.T) {
		obj := float64Error{}
		err := extractor.ExtractHeaders(&obj, req)
		if err == nil {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, nil)
			return
		}
		_, ok := err.(ErrorUnmarshalType)
		if !ok {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, reflect.TypeOf(err))
		}
		exErr := `error unmarshalling [header_name] into [float64] due to strconv.ParseFloat: parsing "header_name": invalid syntax`
		if err.Error() != exErr {
			t.Errorf(`expexted [%v], but received [%v]`, exErr, err.Error())
		}
	})

	t.Run(`test unmarshal type error "uuid.UUID{}"`, func(t *testing.T) {
		obj := uuidError{}
		err := extractor.ExtractHeaders(&obj, req)
		if err == nil {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, nil)
			return
		}
		_, ok := err.(ErrorUnmarshalType)
		if !ok {
			t.Errorf(`expected "ErrorUnmarshalType", but received %v`, reflect.TypeOf(err))
		}
		exErr := `error unmarshalling [header_name] into [uuid] due to invalid UUID length: 11`
		if err.Error() != exErr {
			t.Errorf(`expexted [%v], but received [%v]`, exErr, err.Error())
		}
	})
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

type emptyTag struct {
	Name   string  `param:"-"`
	Age    int32   `param:""`
	Height float32 `param:"test"`
}

type boolError struct {
	Name bool `param:"name"`
}

type int32Error struct {
	Name int32 `param:"name"`
}

type intError struct {
	Name int `param:"name"`
}

type int64Error struct {
	Name int64 `param:"name"`
}

type float32Error struct {
	Name float32 `param:"name"`
}

type float64Error struct {
	Name float64 `param:"name"`
}

type uuidError struct {
	Name uuid.UUID `param:"name"`
}

type unSupportedTypeError struct {
	Name map[string]string `param:"name"`
}
