package queue

import (
	"context"
	"log"
)

// Worker instantiates and calls MessageWorker for every Message in Queue
type Worker struct {
	errChan chan<- error
	queue   *Queue
	factory MessageWorkerFactory
}

// NewWorker returns a worker for a given queue with error channel. The
// MessageWorkerFactory is itself wrapped in a messageWorker for proper
// error handling etc.
func NewWorker(errc chan<- error, queue *Queue, factory MessageWorkerFactory) *Worker {
	return &Worker{
		errChan: errc,
		queue:   queue,
		factory: newMessageWorker(factory),
	}
}

// String returns the name of the worker queue
func (w *Worker) String() string {
	return w.queue.String()
}

// Work performs consumption of messages in the worker's Queue
func (w *Worker) Work(ctx context.Context) error {
	msgs, err := w.queue.Consume()
	if err != nil {
		return err
	}

	// Keep consuming messages until context is cancelled
	for {
		select {
		case <-ctx.Done():
			// Context canceled, stop processing messages
			log.Printf("Stopping worker %s: %s", w, ctx.Err())
			return ctx.Err()
		case msg := <-msgs:
			worker := w.factory(&msg)
			err = worker.Work(ctx)
			if err != nil {
				w.errChan <- err
			}
		}
	}
}
