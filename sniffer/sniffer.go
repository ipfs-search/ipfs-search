package sniffer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/sniffer/filters"
	t "github.com/ipfs-search/ipfs-search/types"
	"golang.org/x/sync/errgroup"
)

// Sniffer is a worker sniffing for provider messages, filtering them and feeding
// them into the crawler's queue.
type Sniffer struct {
	cfg       *Config
	shell     Shell
	queue     Queue
	filter    filters.Filter
	extractor Extractor
}

// New returns a new sniffer.
func New(cfg *Config, shell Shell, queue Queue) (*Sniffer, error) {
	// Initialize filters
	lastSeenFilter := filters.LastSeenFilter(cfg.LastSeenExpiration, cfg.LastSeenPruneLen)
	cidFilter := filters.NewCidFilter()
	filter := filters.MultiFilter(lastSeenFilter, cidFilter)

	// Initialize extractor
	extractor := ProviderExtractor{}

	return &Sniffer{
		cfg:       cfg,
		shell:     shell,
		filter:    filter,
		extractor: &extractor,
	}, nil
}

// Sniff starts sniffing, returning an error when anything goes wrong
func (s *Sniffer) Sniff(ctx context.Context) error {
	// Never timeout, this is a long poll!
	s.shell.SetTimeout(0)

	// Get logger
	logger, err := s.shell.GetLogs(ctx)
	if err != nil {
		return err
	}
	defer logger.Close()

	sniffedProviders := make(chan t.Provider)
	filteredProviders := make(chan t.Provider)

	// Create error group and context
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		return getProviders(ctx, logger, s.extractor, sniffedProviders, s.cfg.LoggerTimeout)
	})
	errg.Go(func() error {
		return filterProviders(ctx, sniffedProviders, filteredProviders, s.filter)
	})
	errg.Go(func() error {
		return queueProviders(ctx, filteredProviders, s.queue)
	})

	return errg.Wait()
}
