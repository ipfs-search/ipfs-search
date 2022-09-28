package pool

import (
	"context"
	"log"

	samqp "github.com/rabbitmq/amqp091-go"

	"github.com/ipfs-search/ipfs-search/components/crawler"
	"github.com/ipfs-search/ipfs-search/components/queue/amqp"
)

func (p *Pool) getQueues(ctx context.Context) (*crawler.Queues, error) {
	amqpConfig := &samqp.Config{
		Dial: p.dialer.Dial,
	}

	log.Println("Connecting to AMQP.")
	amqpConnection, err := amqp.NewConnection(ctx, p.config.AMQPConfig(), amqpConfig, p.Instrumentation)
	if err != nil {
		return nil, err
	}

	log.Println("Creating AMQP channels.")
	fq, err := amqpConnection.NewChannelQueue(ctx, p.config.Queues.Files.Name, p.config.Workers.FileWorkers)
	if err != nil {
		return nil, err
	}

	dq, err := amqpConnection.NewChannelQueue(ctx, p.config.Queues.Directories.Name, p.config.Workers.DirectoryWorkers)
	if err != nil {
		return nil, err
	}

	hq, err := amqpConnection.NewChannelQueue(ctx, p.config.Queues.Hashes.Name, p.config.Workers.HashWorkers)
	if err != nil {
		return nil, err
	}

	return &crawler.Queues{
		Files:       fq,
		Directories: dq,
		Hashes:      hq,
	}, nil
}
