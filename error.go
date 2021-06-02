package paramex

// ErrorUnSupportedParamType created when trying to extract unsupported parameter type
type ErrorUnSupportedParamType struct {
	error
}

// ErrorUnmarshalType created when trying to marshal different type value to another type variable
type ErrorUnmarshalType struct {
	error
}

// ErrorNotAssignable created when sent interface to extract  is not assignable. Not a reference to a Go struct
type ErrorNotAssignable struct {
	error
}

// ErrorUnSupportedType created when sent reference  is not a Go struct type reference
type ErrorUnSupportedType struct {
	error
}
