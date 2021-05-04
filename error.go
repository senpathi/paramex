package paramex

type ErrorUnSupportedParamType struct{ err error }

func (e ErrorUnSupportedParamType) Error() string {
	return e.err.Error()
}

type ErrorUnmarshalType struct{ err error }

func (e ErrorUnmarshalType) Error() string {
	return e.err.Error()
}

type ErrorNotAssignable struct{ err error }

func (e ErrorNotAssignable) Error() string {
	return e.err.Error()
}

type ErrorUnSupportedType struct{ err error }

func (e ErrorUnSupportedType) Error() string {
	return e.err.Error()
}
