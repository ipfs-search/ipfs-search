package factory

import (
	"context"
	"fmt"
	"time"

	"net"

	"github.com/ipfs-search/ipfs-search/components/queue/amqp"
	"github.com/ipfs-search/ipfs-search/components/sniffer"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/utils"

	"github.com/ipfs/go-datastore"
	samqp "github.com/rabbitmq/amqp091-go"
)

func getConfig() (*config.Config, error) {
	cfg, err := config.Get("")
	if err != nil {
		return nil, err
	}

	if err = cfg.Check(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getInstr(cfg *instr.Config) (*instr.Instrumentation, func(), error) {
	instFlusher, err := instr.Install(cfg, "ipfs-sniffer")
	if err != nil {
		return nil, nil, err
	}
	return instr.New(), instFlusher, nil
}

func getQueue(ctx context.Context, cfg *amqp.Config, i *instr.Instrumentation) amqp.PublisherFactory {
	// Retrying dialer for connecting
	dialer := &utils.RetryingDialer{
		Dialer: net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
		},
		Context: ctx,
	}
	samqpConfig := &samqp.Config{
		Dial: dialer.Dial,
	}

	return amqp.PublisherFactory{
		Config:          cfg,
		AMQPConfig:      samqpConfig,
		Queue:           "hashes",
		Instrumentation: i,
	}
}

func getSniffer(cfg *sniffer.Config, ds datastore.Batching, q amqp.PublisherFactory, i *instr.Instrumentation) (*sniffer.Sniffer, error) {
	return sniffer.New(cfg, ds, q, i)
}

// Start initialises a sniffer and all its dependencies and launches it in a goroutine, returning a wrapped context
// and datastore, which should replace the original ones, or an error from initialisation.
func Start(ctx context.Context, ds datastore.Batching) (context.Context, datastore.Batching, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, nil, err
	}

	i, instFlusher, err := getInstr(cfg.InstrConfig())
	if err != nil {
		return nil, nil, err
	}

	// Create context which can be canceled by sniffer so as to propagate failure from sniffer goroutine.
	ctx, cancel := context.WithCancel(ctx)

	q := getQueue(ctx, cfg.AMQPConfig(), i)

	s, err := getSniffer(cfg.SnifferConfig(), ds, q, i)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	// Use batched datastore
	ds = s.Batching()

	// Start sniffer
	go func() {
		// Cancel parent context when done
		defer cancel()
		defer instFlusher()

		err = s.Sniff(ctx)
		fmt.Printf("Sniffer exited: %s\n", err)
	}()

	return ctx, ds, nil
}
