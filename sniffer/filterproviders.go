package sniffer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/sniffer/filters"
	t "github.com/ipfs-search/ipfs-search/types"
)

// filterProviders filters a stream of providers, dropping those
// for which filter returns false
func filterProviders(ctx context.Context, in <-chan t.Provider, out chan<- t.Provider, f filters.Filter) error {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return context error
			return ctx.Err()
		case p := <-in:
			include, err := f.Filter(p)

			if err != nil {
				return err
			}

			if include {
				out <- p
			}
		}
	}
}
