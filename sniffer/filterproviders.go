package sniffer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/sniffer/filters"
	t "github.com/ipfs-search/ipfs-search/types"
)

type providerFilter struct {
	f filters.Filter
}

// filter filters a stream of providers, dropping those for which filter returns false
func (f *providerFilter) filter(ctx context.Context, in <-chan t.Provider, out chan<- t.Provider) error {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return context error
			return ctx.Err()
		case p := <-in:
			include, err := f.f.Filter(p)

			if err != nil {
				return err
			}

			if include {
				out <- p
			}
		}
	}

}
