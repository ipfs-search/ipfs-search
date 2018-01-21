package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

import (
	"github.com/ipfs-search/ipfs-search/worker"
)

// MessageWorkerFactory instantiates a worker for an AMQP message
type MessageWorkerFactory func(msg *amqp.Delivery) *Worker

// MessageWorker contains a worker performing the actual action and a delivery
// with a single AMQP message.
type MessageWorker struct {
	Worker worker.Worker
	*amqp.Delivery
}

// Work initiates the contained worker for a single message, acking if no error and rejecting otherwise
func (m *MessageWorker) Work(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Override original error value on panic
			err = m.recoverPanic(r)
		}
	}()

	log.Printf("Received: %s", m.Body)

	err = m.Worker.Work(ctx)

	if err != nil {
		// Don't retry
		m.Reject(false)

		return
	}

	// Everything went fine, ack the message
	m.Ack(false)

	return
}

func (m *MessageWorker) recoverPanic(r interface{}) (err error) {
	log.Printf("Panic in: %s", m.Body)

	// Permanently remove message from original queue
	m.Reject(false)

	// find out exactly what the error was and set err
	switch x := r.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = fmt.Errorf("Unassertable panic error: %v", r)
	}

	return
}
