package queue

import (
	"context"
	"log"
)

// Worker instantiates and calls MessageWorker for every Message in Queue
type Worker struct {
	ErrChan chan<- error
	Queue   *Queue
	Factory MessageWorkerFactory
}

// String returns the name of the worker queue
func (w *Worker) String() string {
	return w.Queue.String()
}

// Work performs consumption of messages in the worker's Queue
func (w *Worker) Work(ctx context.Context) error {
	msgs, err := w.Queue.Consume()
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
			worker := w.Factory(&msg)
			err = worker.Work(ctx)
			if err != nil {
				w.ErrChan <- err
			}
		}
	}
}
