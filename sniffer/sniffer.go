package sniffer

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
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

const bufSize = 256

type Sniffer struct {
	es eventsource.EventSource
	h  handler.Handler
	f  filter.Filter
	q  queuer.Queuer
}

func New(ds datastore.Batching) (*Sniffer, error) {
	// Test update
	bus := eventbus.NewBus()

	es, err := eventsource.New(bus, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to get eventsource: %w", err)
	}

	sniffedProviders := make(chan t.Provider, bufSize)
	filteredProviders := make(chan t.Provider, bufSize)

	handler := handler.New(sniffedProviders)

	lastSeenFilter := filters.NewLastSeenFilter(60*time.Duration(time.Minute), 16383)
	cidFilter := filters.NewCidFilter()
	f := filters.NewMultiFilter(lastSeenFilter, cidFilter)

	// Create and configure add queue
	conn, err := queue.NewConnection("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Yielded hashes (of which type is unknown), should be added to hashes.
	queue, err := conn.NewChannelQueue("hashes")
	if err != nil {
		return nil, err
	}

	s := Sniffer{
		es: es,
		h:  handler,
		f:  filter.New(f, sniffedProviders, filteredProviders),
		q:  queuer.New(queue, filteredProviders),
	}

	return &s, nil
}

func (s *Sniffer) Batching() datastore.Batching {
	return s.es.Batching()
}

func (s *Sniffer) Sniff(ctx context.Context) error {
	// Create error group and context
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		return s.es.Subscribe(ctx, s.h.HandleFunc)
	})
	errg.Go(func() error {
		return s.f.Filter(ctx)
	})
	errg.Go(func() error {
		return s.q.Queue(ctx)
	})

	return errg.Wait()
}
