package queue

import (
	"context"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/worker"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

// MessageWorkerFactory instantiates a worker for a single AMQP message
type MessageWorkerFactory func(msg *amqp.Delivery) worker.Worker

// newMessageWorker implements MessageWorkerFactory and wraps a factory with
// a messageWorker, such that messages will be properly acked/rejected and
// errors/panics handled
func newMessageWorker(factory MessageWorkerFactory, i *instr.Instrumentation) MessageWorkerFactory {
	return func(msg *amqp.Delivery) worker.Worker {
		return &messageWorker{
			Factory:         factory,
			Delivery:        msg,
			Instrumentation: i,
		}
	}
}

// messageWorker instantiates and wraps a single worker for every message for
// error handling and ack/rejection
type messageWorker struct {
	Factory MessageWorkerFactory
	*amqp.Delivery
	*instr.Instrumentation
}

// Work initiates the contained worker for a single message, acking if no error and rejecting otherwise
func (m *messageWorker) Work(ctx context.Context) (err error) {
	ctx, span := m.Tracer.Start(ctx, "queue.messageWorker.Work",
		trace.WithAttributes(
			label.String("message", string(m.Body)),
			label.Uint("priority", uint(m.Priority)),
		),
	)
	defer span.End()

	// Create new worker for the actual work and perform it
	worker := m.Factory(m.Delivery)
	err = worker.Work(ctx)

	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))

		// Don't retry
		if err := m.Reject(false); err != nil {
			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		}

		return
	}

	span.RecordError(ctx, err, trace.WithErrorStatus(codes.Ok))

	// Everything went fine, ack the message
	if err := m.Ack(false); err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
	}

	return
}
