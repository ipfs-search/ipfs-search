package sniffer

import (
	"context"
	"log"
)

func shouldFilter(p Provider, filters []Filter) bool {
	// The first filter returning false gets a resource skipped
	for _, f := range filters {
		if !f.Filter(p) {
			return false
		}
	}

	return true
}

func filterProviders(ctx context.Context, in <-chan Provider, out chan<- Provider, filters []Filter) error {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return context error
			return ctx.Err()
		case p := <-in:
			if !shouldFilter(p, filters) {
				log.Printf("Filtering %v", p.Resource)
				continue
			}

			out <- p
		}
	}
}
