package index

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"log"
	"reflect"
)

// wantGotEqual asserts that the settings mentioned in want equal the got
func wantGotEqual(want interface{}, got interface{}) bool {
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
			return wantGotEqual(v, gotV)
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

func assertSettings(ctx context.Context, i ManagedIndex, settings interface{}, update bool) error {
	existingSettings, err := i.GetSettings(ctx)
	if err != nil {
		return fmt.Errorf("index %v, getting settings: %w", i, err)
	}

	if wantGotEqual(settings, existingSettings) {
		log.Printf("Index '%v' settings up to date", i)
		return nil
	}

	// Below is debug only
	diff := cmp.Diff(settings, existingSettings)
	log.Printf("Settings do not match (-want +got):\n%s", diff)

	switch update {
	case true:
		if err := i.SetSettings(ctx, settings); err != nil {
			return fmt.Errorf("index %v, updating settings: %w", i, err)
		}

		// Confirm update
		if existingSettings, err = i.GetSettings(ctx); err != nil {
			return fmt.Errorf("index %v, getting settings: %w", i, err)
		}

		// It should match now
		if !wantGotEqual(settings, existingSettings) {
			return fmt.Errorf("index %v, updated settings do not match target settings", i)
		}

		log.Println("settings updated")
	case false:
		log.Println("not updating settings")
	}

	return nil
}

func assertMapping(ctx context.Context, i ManagedIndex, mapping interface{}, update bool) error {
	existingMapping, err := i.GetMapping(ctx)
	if err != nil {
		return fmt.Errorf("index %v, getting mapping: %w", i, err)
	}

	if wantGotEqual(mapping, existingMapping) {
		log.Printf("Index '%v' mapping up to date", i)
		return nil
	}

	// Below is debug only
	diff := cmp.Diff(mapping, existingMapping)
	log.Printf("Mapping does not match (-want +got):\n%s", diff)

	switch update {
	case true:
		if err := i.SetMapping(ctx, mapping); err != nil {
			return fmt.Errorf("index %v, updating mapping: %w", i, err)
		}

		// Confirm update
		if existingMapping, err = i.GetMapping(ctx); err != nil {
			return fmt.Errorf("index %v, getting mapping: %w", i, err)
		}

		// It should match now
		if !wantGotEqual(mapping, existingMapping) {
			return fmt.Errorf("index %v, updated mapping do not match target mapping", i)
		}

		log.Println("mapping updated")
	case false:
		log.Println("not updating mapping")
	}

	return nil
}

// Assert checks for the existence of an index with given settings, creates it if necessary and optionally attempts to update them.
func Assert(ctx context.Context, i ManagedIndex, settings interface{}, mapping interface{}, update bool) error {
	exists, err := i.Exists(ctx)
	if err != nil {
		return fmt.Errorf("index %v, exists: %w", i, err)
	}

	if !exists {
		log.Printf("Creating index \"%v\"", i)
		return i.Create(ctx, settings, mapping)
	}

	if err := assertSettings(ctx, i, settings, update); err != nil {
		return err
	}

	return assertMapping(ctx, i, mapping, update)
}
