package worker

import (
	"context"
	"golang.org/x/sync/errgroup"
)

// Group represents a group of Count workers, created by Factory
type Group struct {
	Count   uint
	Factory Factory
}

// Work starts Count of workers, created by Factory
func (g *Group) Work(ctx context.Context) error {
	var (
		worker Worker
		err    error
	)

	// Create a pool of workers
	for i := uint(0); i < g.Count; i++ {
		worker, err = g.Factory()
		if err != nil {
			return err
		}
	}

	// Create error group and context
	errg, ctx := errgroup.WithContext(ctx)

	// Start the workers, passing them the error group's context
	// This way, if one of the workers returns an error, the Done channel
	// is closed and they'll all stop and they can be signalled to stop
	// by cancelling the parent context.
	for i := uint(0); i < g.Count; i++ {
		errg.Go(func() error {
			return worker.Work(ctx)
		})
	}

	// Block until done, returning an error if and as soon as one of the
	// child contexts errors
	return errg.Wait()
}
