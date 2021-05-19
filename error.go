package paramex

// ErrorUnSupportedParamType created when trying to extract unsupported parameter type
type ErrorUnSupportedParamType struct{ err error }

func (e ErrorUnSupportedParamType) Error() string {
	return e.err.Error()
}

// ErrorUnmarshalType created when trying to marshal different type value to another type variable
type ErrorUnmarshalType struct{ err error }

func (e ErrorUnmarshalType) Error() string {
	return e.err.Error()
}

// ErrorNotAssignable created when sent interface to extract  is not assignable. Not a reference to a Go struct
type ErrorNotAssignable struct{ err error }

func (e ErrorNotAssignable) Error() string {
	return e.err.Error()
}

// ErrorUnSupportedType created when sent reference  is not a Go struct type reference
type ErrorUnSupportedType struct{ err error }

func (e ErrorUnSupportedType) Error() string {
	return e.err.Error()
}
