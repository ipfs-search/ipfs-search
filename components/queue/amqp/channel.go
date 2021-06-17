package amqp

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"

	"github.com/ipfs-search/ipfs-search/instr"
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

	// Declare DLQ
	_, err := c.ch.QueueDeclare(
		deadQueue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
	)
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return nil, err
	}
	// TODO: Move following to documentation.
	// In production, we're using a continuous shovel moving items from the DLQ to the normal queue with lowest possible
	// priority.
	// Ref:
	// PUT {{ _.base_url }}api/parameters/shovel/%2f/ {{ _.queue }}-shovel
	// {
	//   "value": {
	//     "src-protocol": "amqp091",
	//     "src-uri": "amqp://localhost",
	//     "src-queue": "{{ _.queue }}-dead",
	//     "dest-protocol": "amqp091",
	//     "dest-uri": "amqp://localhost",
	//     "dest-queue": "{{ _.queue }}",
	// 		 "dest-publish-properties": {
	// 			 "priority": 1
	// 		 },
	// 		 "delete-after": "never"
	//   }
	// }

	// Declare main qeueue
	_, err = c.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-max-priority":            9,             // Enable all 9 priorities
			"x-message-ttl":             1000 * 60 * 5, // Expire messages after 5 min
			"x-queue-mode":              "lazy",        // Allow RabbitMQ to write queue to disk as fast as possible
			"x-dead-letter-exchange":    "",            // Anything failing or expiring goes here
			"x-dead-letter-routing-key": deadQueue,
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
