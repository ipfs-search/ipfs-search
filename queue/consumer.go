package queue

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// Func is the function consuming a queue
type Func func(interface{}) error

// Consumer takes messages from a TaskQueue and consumes them
type Consumer struct {
	Func    Func
	ErrChan chan<- error
	Queue   *TaskQueue
	Params  interface{}
}

func (c *Consumer) messagePanic(message *amqp.Delivery) {
	if r := recover(); r != nil {
		log.Printf("Panic in: %s", message.Body)

		// Permanently remove message from original queue
		message.Reject(false)

		err, ok := r.(error)

		if !ok {
			err = fmt.Errorf("Unassertable panic error: %v", r)
		}

		c.ErrChan <- err
	}

}

// processMessage processes a single ampq message
func (c *Consumer) processMessage(message *amqp.Delivery) {
	defer c.messagePanic(message)

	log.Printf("Received a message: %s", message.Body)

	err := json.Unmarshal(message.Body, &c.Params)
	if err != nil {
		panic(&err)
	}

	err = c.Func(c.Params)

	if err == nil {
		// Everything went fine, ack the message
		message.Ack(false)
	} else {
		// Send error through channel
		c.ErrChan <- err

		// Don't retry
		message.Reject(false)
	}
}

// Start consuming messages
func (c *Consumer) Start() error {
	msgs, err := c.Queue.c.ch.Consume(
		c.Queue.q.Name, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return err
	}

	// Task loop go routine
	go func() {
		for message := range msgs {
			c.processMessage(&message)
		}
	}()

	log.Printf("Started worker for queue '%s'", c.Queue.q.Name)

	return nil

}
