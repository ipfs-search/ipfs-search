package commands

import (
	"context"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/sniffer"
	"log"
	"time"
)

// Sniff configures and initializes crawling
func Sniff(ctx context.Context, cfg *config.Config) error {
	s, err := sniffer.New(cfg.SnifferConfig())
	if err != nil {
		// Error starting sniffer
		return err
	}

	for {
		err := s.Work(ctx)
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
