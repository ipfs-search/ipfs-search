package commands

import (
	"context"
	"net"
	"time"

	samqp "github.com/streadway/amqp"

	"github.com/ipfs-search/ipfs-search/commands/crawlworker"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/queue/amqp"
	t "github.com/ipfs-search/ipfs-search/types"
)

// AddHash queues a single IPFS hash for indexing
func AddHash(ctx context.Context, cfg *config.Config, hash string) error {
	instFlusher, err := instr.Install("ipfs-search-add")
	if err != nil {
		return err
	}
	defer instFlusher()

	i := instr.New()

	dialer := &crawlworker.RetryingDialer{
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

	provider := t.Provider{
		Resource: resource,
		Date:     time.Now(),
	}

	// Add with highest priority, as this is supposed to be available
	return queue.Publish(ctx, provider, 9)
}
