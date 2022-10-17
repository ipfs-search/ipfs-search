package cache

import (
	"reflect"
)

// GetStructElem expects an interface containing a pointer to a struct and returns the value thereof.
func GetStructElem(s interface{}) reflect.Value {
	// Source value
	src := reflect.ValueOf(s)

	// Ensure src is a pointer pointer
	if src.Kind() != reflect.Pointer {
		panic("not called with pointer")
	}

	// Dereference src
	src = src.Elem()

	// Ensure src is a struct
	if src.Kind() != reflect.Struct {
		panic("not struct pointer")
	}

	return src
}
