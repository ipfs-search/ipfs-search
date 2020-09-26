package queuer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/queue"
	t "github.com/ipfs-search/ipfs-search/types"
	"log"
)

type Queuer struct {
	queue     queue.Publisher
	providers <-chan t.Provider
}

func New(q queue.Publisher, providers <-chan t.Provider) Queuer {
	return Queuer{
		queue:     q,
		providers: providers,
	}
}

func (q *Queuer) Queue(ctx context.Context) error {
	// TODO: Consider running this in a goroutine through an errorgroup
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case p := <-q.providers:
			log.Printf("Queueing %v", p.Resource)

			// Add with highest priority (9), as this is supposed to be available
			err := q.queue.Publish(&crawler.Args{
				Hash: p.ID,
			}, 9)

			if err != nil {
				return err
			}
		}
	}
}
