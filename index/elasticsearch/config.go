package elasticsearch

import (
	"log"
	"reflect"
)

// Config represents the configuration for a specific index.
type Config struct {
	Name     string
	Settings map[string]interface{}
	Mapping  map[string]interface{}
}

// configEqual returns whether the config we want is equal to what we've got
func configEqual(want interface{}, got interface{}) bool {
	log.Printf("Comparing %v to %v", want, got)

	switch wantV := want.(type) {
	case map[string]interface{}:
		// Compare want values, if one fails, yield unequal
		for k, v := range wantV {

			gotV := got.(map[string]interface{})[k]

			if gotV == nil {
				log.Printf("Setting '%v' missing from %v", k, got)
				return false
			}

			// Recurse
			return configEqual(v, gotV)
		}
	default:
		r := reflect.DeepEqual(wantV, got)
		if !r {
			log.Printf("Setting %v not equal to %v", wantV, got)
		}
		return r
	}

	return true
}
