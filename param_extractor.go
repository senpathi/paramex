package paramex

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/google/uuid"
)

type extractorFunc func(key string) (string, bool)

const (
	stringType  = string("")
	boolType    = true
	int32Type   = int32(0)
	intType     = int(0)
	int64Type   = int64(0)
	float32Type = float32(0)
	float64Type = float64(0)
)

// The Extractor interface is implemented to extract http request headers, form values
// and url query values
type Extractor interface {
	// ExtractHeaders extract http headers from sent request and binds to `v`
	// `v` should be a Go struct reference
	ExtractHeaders(v interface{}, req *http.Request) error

	// ExtractQueries extract http url parameters from sent request and binds to v
	// `v` should be a Go struct reference
	ExtractQueries(v interface{}, req *http.Request) error

	// ExtractForms extract http form values from sent request and binds to v
	// `v` should be a Go struct reference
	ExtractForms(v interface{}, req *http.Request) error
}

type extractor struct{}

// NewParamExtractor returns an Extractor which extract
// req.Header, req.FormValue, req.URL.Query values and
// binds them to a Go struct
func NewParamExtractor() Extractor {
	return extractor{}
}

// ExtractHeaders extract http headers from sent request and binds to v
func (p extractor) ExtractHeaders(v interface{}, req *http.Request) error {
	return p.extract(v, func(key string) (string, bool) {
		str := req.Header.Get(key)
		if str == `` {
			return str, false
		}
		return str, true
	})
}

// ExtractQueries extract http url parameters from sent request and binds to v
func (p extractor) ExtractQueries(v interface{}, req *http.Request) error {
	return p.extract(v, func(key string) (string, bool) {
		str := req.URL.Query().Get(key)
		if str == `` {
			return str, false
		}
		return str, true
	})
}

// ExtractForms extract http form values from sent request and binds to v
func (p extractor) ExtractForms(v interface{}, req *http.Request) error {
	return p.extract(v, func(key string) (string, bool) {
		str := req.FormValue(key)
		if str == `` {
			return str, false
		}
		return str, true
	})
}

func (p extractor) extract(v interface{}, keyExtractor extractorFunc) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr || v == nil {
		return ErrorNotAssignable{
			fmt.Errorf(`type of %v is not assignabale, required object reference`, t)}
	}

	elem := reflect.ValueOf(v).Elem()
	if elem.Kind() != reflect.Struct {
		return ErrorUnSupportedType{
			fmt.Errorf(`type of %v is not extractable, required struct object`, elem.Type().String())}
	}

	t = elem.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag, ok := field.Tag.Lookup(`param`)
		if !ok || tag == `-` {
			continue
		}

		valueStr, ok := keyExtractor(tag)
		if !ok {
			continue
		}

		switch field.Type {
		case reflect.TypeOf(stringType):
			elem.Field(i).Set(reflect.ValueOf(valueStr))

		case reflect.TypeOf(boolType):
			value, err := strconv.ParseBool(valueStr)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [bool] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		case reflect.TypeOf(int32Type):
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [int32] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(int32(value)))

		case reflect.TypeOf(intType):
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [int] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		case reflect.TypeOf(int64Type):
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [int64] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		case reflect.TypeOf(float32Type):
			value, err := strconv.ParseFloat(valueStr, 32)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [float32] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(float32(value)))

		case reflect.TypeOf(float64Type):
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [float64] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		case reflect.TypeOf(uuid.UUID{}):
			value, err := uuid.Parse(valueStr)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [uuid] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		default:
			return ErrorUnSupportedParamType{errors.New(`unsupported param extractor type`)}
		}
	}

	return nil
}
