package queue

import (
	"context"
	"github.com/streadway/amqp"
)

// Publisher allows publishing of sniffed items.
type Publisher interface {
	Publish(interface{}, uint8) error
}

// Consumer allows consuming of published items.
type Consumer interface {
	Consume() (<-chan amqp.Delivery, error)
}

// PublisherFactory creates Publishers.
type PublisherFactory interface {
	NewPublisher(context.Context) (Publisher, error)
}
