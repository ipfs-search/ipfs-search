package amqp

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ipfs-search/ipfs-search/instr"
)

// Channel wraps an AMQP channel
type Channel struct {
	ch *amqp.Channel
	*instr.Instrumentation
	MessageTTL time.Duration
}

// Queue creates a named queue on a given chennel
func (c *Channel) Queue(ctx context.Context, name string) (*Queue, error) {
	ctx, span := c.Tracer.Start(ctx, "queue.amqp.Channel.Queue", trace.WithAttributes(attribute.String("queue", name)))
	defer span.End()

	_, err := c.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-max-priority": 9, // Enable all 9 priorities
			"x-message-ttl":  c.MessageTTL.Milliseconds(),
			"x-queue-mode":   "lazy", // Allow RabbitMQ to write queue to disk as fast as possible
		},
	)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &Queue{
		channel:         c,
		name:            name,
		Instrumentation: c.Instrumentation,
	}, nil
}

// Close closes a Channel
func (c *Channel) Close() error {
	return c.ch.Close()
}
