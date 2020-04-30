package sniffer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/queue"
	t "github.com/ipfs-search/ipfs-search/types"
	"log"
)

func addProviders(ctx context.Context, providers <-chan t.Provider, queue *queue.Queue) error {
	for {
		select {
		case <-ctx.Done():
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
