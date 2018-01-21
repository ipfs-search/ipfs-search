package queue

import (
	"context"
)

// WorkerFunc processes queueue messages
type WorkerFunc func(ctx context.Context, msg *MessageWorker) error

// Worker calls Func for every message in Queue, returning errors in ErrChan
type Worker struct {
	ErrChan chan<- error
	Func    WorkerFunc
	Queue   *Queue
}

// String returns the name of the worker queue
func (w *Worker) String() string {
	return w.Queue.String()
}
