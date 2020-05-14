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
	cfg     *Config
	yielder *providerYielder
	filter  *providerFilter
	queuer  *providerQueuer
}

// New returns a new sniffer.
func New(cfg *Config) (*Sniffer, error) {
	// Initialize yielder
	x, err := extractor.New()
	if err != nil {
		return nil, err
	}
	py := &providerYielder{e: x, timeout: cfg.LoggerTimeout}

	// Initialize filter
	lastSeenFilter := filters.NewLastSeenFilter(cfg.LastSeenExpiration, cfg.LastSeenPruneLen)
	cidFilter := filters.NewCidFilter()
	pf := &providerFilter{f: filters.NewMultiFilter(lastSeenFilter, cidFilter)}

	// Initialise queuer
	pq := &providerQueuer{}

	return &Sniffer{
		cfg:     cfg,
		yielder: py,
		filter:  pf,
		queuer:  pq,
	}, nil
}

// Sniff starts sniffing, only returning in error conditions.
func (s *Sniffer) Sniff(ctx context.Context, logger Logger, queue Queue) error {
	sniffedProviders := make(chan t.Provider, s.cfg.BufferSize)
	filteredProviders := make(chan t.Provider, s.cfg.BufferSize)

	// Create error group and context
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		return s.yielder.yield(ctx, logger, sniffedProviders)
	})
	errg.Go(func() error {
		return s.filter.filter(ctx, sniffedProviders, filteredProviders)
	})
	errg.Go(func() error {
		return s.queuer.queue(ctx, filteredProviders, queue)
	})

	return errg.Wait()
}
