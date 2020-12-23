package commands

// import (
// 	"context"
// 	"github.com/ipfs-search/ipfs-search/config"
// 	"github.com/ipfs-search/ipfs-search/crawler/factory"
// 	"github.com/ipfs-search/ipfs-search/instr"
// 	"github.com/ipfs-search/ipfs-search/worker"
// 	"golang.org/x/sync/errgroup"
// 	"log"

// 	"go.opentelemetry.io/otel/api/trace"
// 	"go.opentelemetry.io/otel/codes"
// )

// // Crawl configures and initializes crawling
// func Crawl(ctx context.Context, cfg *config.Config) error {
// 	instFlusher, err := instr.Install("ipfs-crawler")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer instFlusher()

// 	instrumentation := instr.New()
// 	tracer := instrumentation.Tracer

// 	ctx, span := tracer.Start(ctx, "commands.Crawl")
// 	defer span.End()

// 	errc := make(chan error, 1)

// 	// Create error group and context
// 	// The derived Context is canceled the first time a function passed to Go returns a non-nil error or the
// 	// first time Wait returns, whichever occurs first.
// 	errg, ctx := errgroup.WithContext(ctx)

// 	startWorkers := func(ctx context.Context, cfg *config.Config, errc chan<- error) error {
// 		ctx, span := tracer.Start(ctx, "commands.startWorkers")
// 		defer span.End()

// 		factory, err := factory.New(ctx, cfg.FactoryConfig(), instrumentation, errc)
// 		if err != nil {
// 			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
// 			return err
// 		}

// 		hashGroup := worker.Group{
// 			Count:   cfg.Crawler.HashWorkers,
// 			Wait:    cfg.Crawler.HashWait,
// 			Factory: factory.NewHashWorker,
// 		}
// 		fileGroup := worker.Group{
// 			Count:   cfg.Crawler.FileWorkers,
// 			Wait:    cfg.Crawler.FileWait,
// 			Factory: factory.NewFileWorker,
// 		}

// 		// Start work loop
// 		errg.Go(func() error {
// 			ctx, span := tracer.Start(ctx, "commands.hashWorkers")
// 			defer span.End()
// 			err := hashGroup.Work(ctx)
// 			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
// 			return err

// 		})
// 		errg.Go(func() error {
// 			ctx, span := tracer.Start(ctx, "commands.fileWorkers")
// 			defer span.End()
// 			err := fileGroup.Work(ctx)
// 			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
// 			return err
// 		})

// 		return nil
// 	}

// 	if err := startWorkers(ctx, cfg, errc); err != nil {
// 		return err
// 	}

// 	log.Printf("Workers started")
// 	span.AddEvent(ctx, "workers-started")

// 	// Log errors, wait for context to cancel
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			// Errorgroup context closed (parent or error ocurred).
// 			err := ctx.Err()
// 			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))

// 			log.Printf("Shutting down: %s", err)
// 			log.Print("Waiting for workers to finish")

// 			// Wait blocks until all function calls from the Go method
// 			// have returned, then returns the first non-nil error (if any) from them.
// 			err = errg.Wait()
// 			log.Printf("Error group finished: %s", err)
// 			return err
// 		case err := <-errc:
// 			// Log errors
// 			log.Printf("%T: %v", err, err)
// 		}
// 	}
// }
