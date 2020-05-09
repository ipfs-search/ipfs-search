package sniffer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	t "github.com/ipfs-search/ipfs-search/types"
	"log"
)

func queueProviders(ctx context.Context, providers <-chan t.Provider, queue Queue) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case p := <-providers:
			log.Printf("Queueing %v", p.Resource)

			// Add with highest priority (9), as this is supposed to be available
			err := queue.Publish(&crawler.Args{
				Hash: p.Id,
			}, 9)

			if err != nil {
				return err
			}
		}
	}
}
