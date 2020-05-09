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

	// Create and configure add queue
	conn, err := queue.NewConnection(cfg.AMQPURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Yielded hashes (of which type is unknown), should be added to hashes.
	q, err := conn.NewChannelQueue("hashes")
	if err != nil {
		return err
	}

	s, err := sniffer.New(cfg.SnifferConfig(), sh, q)
	if err != nil {
		// Error starting sniffer
		return err
	}

	for {
		err := s.Sniff(ctx)
		log.Printf("Sniffer completed, error: %v", err)

		select {
		case <-ctx.Done():
			// Context cancelled, return error
			return err
		case <-time.After(1 * time.Second):
			// Wait a second, preventing tight restart loops
			log.Printf("Restarting sniffer, error: %v", err)
		}
	}
}
