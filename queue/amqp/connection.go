package amqp

import (
	"github.com/streadway/amqp"
)

// Connection wraps an AMQP connection
type Connection struct {
	conn *amqp.Connection
}

// NewConnection returns new AMQP connection
func NewConnection(url string) (*Connection, error) {
	connection, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	return &Connection{conn: connection}, nil
}

// Channel creates an AMQP channel
func (c *Connection) Channel() (*Channel, error) {
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