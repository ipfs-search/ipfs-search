package config

import (
	"fmt"
	"reflect"
)

// findZeroElements returns a slice of all (nested) struct fields with a zero value.
func findZeroElements(s interface{}) []string {
	var output []string

	v := reflect.ValueOf(s)

	// Iterate over fields
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		name := v.Type().Field(i).Tag.Get("yaml")

		switch f.Kind() {
		case reflect.Struct:
			// It's a struct - recurse!
			for _, new_e := range findZeroElements(f.Interface()) {
				output = append(output, fmt.Sprintf("%s.%s", name, new_e))
			}
		default:
			if f.Interface() == reflect.Zero(f.Type()).Interface() {
				output = append(output, name)
			}
		}
	}

	return output
}
