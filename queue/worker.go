package queue

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// WorkerFunc processes queueue messages
type WorkerFunc func(ctx context.Context, msg *amqp.Delivery) error

// Worker calls Func for every message in Queue, returning errors in ErrChan
type Worker struct {
	ErrChan chan error
	Func    WorkerFunc
	Queue   *Queue
}

// messagePanic handles panic in a single message
func (w *Worker) messagePanic(msg *amqp.Delivery) {
	if r := recover(); r != nil {
		log.Printf("Panic in: %s", msg.Body)

		// Permanently remove msg from original queue
		msg.Reject(false)

		err, ok := r.(error)

		if !ok {
			err = fmt.Errorf("Unassertable panic error: %v", r)
		}

		w.ErrChan <- err
	}

}

// procesMessage processes a single message
func (w *Worker) processMessage(ctx context.Context, msg *amqp.Delivery) (err error) {
	defer w.messagePanic(msg)

	log.Printf("Received a msg: %s", msg.Body)

	err = w.Func(ctx, msg)

	if err != nil {
		// Don't retry
		msg.Reject(false)

		return
	}

	// Everything went fine, ack the msg
	msg.Ack(false)

	return
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
			return ctx.Err()
		case msg := <-msgs:
			w.ErrChan <- w.processMessage(ctx, &msg)
		}
	}
}
