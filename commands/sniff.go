package commands

import (
	"context"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/sniffer"
	"log"
)

// Sniff configures and initializes crawling
func Sniff(ctx context.Context, cfg *config.Config) error {
	s, err := sniffer.New(cfg.SnifferConfig())
	if err != nil {
		return err
	}

	log.Printf("Starting sniffer")

	// TODO: Implement exponential backoff etc here (and elsewhere).
	for {
		err := s.Work(ctx)

		log.Printf("Restarting sniffer, error: %v", err)
	}
}
