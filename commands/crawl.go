package commands

import (
	"context"

	"github.com/ipfs-search/ipfs-search/components/worker/pool"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/instr"

	"log"
)

// Crawl configures and initializes crawling
func Crawl(ctx context.Context, cfg *config.Config) error {
	instFlusher, err := instr.Install(cfg.InstrConfig(), "ipfs-crawler")
	if err != nil {
		log.Fatal(err)
	}
	defer instFlusher(ctx)

	i := instr.New()

	ctx, span := i.Tracer.Start(ctx, "commands.Crawl")
	defer span.End()

	pool := pool.New(cfg, i)

	if err := pool.Init(ctx); err != nil {
		return err
	}

	return pool.Start(ctx)
}
