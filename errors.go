package qs

import (
	"fmt"
	"reflect"
)

// IsRequiredFieldError returns ok==false if the given error wasn't caused by a
// required field that was missing from the query string.
// Otherwise it returns the name of the missing required field with ok==true.
func IsRequiredFieldError(e error) (fieldName string, ok bool) {
	if re, ok := e.(*reqError); ok {
		return re.FieldName, true
	}
	return "", false
}

// reqError is returned when a struct field marked with the 'req' option isn't
// in the unmarshaled url.Values or query string.
type reqError struct {
	Message   string
	FieldName string
}

func (e *reqError) Error() string {
	return e.Message
}

type wrongTypeError struct {
	Actual   reflect.Type
	Expected reflect.Type
}

func (e *wrongTypeError) Error() string {
	return fmt.Sprintf("received type %v, want %v", e.Actual, e.Expected)
}

type wrongKindError struct {
	Actual   reflect.Type
	Expected reflect.Kind
}

func (e *wrongKindError) Error() string {
	return fmt.Sprintf("received type %v of kind %v, want kind %v",
		e.Actual, e.Actual.Kind(), e.Expected)
}

type unhandledTypeError struct {
	Type reflect.Type
}

func (e *unhandledTypeError) Error() string {
	return fmt.Sprintf("unhandled type: %v", e.Type)
}
