package amqp

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

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
		trace.WithAttributes(attribute.String("amqp_url", f.Config.URL)),
		trace.WithAttributes(attribute.String("queue", f.Queue)),
	)
	defer span.End()

	// Create and configure add queue
	conn, err := NewConnection(ctx, f.Config, f.AMQPConfig, f.Instrumentation)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Close connection when context closes
	go func() {
		<-ctx.Done()
		span.AddEvent("closing-amqp-context-closed")
		log.Printf("Closing AMQP connection; context closed")
		conn.Close()
	}()

	return conn.NewChannelQueue(ctx, f.Queue, 1)
}
