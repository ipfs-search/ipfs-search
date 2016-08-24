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

func receive_message(worker func(interface{}) error, d *amqp.Delivery, params interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Reject message, don't retry
			d.Reject(false)

			log.Println("Panic in: %s", d.Body)

			var ok bool
			err, ok = r.(error)

			if !ok {
				err = fmt.Errorf("%T: %v", r)
			}

			return
		}
	}()

	log.Printf("Received a message: %s", d.Body)

	err = json.Unmarshal(d.Body, &params)
	if err != nil {
		panic(&err)
	}

	err = worker(params)
	if err != nil {
		// Reject, retry
		d.Reject(true)
		return
	}

	d.Ack(false)

	return
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

	go func() {
		for d := range msgs {
			err := receive_message(worker, &d, params)
			if err != nil {
				errc <- err
			}
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
