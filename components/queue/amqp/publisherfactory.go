package amqp

import (
	"context"
	"log"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"

	"github.com/ipfs-search/ipfs-search/components/queue"
	"github.com/ipfs-search/ipfs-search/instr"
)

// PublisherFactory automates creation of AMQP Publishers.
type PublisherFactory struct {
	*Config
	AMQPConfig *amqp.Config
	Queue      string
	*instr.Instrumentation
}

// NewPublisher generates a new publisher or returns an error.
func (f PublisherFactory) NewPublisher(ctx context.Context) (queue.Publisher, error) {
	ctx, span := f.Tracer.Start(ctx, "queue.amqp.NewPublisher",
		trace.WithAttributes(label.String("amqp_url", f.Config.URL)),
		trace.WithAttributes(label.String("queue", f.Queue)),
	)
	defer span.End()

	// Create and configure add queue
	conn, err := NewConnection(ctx, f.Config, f.AMQPConfig, f.Instrumentation)
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return nil, err
	}

	// Close connection when context closes
	go func() {
		<-ctx.Done()
		span.AddEvent(ctx, "closing-amqp-context-closed")
		log.Printf("Closing AMQP connection; context closed")
		conn.Close()
	}()

	return conn.NewChannelQueue(ctx, f.Queue, 1)
}
