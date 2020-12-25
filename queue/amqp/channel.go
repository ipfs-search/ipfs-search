package amqp

import (
	"context"
	"fmt"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

// Channel wraps an AMQP channel
type Channel struct {
	ch *amqp.Channel
	*instr.Instrumentation
}

// Queue creates a named queue on a given chennel
func (c *Channel) Queue(ctx context.Context, name string) (*Queue, error) {
	ctx, span := c.Tracer.Start(ctx, "queue.amqp.Channel.Queue", trace.WithAttributes(label.String("queue", name)))
	defer span.End()

	deadQueue := fmt.Sprintf("%s-dead", name)

	_, err := c.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-max-priority": 9,                   // Enable all 9 priorities
			"x-message-ttl":  1000 * 60 * 60 * 24, // Expire messages after 24 hours
		},
	)
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
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
