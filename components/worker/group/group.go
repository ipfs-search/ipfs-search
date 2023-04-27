package group

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"
	samqp "github.com/rabbitmq/amqp091-go"

	"github.com/ipfs-search/ipfs-search/components/worker"
	"github.com/ipfs-search/ipfs-search/instr"
)

// WorkerGetter returns a named worker.
type WorkerGetter func(name string) *worker.Worker

// Group wraps a group of workers in errgroup.
type Group struct {
	errg *errgroup.Group
	ctx context.Context

	workerGetter WorkerGetter

	*instr.Instrumentation
}

// New returns a new workergroup which will instantiate workers using workgerGetter.
func New(ctx context.Context, workerGetter WorkerGetter, i *instr.Instrumentation) *Group {
	errg, ctx := errgroup.WithContext(ctx)
	return &Group{
		errg, ctx, workerGetter, i,
	}
}

// Go starts a pool with a number of workers for delivieries and a name.
func (g *Group) Go(deliveries <-chan samqp.Delivery, workers int, poolName string) {
	ctx, span := g.Tracer.Start(g.ctx, "crawler.pool.start")
	defer span.End()

	log.Printf("Starting %d workers for %s", workers, poolName)

	for i := 0; i < workers; i++ {
		name := fmt.Sprintf("%s-%d", poolName, i)
		worker := g.workerGetter(name)

		g.errg.Go(func() error {
			return worker.Start(ctx, deliveries)
		})
	}
}

// Wait wraps the underlying errorgroup's wait.
func (g *Group) Wait() error {
	return g.errg.Wait()
}
