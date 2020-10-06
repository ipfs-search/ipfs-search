package amqp

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

// Queue wraps an channel/queue for tasks
type Queue struct {
	name    string
	channel *Channel
}

// String returns the name of the queue
func (q *Queue) String() string {
	return q.name
}

// Publish adds a task with specified params to the Queue
// priority: higher number, higher priority
func (q *Queue) Publish(params interface{}, priority uint8) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return q.channel.ch.Publish(
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
}

// Consume consumes messages from a queue
func (q *Queue) Consume() (<-chan amqp.Delivery, error) {
	return q.channel.ch.Consume(
		q.name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}
