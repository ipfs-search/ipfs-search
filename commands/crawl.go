package commands

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawlworker"
	"github.com/ipfs-search/ipfs-search/worker"
	"golang.org/x/sync/errgroup"
	"log"
)

// Crawl configures and initializes crawling
func Crawl() error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	errc := make(chan error, 1)

	factory, err := crawlworker.NewFactory(config, errc)
	if err != nil {
		return err
	}

	hashGroup := worker.Group{
		Count:   config.HashWorkers,
		Factory: factory.NewHashWorker,
	}
	fileGroup := worker.Group{
		Count:   config.FileWorkers,
		Factory: factory.NewFileWorker,
	}

	// Create error group and context
	errg, ctx := errgroup.WithContext(context.Background())

	// Start work loop
	errg.Go(func() error { return hashGroup.Work(ctx) })
	errg.Go(func() error { return fileGroup.Work(ctx) })

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// Start block on errc, logging messages
mainloop:
	for {
		select {
		// TODO: Cancel on QUIT signal
		case <-ctx.Done():
			// Context canceled, stop
			log.Printf("Context cancelled: %s", ctx.Err())
			break mainloop
		case err = <-errc:
			// Print errors
			log.Printf("%T: %v", err, err)
		}
	}

	// Wait until all processes have finished
	err = errg.Wait()
	log.Printf("Error group finished: %s", err)
	return err
}
