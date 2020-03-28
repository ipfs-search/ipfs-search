package sniffer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs/go-ipfs-api"
	"log"
)

func getProviders(ctx context.Context, l shell.Logger, providers chan<- Provider) error {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return context error
			return ctx.Err()
		default:
			// Note: this one is blocking, and might stall. We should have a timeout on this
			// or something!
			log.Printf("Waiting for next messaage")
			msg, err := l.Next()
			if err != nil {
				return err
			}

			provider, err := Message(msg).ResourceProvider()
			if err != nil {
				return err
			}

			if provider != nil {
				providers <- *provider
			}
		}
	}
}

func addProviders(ctx context.Context, providers <-chan Provider, queue *queue.Queue) error {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return context error
			return ctx.Err()
		case p := <-providers:
			// Add with highest priority, as this is supposed to be available
			log.Printf("Queueing %v", p.Resource)

			err := queue.Publish(&crawler.Args{
				Hash: p.Id,
			}, 9)

			if err != nil {
				return err
			}
		}
	}
}

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
