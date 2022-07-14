// Package queue provides queueing semantics.
package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher allows publishing of sniffed items.
type Publisher interface {
	Publish(context.Context, interface{}, uint8) error
}

// Consumer allows consuming of published items.
type Consumer interface {
	Consume(context.Context) (<-chan amqp.Delivery, error)
}

// PublisherFactory creates Publishers.
type PublisherFactory interface {
	NewPublisher(context.Context) (Publisher, error)
}

// Queue allows both publishing as well as consuming.
type Queue interface {
	Publisher
	Consumer
}
