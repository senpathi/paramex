package paramex

import (
	"errors"
	"testing"
)

func Test_Errors(t *testing.T) {
	t.Run(`test_ErrorUnSupportedParamType`, func(t *testing.T) {
		err := ErrorUnSupportedParamType{
			errors.New(`test error`),
		}
		if err.Error() != `test error` {
			t.Errorf(`expexted "test error", received "%v"`, err.Error())
		}
	})

	t.Run(`test_ErrorUnmarshalType`, func(t *testing.T) {
		err := ErrorUnmarshalType{
			errors.New(`test error`),
		}
		if err.Error() != `test error` {
			t.Errorf(`expexted "test error", received "%v"`, err.Error())
		}
	})

	t.Run(`test_ErrorNotAssignable`, func(t *testing.T) {
		err := ErrorNotAssignable{
			errors.New(`test error`),
		}
		if err.Error() != `test error` {
			t.Errorf(`expexted "test error", received "%v"`, err.Error())
		}
	})

	t.Run(`test_ErrorUnSupportedType`, func(t *testing.T) {
		err := ErrorUnSupportedType{
			errors.New(`test error`),
		}
		if err.Error() != `test error` {
			t.Errorf(`expexted "test error", received "%v"`, err.Error())
		}
	})
}
