package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/ipfs-search/ipfs-search/worker"
	"github.com/streadway/amqp"
	"log"
)

// MessageWorkerFactory instantiates a worker for a single AMQP message
type MessageWorkerFactory func(msg *amqp.Delivery) worker.Worker

// newMessageWorker implements MessageWorkerFactory and wraps a factory with
// a messageWorker, such that messages will be properly acked/rejected and
// errors/panics handled
func newMessageWorker(factory MessageWorkerFactory) MessageWorkerFactory {
	return func(msg *amqp.Delivery) worker.Worker {
		return &messageWorker{
			Factory:  factory,
			Delivery: msg,
		}
	}
}

// messageWorker instantiates and wraps a single worker for every message for
// error handling and ack/rejection
type messageWorker struct {
	Factory MessageWorkerFactory
	*amqp.Delivery
}

// Work initiates the contained worker for a single message, acking if no error and rejecting otherwise
func (m *messageWorker) Work(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Override original error value on panic
			err = m.recoverPanic(r)
		}
	}()

	log.Printf("Received: %s", m.Body)

	// Create new worker for the actual work and perform it
	worker := m.Factory(m.Delivery)
	err = worker.Work(ctx)

	if err != nil {
		// Don't retry
		m.Reject(false)

		return
	}

	// Everything went fine, ack the message
	m.Ack(false)

	return
}

func (m *messageWorker) recoverPanic(r interface{}) (err error) {
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
