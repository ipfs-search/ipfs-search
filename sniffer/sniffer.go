package sniffer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/sniffer/extractor"
	"github.com/ipfs-search/ipfs-search/sniffer/filters"
	t "github.com/ipfs-search/ipfs-search/types"
	"golang.org/x/sync/errgroup"
)

// Sniffer is a worker sniffing for provider messages, filtering them and feeding
// them into the crawler's queue.
type Sniffer struct {
	cfg       *Config
	queue     Queue
	filter    filters.Filter
	extractor Extractor
}

// New returns a new sniffer.
func New(cfg *Config, queue Queue) (*Sniffer, error) {
	// Initialize filters
	lastSeenFilter := filters.NewLastSeenFilter(cfg.LastSeenExpiration, cfg.LastSeenPruneLen)
	cidFilter := filters.NewCidFilter()
	f := filters.NewMultiFilter(lastSeenFilter, cidFilter)

	// Initialize extractor
	x, err := extractor.New()
	if err != nil {
		return nil, err
	}

	return &Sniffer{
		cfg:       cfg,
		queue:     queue,
		filter:    f,
		extractor: x,
	}, nil
}

// Sniff starts sniffing, only returning in error conditions.
func (s *Sniffer) Sniff(ctx context.Context, logger Logger) error {
	sniffedProviders := make(chan t.Provider, s.cfg.BufferSize)
	filteredProviders := make(chan t.Provider, s.cfg.BufferSize)

	// Create error group and context
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		return yieldProviders(ctx, logger, s.extractor, sniffedProviders, s.cfg.LoggerTimeout)
	})
	errg.Go(func() error {
		return filterProviders(ctx, sniffedProviders, filteredProviders, s.filter)
	})
	errg.Go(func() error {
		return queueProviders(ctx, filteredProviders, s.queue)
	})

	return errg.Wait()
}
