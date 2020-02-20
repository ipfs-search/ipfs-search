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
			err := queue.Publish(&crawler.Args{
				Hash: p.Id,
			}, 9)

			if err != nil {
				return err
			}
		}
	}
}

func filterProviders(ctx context.Context, input <-chan Provider, output chan<- Provider, filters []Filter) error {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return context error
			return ctx.Err()
		case i := <-input:
			for _, f := range filters {
				if !f.Filter(i) {
					log.Printf("Disgarding provider: %v", i)
					continue
				}
			}

			output <- i
		}
	}
}
