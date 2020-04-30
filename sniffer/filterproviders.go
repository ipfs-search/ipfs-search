package sniffer

import (
	"context"
	t "github.com/ipfs-search/ipfs-search/types"
)

// Filter takes a provider, returning true if it is to be included or false when
// it is to be discarded.
type Filter interface {
	Filter(t.Provider) (bool, error)
}

// Return false for the first filter returning false, true otherwise
func shouldInclude(p t.Provider, filters []Filter) (bool, error) {
	for _, f := range filters {
		include, err := f.Filter(p)

		if err != nil {
			return false, err
		}

		if !include {
			return false, nil
		}
	}

	return true, nil
}

// filterProviders filters a stream of providers, retaining only those for
// which all filters return true.
func filterProviders(ctx context.Context, in <-chan t.Provider, out chan<- t.Provider, filters []Filter) error {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return context error
			return ctx.Err()
		case p := <-in:
			include, err := shouldInclude(p, filters)

			if err != nil {
				return err
			}

			if include {
				out <- p
			}
		}
	}
}
