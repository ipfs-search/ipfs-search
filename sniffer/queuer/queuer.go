package queuer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/queue"
	t "github.com/ipfs-search/ipfs-search/types"
	"log"
	"time"
)

const queueTimeout = 5 * time.Minute

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
	var err error

	for {
		// Never wait more than queueTimeout for a message
		ctx, cancel := context.WithTimeout(ctx, queueTimeout)

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case p := <-q.providers:
			log.Printf("Queueing %v", p.Resource)

			// Add with highest priority (9), as this is supposed to be available
			err = q.queue.Publish(&crawler.Args{
				Hash: p.ID,
			}, 9)
		}

		cancel()

		if err != nil {
			return err
		}
	}
}
