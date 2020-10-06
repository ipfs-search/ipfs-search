package amqp

import (
	"context"
	"log"

	"github.com/ipfs-search/ipfs-search/queue"
)

// PublisherFactory automates creation of AMQP Publishers.
type PublisherFactory struct {
	AMQPURL string
	Queue   string
}

func (f PublisherFactory) NewPublisher(ctx context.Context) (queue.Publisher, error) {
	// Create and configure add queue
	conn, err := NewConnection(f.AMQPURL)
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
