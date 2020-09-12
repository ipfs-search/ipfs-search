package filter

import (
	"context"
	"github.com/ipfs-search/ipfs-search/sniffer/filters"
	t "github.com/ipfs-search/ipfs-search/types"
)

type Filter struct {
	f   filters.Filter
	in  <-chan t.Provider
	out chan<- t.Provider
}

func New(f filters.Filter, in <-chan t.Provider, out chan<- t.Provider) Filter {
	return Filter{
		f:   f,
		in:  in,
		out: out,
	}
}

// Filter filters a stream of providers, dropping those for which filter returns false
func (f *Filter) Filter(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return context error
			return ctx.Err()
		case p := <-f.in:
			include, err := f.f.Filter(p)

			if err != nil {
				return err
			}

			if include {
				f.out <- p
			}
		}
	}

}
