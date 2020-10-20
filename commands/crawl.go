package commands

import (
	"context"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler/factory"
	"github.com/ipfs-search/ipfs-search/worker"
	"golang.org/x/sync/errgroup"
	"log"
)

func startWorkers(ctx context.Context, cfg *config.Config, errc chan<- error) (*errgroup.Group, error) {
	factory, err := factory.New(ctx, cfg.FactoryConfig(), errc)
	if err != nil {
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
	errg.Go(func() error { return hashGroup.Work(ctx) })
	errg.Go(func() error { return fileGroup.Work(ctx) })

	return errg, nil
}

// Crawl configures and initializes crawling
func Crawl(ctx context.Context, cfg *config.Config) error {
	errc := make(chan error, 1)

	errg, err := startWorkers(ctx, cfg, errc)
	if err != nil {
		return err
	}

	log.Printf("Workers started")

	// Log messages, wait for context break
	go errorLoop(errc)
	err = block(ctx)

	log.Printf("Shutting down: %s", err)
	log.Print("Waiting for processes to finish")

	err = errg.Wait()
	log.Printf("Error group finished: %s", err)
	return err
}
