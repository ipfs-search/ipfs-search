package queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection

// TaskChannel wraps an AMQP channel for tasks
type TaskChannel struct {
	ch *amqp.Channel
}

// NewChannel initialises an AMQP channel
func NewChannel() (*TaskChannel, error) {
	var err error

	if conn == nil {
		conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			return nil, err
		}
	}

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

	return &TaskChannel{
		ch: ch,
	}, nil
}

// Close closes a TaskChannel
func (c *TaskChannel) Close() error {
	if conn != nil {
		// Connection exists, defer close
		defer conn.Close()
	}

	if c.ch != nil {
		err := c.ch.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// TaskQueue wraps an channel/queue for tasks
type TaskQueue struct {
	c *TaskChannel
	q *amqp.Queue
}

// NewTaskQueue creates a named queue on a given chennel
func NewTaskQueue(c *TaskChannel, queueName string) (*TaskQueue, error) {
	q, err := c.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	tq := TaskQueue{
		c: c,
		q: &q,
	}

	return &tq, nil
}

// AddTask adds a task with specified params to the TaskQueue
func (t TaskQueue) AddTask(params interface{}) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = t.c.ch.Publish(
		"",       // exchange
		t.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		return err
	}

	return nil
}
