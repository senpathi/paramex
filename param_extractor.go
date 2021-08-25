package paramex

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/google/uuid"
)

type extractorFunc func(key string, array bool) (interface{}, bool)

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
	return p.extract(v, func(key string, array bool) (interface{}, bool) {
		if array {
			return nil, false
		}
		str := req.Header.Get(key)
		if str == `` {
			return str, false
		}
		return str, true
	})
}

// ExtractQueries extract http url parameters from sent request and binds to v
func (p extractor) ExtractQueries(v interface{}, req *http.Request) error {
	return p.extract(v, func(key string, array bool) (interface{}, bool) {
		str := req.URL.Query()[key]
		if len(str) == 0 {
			return "", false
		}

		if !array {
			return str[0], true
		}
		return str, true
	})
}

// ExtractForms extract http form values from sent request and binds to v
func (p extractor) ExtractForms(v interface{}, req *http.Request) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}
	return p.extract(v, func(key string, array bool) (interface{}, bool) {
		str := req.PostForm[key]
		if len(str) == 0 {
			return "", false
		}

		if !array {
			return str[0], true
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

		_, ok = keyExtractor(tag, false)
		if !ok {
			continue
		}

		switch field.Type {
		case reflect.TypeOf(stringType):
			valueStr, _ := keyExtractor(tag, false)
			elem.Field(i).Set(reflect.ValueOf(valueStr.(string)))

		case reflect.TypeOf(boolType):
			valueStr, _ := keyExtractor(tag, false)
			value, err := strconv.ParseBool(valueStr.(string))
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [bool] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		case reflect.TypeOf(int32Type):
			valueStr, _ := keyExtractor(tag, false)
			value, err := strconv.Atoi(valueStr.(string))
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [int32] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(int32(value)))

		case reflect.TypeOf(intType):
			valueStr, _ := keyExtractor(tag, false)
			value, err := strconv.Atoi(valueStr.(string))
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [int] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		case reflect.TypeOf(int64Type):
			valueStr, _ := keyExtractor(tag, false)
			value, err := strconv.ParseInt(valueStr.(string), 10, 64)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [int64] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		case reflect.TypeOf(float32Type):
			valueStr, _ := keyExtractor(tag, false)
			value, err := strconv.ParseFloat(valueStr.(string), 32)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [float32] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(float32(value)))

		case reflect.TypeOf(float64Type):
			valueStr, _ := keyExtractor(tag, false)
			value, err := strconv.ParseFloat(valueStr.(string), 64)
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [float64] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		case reflect.TypeOf(uuid.UUID{}):
			valueStr, _ := keyExtractor(tag, false)
			value, err := uuid.Parse(valueStr.(string))
			if err != nil {
				return ErrorUnmarshalType{
					fmt.Errorf(`error unmarshalling [%v] into [uuid] due to %v`, valueStr, err)}
			}
			elem.Field(i).Set(reflect.ValueOf(value))

		case reflect.TypeOf([]string{stringType}):
			valueStr, ok := keyExtractor(tag, true)
			if !ok {
				return ErrorUnSupportedParamType{
					fmt.Errorf(`error unmarshalling []string into "%v", unsupported param type`, tag),
				}
			}
			elem.Field(i).Set(reflect.ValueOf(valueStr.([]string)))

		default:
			return ErrorUnSupportedParamType{errors.New(`unsupported param extractor type`)}
		}
	}

	return nil
}
