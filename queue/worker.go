package queue

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// WorkerMessage wraps amqp delivery
type WorkerMessage struct {
	*Worker
	*amqp.Delivery
}

// WorkerFunc processes queueue messages
type WorkerFunc func(ctx context.Context, msg *WorkerMessage) error

// Worker calls Func for every message in Queue, returning errors in ErrChan
type Worker struct {
	ErrChan chan error
	Func    WorkerFunc
	Queue   *Queue
}

// Process handles a single message, acking if no error and rejecting otherwise
func (m *WorkerMessage) Process(ctx context.Context) (err error) {
	defer func() { err = m.recoverPanic() }()

	log.Printf("Received a msg: %s", m.Body)

	err = m.Worker.Func(ctx, m)

	if err != nil {
		// Don't retry
		m.Reject(false)

		return
	}

	// Everything went fine, ack the message
	m.Ack(false)

	return
}

func (m *WorkerMessage) recoverPanic() error {
	if r := recover(); r != nil {
		log.Printf("Panic in: %s", m.Body)

		// Permanently remove message from original queue
		m.Reject(false)

		err, ok := r.(error)

		if !ok {
			err = fmt.Errorf("Unassertable panic error: %v", r)
		}

		return err
	}

	return nil
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
			return ctx.Err()
		case msg := <-msgs:
			// Keep going on forever
			message := &WorkerMessage{
				Worker:   w,
				Delivery: &msg,
			}
			w.ErrChan <- message.Process(ctx)
		}
	}
}
