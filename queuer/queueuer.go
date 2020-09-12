package queuer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	t "github.com/ipfs-search/ipfs-search/types"
	"log"
)

// Queue allows publishing of sniffed items.
type Queue interface {
	Publish(interface{}, uint8) error
}

type Queuer struct {
	queue     Queue
	providers <-chan t.Provider
}

func New(q Queue, providers <-chan t.Provider) Queuer {
	return Queuer{
		queue:     q,
		providers: providers,
	}
}

func (q *Queuer) Queue(ctx context.Context) error {
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
