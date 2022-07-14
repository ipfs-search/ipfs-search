package pool

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	samqp "github.com/rabbitmq/amqp091-go"

	"github.com/ipfs-search/ipfs-search/components/crawler"
	"github.com/ipfs-search/ipfs-search/components/worker"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/utils"
)

type consumeChans struct {
	Files       <-chan samqp.Delivery
	Directories <-chan samqp.Delivery
	Hashes      <-chan samqp.Delivery
}

// Pool represents a pool of pools.
type Pool struct {
	config  *config.Config
	dialer  *utils.RetryingDialer
	crawler *crawler.Crawler

	*consumeChans
	*instr.Instrumentation
}

func (p *Pool) startWorkers(ctx context.Context, deliveries <-chan samqp.Delivery, workers int, poolName string) {
	ctx, span := p.Tracer.Start(ctx, "crawler.pool.start")
	defer span.End()

	log.Printf("Starting %d workers for %s", workers, poolName)

	for i := 0; i < workers; i++ {
		name := fmt.Sprintf("%s-%d", poolName, i)
		cfg := &worker.Config{
			Name:         name,
			MaxLoadRatio: p.config.MaxLoadRatio,
			ThrottleMin:  p.config.ThrottleMin,
			ThrottleMax:  p.config.ThrottleMax,
		}

		worker := worker.New(cfg, p.crawler, p.Instrumentation)
		go worker.Start(ctx, deliveries)
	}
}

// Start launches the pool.
func (p *Pool) Start(ctx context.Context) {
	ctx, span := p.Tracer.Start(ctx, "crawler.pool.Start")
	defer span.End()

	p.startWorkers(ctx, p.consumeChans.Files, p.config.Workers.FileWorkers, "files")
	p.startWorkers(ctx, p.consumeChans.Hashes, p.config.Workers.HashWorkers, "hashes")
	p.startWorkers(ctx, p.consumeChans.Directories, p.config.Workers.DirectoryWorkers, "directories")
}

func (p *Pool) init(ctx context.Context) error {
	var err error

	p.dialer = &utils.RetryingDialer{
		Dialer: net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
		},
		Context: ctx,
	}

	log.Println("Initializing crawler.")
	if p.crawler, err = p.getCrawler(ctx); err != nil {
		return err
	}

	log.Println("Initializing consuming channels.")
	if p.consumeChans, err = p.getConsumeChans(ctx); err != nil {
		return err
	}

	return nil
}

// New initializes and returns a new pool.
func New(ctx context.Context, c *config.Config, i *instr.Instrumentation) (*Pool, error) {
	if i == nil {
		panic("Instrumentation cannot be null.")
	}

	if c == nil {
		panic("Config cannot be nil.")
	}

	p := &Pool{
		config:          c,
		Instrumentation: i,
	}

	err := p.init(ctx)

	return p, err
}
