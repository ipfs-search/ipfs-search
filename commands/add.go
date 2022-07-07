package commands

import (
	"context"
	"net"
	"time"

	samqp "github.com/rabbitmq/amqp091-go"

	"github.com/ipfs-search/ipfs-search/components/queue/amqp"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs-search/ipfs-search/utils"
)

// AddHash queues a single IPFS hash for indexing
func AddHash(ctx context.Context, cfg *config.Config, hash string) error {
	instFlusher, err := instr.Install(cfg.InstrConfig(), "ipfs-crawler add")
	if err != nil {
		return err
	}
	defer instFlusher()

	i := instr.New()

	dialer := &utils.RetryingDialer{
		Dialer: net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
		},
		Context: ctx,
	}

	amqpConfig := &samqp.Config{
		Dial: dialer.Dial,
	}

	f := amqp.PublisherFactory{
		Config:          cfg.AMQPConfig(),
		Queue:           "hashes",
		AMQPConfig:      amqpConfig,
		Instrumentation: i,
	}

	queue, err := f.NewPublisher(ctx)
	if err != nil {
		return err
	}

	resource := &t.Resource{
		Protocol: t.IPFSProtocol,
		ID:       hash,
	}

	r := t.AnnotatedResource{
		Resource: resource,
		Source:   t.ManualSource,
	}

	// TODO: Use provider here

	// Add with highest priority, as this is supposed to be available
	return queue.Publish(ctx, &r, 9)
}
