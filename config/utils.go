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
			for _, newE := range findZeroElements(f.Interface()) {
				output = append(output, fmt.Sprintf("%s.%s", name, newE))
			}
		case reflect.Map:
			// Map type, require non-zero length
			if f.Len() == 0 {
				output = append(output, name)
			}

			// Recurse for map values
			// for _, e := range f.MapKeys() {
			// 	v := f.MapIndex(e)
			// 	findZeroElements(v.Interface())
			// }
		default:
			if f.Interface() == reflect.Zero(f.Type()).Interface() {
				output = append(output, name)
			}
		}
	}

	return output
}
