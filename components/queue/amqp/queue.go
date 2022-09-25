package amqp

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ipfs-search/ipfs-search/components/queue"
	"github.com/ipfs-search/ipfs-search/instr"
)

// Queue wraps an channel/queue for tasks
type Queue struct {
	name    string
	channel *Channel
	*instr.Instrumentation
}

// String returns the name of the queue
func (q *Queue) String() string {
	return q.name
}

// Publish adds a task with specified params to the Queue
// priority: higher number, higher priority
// TODO: Add context parameter, allow for timeouts etc
func (q *Queue) Publish(ctx context.Context, params interface{}, priority uint8) error {
	ctx, span := q.Tracer.Start(ctx, "queue.amqp.Publish",
		trace.WithAttributes(
			attribute.String("queue", q.name),
			// attribute.Any("params", params)),
			attribute.Int("priority", int(priority))),
	)
	defer span.End()

	body, err := json.Marshal(params)
	if err != nil {
		span.RecordError(err)
		return err
	}

	err = q.channel.ch.Publish(
		"",     // exchange
		q.name, // routing key
		true,   // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "application/json",
			Body:         body,
			Priority:     priority,
		})

	if err != nil {
		span.RecordError(err)
	}

	return err
}

// Consume consumes messages from a queue
func (q *Queue) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	ctx, span := q.Tracer.Start(ctx, "queue.amqp.Consume")
	defer span.End()

	c, err := q.channel.ch.Consume(
		q.name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return c, err
}

// Compile-time assurance that implementation satisfies interface.
var _ queue.Queue = &Queue{}
