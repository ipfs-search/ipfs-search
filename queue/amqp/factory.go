package amqp

import (
	"context"
	"log"

	"github.com/ipfs-search/ipfs-search/queue"
)

// AMQPublisherFactory automates creation of AMQP Publishers.
type AMQPPublisherFactory struct {
	AMQPURL string
	Queue   string
}

func (f AMQPPublisherFactory) NewPublisher(ctx context.Context) (queue.Publisher, error) {
	// Create and configure add queue
	conn, err := queue.NewConnection(f.AMQPURL)
	if err != nil {
		return nil, err
	}

	// Close connection when context closes
	go func() {
		<-ctx.Done()
		log.Printf("Closing AMQP connection; context closed")
		conn.Close()
	}()

	return conn.NewChannelQueue(f.Queue)
}
