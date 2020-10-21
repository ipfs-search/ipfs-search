package queue

import (
	"context"
	"fmt"
	"github.com/ipfs-search/ipfs-search/instr"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"log"
)

// Worker instantiates and calls MessageWorker for every Message in Queue
type Worker struct {
	errChan chan<- error
	queue   Consumer
	factory MessageWorkerFactory
	*instr.Instrumentation
}

// NewWorker returns a worker for a given queue with error channel. The
// MessageWorkerFactory is itself wrapped in a messageWorker for proper
// error handling etc.
func NewWorker(errc chan<- error, queue Consumer, factory MessageWorkerFactory, i *instr.Instrumentation) *Worker {
	return &Worker{
		errChan:         errc,
		queue:           queue,
		factory:         newMessageWorker(factory, i),
		Instrumentation: i,
	}
}

// String returns the name of the worker queue
func (w *Worker) String() string {
	return fmt.Sprintf("%s", w.queue)
}

// Work performs consumption of messages in the worker's Queue
func (w *Worker) Work(ctx context.Context) error {
	ctx, span := w.Tracer.Start(ctx, "queue.Worker.Work")
	defer span.End()

	msgs, err := w.queue.Consume(ctx)
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}

	// Keep consuming messages until context is cancelled
	for {
		select {
		case <-ctx.Done():
			// Context canceled, stop processing messages
			log.Printf("Stopping worker %s: %s", w, ctx.Err())
			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
			return ctx.Err()
		case msg := <-msgs:
			msgWorker := w.factory(&msg)
			err = msgWorker.Work(ctx)
			if err != nil {
				span.RecordError(ctx, err, trace.WithErrorStatus(codes.Ok))
				w.errChan <- err
			}
		}
	}
}
