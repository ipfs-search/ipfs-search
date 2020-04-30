package sniffer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs/go-ipfs-api"
	"golang.org/x/sync/errgroup"
)

// Sniffer is a worker sniffing for provider messages, filtering them and feeding
// them into the crawler's queue.
type Sniffer struct {
	Shell  *shell.Shell
	Config *Config
}

// New returns a new sniffer.
func New(cfg *Config) (*Sniffer, error) {
	// Create and configure Ipfs shell
	sh := shell.NewShell(cfg.IpfsAPI)

	// Never timeout, this is a long poll!
	sh.SetTimeout(0)

	return &Sniffer{
		Shell:  sh,
		Config: cfg,
	}, nil
}

// Work starts a blocking sniffer, returning an error when anything goes wrong
func (s *Sniffer) Work(ctx context.Context) error {
	logger, err := s.Shell.GetLogs(ctx)
	if err != nil {
		return err
	}

	defer logger.Close()

	// Create and configure add queue
	conn, err := queue.NewConnection(s.Config.AMQPURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	queue, err := conn.NewChannelQueue("hashes")
	if err != nil {
		return err
	}

	sniffedProviders := make(chan Provider)
	filteredProviders := make(chan Provider)

	lastSeenFilter := NewLastSeenFilter(s.Config.LastSeenExpiration, s.Config.LastSeenPruneLen)
	cidFilter := NewCidFilter()
	filters := []Filter{lastSeenFilter, cidFilter}

	// Create error group and context
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error { return getProviders(ctx, logger, sniffedProviders, s.Config.LoggerTimeout) })
	errg.Go(func() error { return filterProviders(ctx, sniffedProviders, filteredProviders, filters) })
	errg.Go(func() error { return addProviders(ctx, filteredProviders, queue) })

	return errg.Wait()
}
