package queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type TaskChannel struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewChannel() (*TaskChannel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
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
		conn: conn,
		ch:   ch,
	}, nil
}

func (c *TaskChannel) Close() error {
	if c.conn != nil {
		// Connection exists, defer close
		defer c.conn.Close()
	}

	if c.ch != nil {
		err := c.ch.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

type TaskQueue struct {
	c *TaskChannel
	q *amqp.Queue
}

func NewTaskQueue(c *TaskChannel, queue_name string) (*TaskQueue, error) {
	q, err := c.ch.QueueDeclare(
		queue_name, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
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

func (t TaskQueue) StartConsumer(worker func(map[string]interface{}) error, errc chan error) error {
	var params map[string]interface{}

	msgs, err := t.c.ch.Consume(
		t.q.Name, // queue
		"",       // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			err := json.Unmarshal(d.Body, &params)
			if err != nil {
				// Message is fucked up, don't retry
				d.Reject(false)
				errc <- err
			}

			err = worker(params)
			if err != nil {
				// Reject but do retry
				d.Reject(true)
				errc <- err
			}

			d.Ack(false)
		}
	}()

	log.Printf("Started worker for queue '%s'", t.q.Name)

	return nil
}

func (t TaskQueue) AddTask(params map[string]interface{}) error {
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
