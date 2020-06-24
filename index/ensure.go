package index

import (
	"context"
	"fmt"
	"log"
)

func ensureExists(ctx context.Context, i ManagedIndex) error {
	exists, err := i.Exists(ctx)
	if err != nil {
		return fmt.Errorf("index %v, exists: %w", i, err)
	}

	if !exists {
		log.Printf("Creating index \"%v\"", i)
		err := i.Create(ctx)

		if err != nil {
			return fmt.Errorf("index %v, create: %w", i, err)
		}
	}

	return nil
}

func ensureConfigUpToDate(ctx context.Context, i ManagedIndex) error {
	equal, err := i.ConfigUpToDate(ctx)
	if err != nil {
		return err
	}

	if !equal {
		// Attempt update
		if err := i.ConfigUpdate(ctx); err != nil {
			return err
		}

		// Verify update
		equal, err := i.ConfigUpToDate(ctx)
		if err != nil {
			return err
		}

		if equal {
			log.Printf("Configuration updated for %v", i)
		} else {
			return fmt.Errorf("configuration update result not equal for index '%v'", i)
		}
	}

	return nil
}

// ensureExistsAndUpdated ensures that an index exists and that the configuration is up to date.
func ensureExistsAndUpdated(ctx context.Context, i ManagedIndex) error {
	if err := ensureExists(ctx, i); err != nil {
		return err
	}

	if err := ensureConfigUpToDate(ctx, i); err != nil {
		return err
	}

	return nil
}

// EnsureExistsAndUpdatedMulti ensures, for a string mapping of indexes, that they exist and that their configuration is up to date.
func EnsureExistsAndUpdatedMulti(ctx context.Context, indexes ...ManagedIndex) error {
	for _, i := range indexes {
		if err := ensureExistsAndUpdated(ctx, i); err != nil {
			return err
		}
	}

	return nil
}
