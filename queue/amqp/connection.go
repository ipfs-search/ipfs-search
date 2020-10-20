package amqp

import (
	"github.com/streadway/amqp"
	"log"
)

// Connection wraps an AMQP connection
type Connection struct {
	conn *amqp.Connection
}

// NewConnection returns new AMQP connection
func NewConnection(url string) (*Connection, error) {
	amqpConn, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	blockings := amqpConn.NotifyBlocked(make(chan amqp.Blocking))
	go func() {
		for b := range blockings {
			if b.Active {
				log.Printf("TCP blocked: %q", b.Reason)
			} else {
				log.Printf("TCP unblocked")
			}
		}
	}()

	return &Connection{conn: amqpConn}, nil
}

// Channel creates an AMQP channel
func (c *Connection) Channel() (*Channel, error) {
	log.Printf("Creating AMQP channel")

	// Create channel
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	// Set Qos
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	return &Channel{
		ch: ch,
	}, nil
}

// NewChannelQueue returns a new queue on a new channel
func (c *Connection) NewChannelQueue(name string) (*Queue, error) {
	ch, err := c.Channel()
	if err != nil {
		return nil, err
	}

	return ch.Queue(name)
}

// Close closes the channel
func (c *Connection) Close() error {
	return c.conn.Close()
}
