package commands

import (
	"context"

	"github.com/ipfs-search/ipfs-search/components/crawler/worker"
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

	c, err := worker.NewPool(ctx, cfg, i)
	if err != nil {
		return err
	}

	c.Start(ctx)

	// Context closure or panic is the only way to stop crawling
	<-ctx.Done()

	return ctx.Err()
}
