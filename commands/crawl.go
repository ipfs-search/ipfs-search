package commands

import (
	"context"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler/factory"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/worker"
	"golang.org/x/sync/errgroup"
	"log"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
)

// Crawl configures and initializes crawling
func Crawl(ctx context.Context, cfg *config.Config) error {
	instrumentation := instr.New()
	tracer := instrumentation.Tracer

	ctx, span := tracer.Start(ctx, "commands.Crawl")
	defer span.End()

	errc := make(chan error, 1)

	startWorkers := func(ctx context.Context, cfg *config.Config, errc chan<- error) (*errgroup.Group, error) {
		ctx, span := tracer.Start(ctx, "commands.startWorkers")
		defer span.End()

		factory, err := factory.New(ctx, cfg.FactoryConfig(), instrumentation, errc)
		if err != nil {
			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
			return nil, err
		}

		hashGroup := worker.Group{
			Count:   cfg.Crawler.HashWorkers,
			Wait:    cfg.Crawler.HashWait,
			Factory: factory.NewHashWorker,
		}
		fileGroup := worker.Group{
			Count:   cfg.Crawler.FileWorkers,
			Wait:    cfg.Crawler.FileWait,
			Factory: factory.NewFileWorker,
		}

		// Create error group and context
		errg, ctx := errgroup.WithContext(ctx)

		// Start work loop
		errg.Go(func() error {
			ctx, span := tracer.Start(ctx, "commands.hashWorkers")
			defer span.End()
			return hashGroup.Work(ctx)
		})
		errg.Go(func() error {
			ctx, span := tracer.Start(ctx, "commands.fileWorkers")
			defer span.End()
			return fileGroup.Work(ctx)
		})

		return errg, nil
	}

	errg, err := startWorkers(ctx, cfg, errc)
	if err != nil {
		return err
	}

	log.Printf("Workers started")
	span.AddEvent(ctx, "workers-started")

	// Log errors, wait for context to cancel
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))

			log.Printf("Shutting down: %s", err)
			log.Print("Waiting for processes to finish")

			err = errg.Wait()
			log.Printf("Error group finished: %s", err)
			return err
		case err := <-errc:
			// Log errors
			log.Printf("%T: %v", err, err)
		}
	}
}
