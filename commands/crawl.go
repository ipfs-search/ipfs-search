package commands

import (
	"context"

	"github.com/ipfs-search/ipfs-search/commands/crawlworker"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/instr"

	"log"
	// "go.opentelemetry.io/otel/api/trace"
	// "go.opentelemetry.io/otel/codes"
)

// Crawl configures and initializes crawling
func Crawl(ctx context.Context, cfg *config.Config) error {
	instFlusher, err := instr.Install("ipfs-crawler")
	if err != nil {
		log.Fatal(err)
	}
	defer instFlusher()

	instr := instr.New()
	tracer := instr.Tracer

	ctx, span := tracer.Start(ctx, "commands.Crawl")
	defer span.End()

	c := crawlworker.New(cfg, instr)
	c.Start(ctx)

	// Context closure or panic is the only way to stop crawling
	<-ctx.Done()

	return ctx.Err()
}
