package queue

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var conn *amqp.Connection

type TaskChannel struct {
	ch *amqp.Channel
}

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

func (t TaskQueue) StartConsumer(worker func(interface{}) error, params interface{}, errc chan error) error {
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

	// Task loop go routine
	go func() {
		for d := range msgs {

			// Anonymous function to catch panics without disrupting msgs loop
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Panic in: %s", d.Body)

						// Permanently remove message from original queue
						d.Reject(false)

						var ok bool
						err, ok = r.(error)

						if !ok {
							err = fmt.Errorf("Unassertable panic error: %v", r)
						}

						errc <- err
					}
				}()

				log.Printf("Received a message: %s", d.Body)

				err = json.Unmarshal(d.Body, &params)
				if err != nil {
					panic(&err)
				}

				err = worker(params)

				if err == nil {
					// Everything went fine, ack the message
					d.Ack(false)
				} else {
					// Send error through channel
					errc <- err

					// Don't retry
					d.Reject(false)
				}
			}()
		}
	}()

	log.Printf("Started worker for queue '%s'", t.q.Name)

	return nil
}

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
