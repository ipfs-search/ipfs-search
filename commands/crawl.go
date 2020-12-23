package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/extractor/tika"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/protocol/ipfs"
	"github.com/ipfs-search/ipfs-search/queue/amqp"
	t "github.com/ipfs-search/ipfs-search/types"
	samqp "github.com/streadway/amqp"

	"log"
	// "go.opentelemetry.io/otel/api/trace"
	// "go.opentelemetry.io/otel/codes"
)

func getIndexes(ctx context.Context, cfg *config.Config, instrumentation *instr.Instrumentation) (*crawler.Indexes, error) {
	esClient, err := getElasticClient(cfg.ElasticSearch.URL)
	if err != nil {
		return nil, err
	}

	return &crawler.Indexes{
		Files: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: cfg.Indexes.Files.Name},
		),
		Directories: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: cfg.Indexes.Directories.Name},
		),
		Invalids: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: cfg.Indexes.Invalids.Name},
		),
	}, nil
}

func getQueues(ctx context.Context, cfg *config.Config, instrumentation *instr.Instrumentation) (*crawler.Queues, error) {
	amqpConnection, err := amqp.NewConnection(ctx, cfg.AMQP.URL, instrumentation)
	if err != nil {
		return nil, err
	}

	fq, err := amqpConnection.NewChannelQueue(ctx, cfg.Queues.Files.Name)
	if err != nil {
		return nil, err
	}

	dq, err := amqpConnection.NewChannelQueue(ctx, cfg.Queues.Directories.Name)
	if err != nil {
		return nil, err
	}

	hq, err := amqpConnection.NewChannelQueue(ctx, cfg.Queues.Hashes.Name)
	if err != nil {
		return nil, err
	}

	return &crawler.Queues{
		Files:       fq,
		Directories: dq,
		Hashes:      hq,
	}, nil
}

func getCrawler(ctx context.Context, cfg *config.Config, instrumentation *instr.Instrumentation) (*crawler.Crawler, error) {
	httpClient := getHttpClient()

	queues, err := getQueues(ctx, cfg, instrumentation)
	if err != nil {
		return nil, err
	}

	indexes, err := getIndexes(ctx, cfg, instrumentation)
	if err != nil {
		return nil, err
	}

	protocol := ipfs.New(cfg.IPFSConfig(), httpClient, instrumentation)
	extractor := tika.New(cfg.TikaConfig(), httpClient, protocol, instrumentation)

	return crawler.New(cfg.CrawlerConfig(), indexes, queues, protocol, extractor), nil
}

type consumeChans struct {
	Files       <-chan samqp.Delivery
	Directories <-chan samqp.Delivery
	Hashes      <-chan samqp.Delivery
}

func getConsumeChans(ctx context.Context, cfg *config.Config, instrumentation *instr.Instrumentation) (*consumeChans, error) {
	var c consumeChans

	queues, err := getQueues(ctx, cfg, instrumentation)
	if err != nil {
		return nil, err
	}

	c.Files, err = queues.Files.Consume(ctx)
	if err != nil {
		return nil, err
	}

	c.Directories, err = queues.Directories.Consume(ctx)
	if err != nil {
		return nil, err
	}

	c.Hashes, err = queues.Hashes.Consume(ctx)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func crawlDelivery(ctx context.Context, d samqp.Delivery, c *crawler.Crawler) error {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{},
	}

	if err := json.Unmarshal(d.Body, r); err != nil {
		return err
	}

	if !r.IsValid() {
		return fmt.Errorf("Invalid resource: %v", r)
	}

	fmt.Printf("Crawling: %v\n", r)

	return c.Crawl(ctx, r)
}

func work(ctx context.Context, consumeChan <-chan samqp.Delivery, c *crawler.Crawler) {
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-consumeChan:
			if !ok {
				// This is a fatal error; it should never happen - crash the program!
				panic("consume channel closed")
			}
			if err := crawlDelivery(ctx, d, c); err != nil {
				shouldRetry := crawler.IsTemporaryErr(err)

				if err := d.Reject(shouldRetry); err != nil {
					// span.RecordError(ctx, err)
				}
				log.Printf("Error '%s' in delivery '%s'", err, d.Body)
				// span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
			} else {
				if err := d.Ack(false); err != nil {
					// span.RecordError(ctx, err)
				}
			}
		}
	}
}

func makeWorkers(ctx context.Context, consumeChan <-chan samqp.Delivery, c *crawler.Crawler, n uint) {
	var i uint
	for i = 0; i < n; i++ {
		go work(ctx, consumeChan, c)
	}
}

// Crawl configures and initializes crawling
func Crawl(ctx context.Context, cfg *config.Config) error {
	instFlusher, err := instr.Install("ipfs-crawler")
	if err != nil {
		log.Fatal(err)
	}
	defer instFlusher()

	instrumentation := instr.New()
	tracer := instrumentation.Tracer

	ctx, span := tracer.Start(ctx, "commands.Crawl")
	defer span.End()

	c, err := getCrawler(ctx, cfg, instrumentation)
	if err != nil {
		return err
	}

	consumeChans, err := getConsumeChans(ctx, cfg, instrumentation)

	makeWorkers(ctx, consumeChans.Files, c, cfg.Workers.FileWorkers)
	makeWorkers(ctx, consumeChans.Hashes, c, cfg.Workers.HashWorkers)
	makeWorkers(ctx, consumeChans.Directories, c, cfg.Workers.DirectoryWorkers)

	// Context closure or panic is the only way to stop crawling
	<-ctx.Done()

	return ctx.Err()
}
