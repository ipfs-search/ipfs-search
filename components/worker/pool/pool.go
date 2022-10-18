package pool

import (
	"context"
	"log"
	"net"
	"time"

	samqp "github.com/rabbitmq/amqp091-go"

	"github.com/ipfs-search/ipfs-search/components/crawler"
	"github.com/ipfs-search/ipfs-search/components/worker"
	"github.com/ipfs-search/ipfs-search/components/worker/group"
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

func (p *Pool) getWorkerConfig(name string) *worker.Config {
	return &worker.Config{
		Name:         name,
		MaxLoadRatio: p.config.MaxLoadRatio,
		ThrottleMin:  p.config.ThrottleMin,
		ThrottleMax:  p.config.ThrottleMax,
	}
}

func (p *Pool) getWorker(name string) *worker.Worker {
	cfg := p.getWorkerConfig(name)
	return worker.New(cfg, p.crawler, p.Instrumentation)
}

// Start launches the pool.
func (p *Pool) Start(ctx context.Context) error {
	ctx, span := p.Tracer.Start(ctx, "crawler.pool.Start")
	defer span.End()

	g := group.New(ctx, p.getWorker, p.Instrumentation)

	g.Go(p.consumeChans.Files, p.config.Workers.FileWorkers, "files")
	g.Go(p.consumeChans.Hashes, p.config.Workers.HashWorkers, "hashes")
	g.Go(p.consumeChans.Directories, p.config.Workers.DirectoryWorkers, "directories")

	return g.Wait()
}

// Init initializes a worker poool.
func (p *Pool) Init(ctx context.Context) error {
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

// New returns a new pool.
func New(c *config.Config, i *instr.Instrumentation) *Pool {
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

	return p
}
