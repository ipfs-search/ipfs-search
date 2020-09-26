package sniffer

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"time"

	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs-search/ipfs-search/sniffer/filters"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs-search/ipfs-sniffer/eventsource"
	"github.com/ipfs-search/ipfs-sniffer/filter"
	"github.com/ipfs-search/ipfs-sniffer/handler"
	"github.com/ipfs-search/ipfs-sniffer/queuer"

	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-eventbus"
)

const bufSize = 512

type Sniffer struct {
	es eventsource.EventSource
}

func New(ds datastore.Batching) (*Sniffer, error) {
	// Test update
	bus := eventbus.NewBus()

	es, err := eventsource.New(bus, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to get eventsource: %w", err)
	}

	s := Sniffer{
		es: es,
	}

	return &s, nil
}

func (s *Sniffer) Batching() datastore.Batching {
	return s.es.Batching()
}

func (s *Sniffer) Sniff(ctx context.Context) {
	sniffedProviders := make(chan t.Provider, bufSize)
	filteredProviders := make(chan t.Provider, bufSize)

	// Create error group and context
	for {
		errg, ctx := errgroup.WithContext(ctx)
		errg.Go(func() error {
			h := handler.New(sniffedProviders)

			return s.es.Subscribe(ctx, h.HandleFunc)
		})
		errg.Go(func() error {
			lastSeenFilter := filters.NewLastSeenFilter(60*time.Duration(time.Minute), 32768)
			cidFilter := filters.NewCidFilter()
			mutliFilter := filters.NewMultiFilter(lastSeenFilter, cidFilter)
			f := filter.New(mutliFilter, sniffedProviders, filteredProviders)

			return f.Filter(ctx)
		})
		errg.Go(func() error {
			// Create and configure add queue
			conn, err := queue.NewConnection("amqp://guest:guest@localhost:5672/")
			if err != nil {
				return err
			}
			defer conn.Close()

			// Yielded hashes (of which type is unknown), should be added to hashes.
			queue, err := conn.NewChannelQueue("hashes")
			if err != nil {
				return err
			}

			q := queuer.New(queue, filteredProviders)

			return q.Queue(ctx)
		})

		err := errg.Wait()

		log.Printf("Wait group exited, error: %s", err)
		log.Printf("Stubbornly restarting in 1s")
		time.Sleep(time.Second)
	}

}
