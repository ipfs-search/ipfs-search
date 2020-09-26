package sniffer

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"time"

	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs-search/ipfs-search/sniffer/eventsource"
	"github.com/ipfs-search/ipfs-search/sniffer/handler"
	filters "github.com/ipfs-search/ipfs-search/sniffer/providerfilters"
	"github.com/ipfs-search/ipfs-search/sniffer/queuer"
	filter "github.com/ipfs-search/ipfs-search/sniffer/streamfilter"
	t "github.com/ipfs-search/ipfs-search/types"

	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-eventbus"
)

// Sniffer allows sniffing Batching datastore's events, effectively allowing sniffing of the IPFS DHT.
type Sniffer struct {
	cfg *Config
	es  eventsource.EventSource
	pub queue.PublisherFactory
}

// New creates a new Sniffer or returns an error.
func New(cfg *Config, ds datastore.Batching, pub queue.PublisherFactory) (*Sniffer, error) {
	bus := eventbus.NewBus()

	es, err := eventsource.New(bus, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to get eventsource: %w", err)
	}

	s := Sniffer{
		cfg: cfg,
		es:  es,
		pub: pub,
	}

	return &s, nil
}

// Batching returns the datastore wrapped with sniffing hooks.
func (s *Sniffer) Batching() datastore.Batching {
	return s.es.Batching()
}

// Sniff starts sniffing until the context is closed - it restarts itself on intermittant errors.
func (s *Sniffer) Sniff(ctx context.Context) error {
	sniffedProviders := make(chan t.Provider, s.cfg.BufferSize)
	filteredProviders := make(chan t.Provider, s.cfg.BufferSize)

	// Create error group and context
	for {
		errg, errCtx := errgroup.WithContext(ctx)
		errg.Go(func() error {
			h := handler.New(sniffedProviders)

			return s.es.Subscribe(errCtx, h.HandleFunc)
		})
		errg.Go(func() error {
			lastSeenFilter := filters.NewLastSeenFilter(s.cfg.LastSeenExpiration, s.cfg.LastSeenPruneLen)
			cidFilter := filters.NewCidFilter()
			mutliFilter := filters.NewMultiFilter(lastSeenFilter, cidFilter)
			f := filter.New(mutliFilter, sniffedProviders, filteredProviders)

			return f.Filter(errCtx)
		})
		errg.Go(func() error {
			publisher, err := s.pub.NewPublisher(ctx)
			if err != nil {
				return err
			}

			q := queuer.New(publisher, filteredProviders)

			return q.Queue(errCtx)
		})

		// Wait until context lose
		err := errg.Wait()

		// Closing the parent context should cause a return.
		if err := ctx.Err(); err != nil {
			log.Printf("Parent context closed with error '%s', returning error", err)
			return err
		}

		log.Printf("Wait group exited with error '%s', restarting", err)

		// TODO: Add circuit breaker here
		log.Printf("Stubbornly restarting in 1s")
		time.Sleep(time.Second)
	}
}
