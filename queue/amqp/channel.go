package amqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// Channel wraps an AMQP channel
type Channel struct {
	ch *amqp.Channel
}

// Queue creates a named queue on a given chennel
func (c *Channel) Queue(name string) (*Queue, error) {
	log.Printf("Creating AMQP queue '%s'", name)

	deadQueue := fmt.Sprintf("%s-dead", name)

	_, err := c.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-max-priority":            9,                   // Enable all 9 priorities
			"x-message-ttl":             1000 * 60 * 60 * 24, // Expire messages after 24 hours
			"x-dead-letter-exchange":    "",                  // Anything failing or expiring goes here
			"x-dead-letter-routing-key": deadQueue,
			"x-max-length":              100 * 1000, // Max 100.000 messages; prevent resource exhaustion in queue.
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = c.ch.QueueDeclare(
		deadQueue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-message-ttl": 1000 * 60 * 60 * 24, // Expire messages after 24 hours,
			"x-max-length":  100 * 1000,          // Max 100.000 messages; prevent resource exhaustion in queue.
		},
	)
	if err != nil {
		return nil, err
	}

	return &Queue{
		channel: c,
		name:    name,
	}, nil
}

// Close closes a Channel
func (c *Channel) Close() error {
	return c.ch.Close()
}
