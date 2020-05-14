package commands

import (
	"context"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs-search/ipfs-search/sniffer"
	"github.com/ipfs/go-ipfs-api"
	"log"
	"time"
)

// Sniff configures and initializes crawling
func Sniff(ctx context.Context, cfg *config.Config) error {
	// Initialize IPFS shell
	sh := shell.NewShell(cfg.IpfsAPI)

	// Never timeout, the logger does a long poll!
	sh.SetTimeout(0)

	// Create and configure add queue
	conn, err := queue.NewConnection(cfg.AMQPURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Yielded hashes (of which type is unknown), should be added to hashes.
	queue, err := conn.NewChannelQueue("hashes")
	if err != nil {
		return err
	}

	s, err := sniffer.New(cfg.SnifferConfig())
	if err != nil {
		// Error starting sniffer
		return err
	}

	for {
		// Get a new logger everytime; the logger tends to hang, some times
		logger, err := sh.GetLogs(ctx)
		if err != nil {
			// Error opening logger
			return err
		}

		err = s.Sniff(ctx, logger, queue)
		log.Printf("Sniffer completed, error: %v", err)

		// We're done with the current logger
		logger.Close()

		select {
		case <-ctx.Done():
			// Context cancelled from above, return error
			return err
		case <-time.After(1 * time.Second):
			// Wait a second, preventing tight restart loops
			log.Printf("Restarting sniffer, error: %v", err)
		}
	}
}
