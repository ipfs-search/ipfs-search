package worker

import (
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

// Factory returns a worker
type Factory func(context.Context) (Worker, error)

// Group represents a group of Count workers, created by Factory
type Group struct {
	Count   uint
	Factory Factory
	Wait    time.Duration // Time to wait between starting workers
}

// Work starts Count of workers, created by Factory
func (g *Group) Work(ctx context.Context) error {
	// Create error group and context
	errg, ctx := errgroup.WithContext(ctx)

	// Create a pool of workers within errorgroup
	for i := uint(0); i < g.Count; i++ {
		worker, err := g.Factory(ctx)
		if err != nil {
			return err
		}

		errg.Go(func() error {
			log.Printf("Starting worker %s (%d)", worker, i+1)
			return worker.Work(ctx)
		})

		time.Sleep(g.Wait)
	}

	// Block until done, returning an error if and as soon as one of the
	// child contexts errors
	return errg.Wait()
}
