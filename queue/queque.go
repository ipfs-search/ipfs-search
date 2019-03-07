package queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

// Connection wraps an AMQP connection
type Connection struct {
	*amqp.Connection
}

// NewConnection returns new AMQP connection
func NewConnection(url string) (*Connection, error) {
	connection, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	return &Connection{Connection: connection}, nil
}

// NewChannel initialises an AMQP channel
func (conn *Connection) NewChannel() (*Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	return &Channel{
		Channel: ch,
	}, nil
}

// Channel wraps an AMQP channel
type Channel struct {
	*amqp.Channel
}

// Close closes a Channel
func (c *Channel) Close() error {
	err := c.Close()
	if err != nil {
		return err
	}

	return nil
}

// Queue wraps an channel/queue for tasks
type Queue struct {
	Channel *Channel
	*amqp.Queue
}

// String returns the name of the queue
func (q *Queue) String() string {
	return q.Name
}

// NewQueue creates a named queue on a given chennel
func (c *Channel) NewQueue(name string) (*Queue, error) {
	args := amqp.Table{
		"x-max-priority": 9,                   // Enable all 9 priorities
		"x-message-ttl":  1000 * 60 * 60 * 24, // Expire messages after 24 hours
	}

	q, err := c.Channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Queue{
		Channel: c,
		Queue:   &q,
	}, nil
}

// NewChannelQueue returns a new queue on a new channel
func (conn *Connection) NewChannelQueue(name string) (*Queue, error) {
	channel, err := conn.NewChannel()
	if err != nil {
		return nil, err
	}

	return channel.NewQueue(name)
}

// Publish adds a task with specified params to the Queue
// priority: higher number, higher priority
func (q *Queue) Publish(params interface{}, priority uint8) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = q.Channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Priority:     priority,
		})
	if err != nil {
		return err
	}

	return nil
}

// Consume consumes messages from a queue
func (q *Queue) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := q.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
