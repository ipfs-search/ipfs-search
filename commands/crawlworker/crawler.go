package crawlworker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/olivere/elastic/v7"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/extractor/tika"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/protocol/ipfs"
	"github.com/ipfs-search/ipfs-search/queue/amqp"
	t "github.com/ipfs-search/ipfs-search/types"

	samqp "github.com/streadway/amqp"
)

type Worker struct {
	config     *config.Config
	httpClient *http.Client
	instr      *instr.Instrumentation
}

func (f *Worker) getCrawler(ctx context.Context) (*crawler.Crawler, error) {
	queues, err := f.getQueues(ctx)
	if err != nil {
		return nil, err
	}

	indexes, err := f.getIndexes(ctx)
	if err != nil {
		return nil, err
	}

	protocol := ipfs.New(f.config.IPFSConfig(), f.httpClient, f.instr)
	extractor := tika.New(f.config.TikaConfig(), f.httpClient, protocol, f.instr)

	return crawler.New(f.config.CrawlerConfig(), indexes, queues, protocol, extractor), nil
}

func getHttpClient() *http.Client {
	// TODO: Get more advanced client with circuit breaking etc. over manual
	// retrying get etc.
	// Ref: https://github.com/gojek/heimdall#creating-a-hystrix-like-circuit-breaker
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func New(c *config.Config, i *instr.Instrumentation) *Worker {
	return &Worker{
		config:     c,
		httpClient: getHttpClient(),
		instr:      i,
	}
}

func (f *Worker) getElasticClient() (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(f.config.ElasticSearch.URL),
		elastic.SetHttpClient(f.httpClient),
	)
}

func (f *Worker) getIndexes(ctx context.Context) (*crawler.Indexes, error) {
	esClient, err := f.getElasticClient()
	if err != nil {
		return nil, err
	}

	return &crawler.Indexes{
		Files: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: f.config.Indexes.Files.Name},
		),
		Directories: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: f.config.Indexes.Directories.Name},
		),
		Invalids: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: f.config.Indexes.Invalids.Name},
		),
	}, nil
}

func (f *Worker) getQueues(ctx context.Context) (*crawler.Queues, error) {
	amqpConnection, err := amqp.NewConnection(ctx, f.config.AMQPConfig(), f.instr)
	if err != nil {
		return nil, err
	}

	fq, err := amqpConnection.NewChannelQueue(ctx, f.config.Queues.Files.Name)
	if err != nil {
		return nil, err
	}

	dq, err := amqpConnection.NewChannelQueue(ctx, f.config.Queues.Directories.Name)
	if err != nil {
		return nil, err
	}

	hq, err := amqpConnection.NewChannelQueue(ctx, f.config.Queues.Hashes.Name)
	if err != nil {
		return nil, err
	}

	return &crawler.Queues{
		Files:       fq,
		Directories: dq,
		Hashes:      hq,
	}, nil
}

func (w *Worker) crawlDelivery(ctx context.Context, c *crawler.Crawler, d samqp.Delivery) error {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{},
	}

	if err := json.Unmarshal(d.Body, r); err != nil {
		return err
	}

	if !r.IsValid() {
		return fmt.Errorf("Invalid resource: %v", r)
	}

	log.Printf("Crawling: %v\n", r)

	return c.Crawl(ctx, r)
}

func (w *Worker) startWorker(ctx context.Context, c *crawler.Crawler, deliveries <-chan samqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-deliveries:
			if !ok {
				// This is a fatal error; it should never happen - crash the program!
				panic("unexpected channel close")
			}
			if err := w.crawlDelivery(ctx, c, d); err != nil {
				shouldRetry := crawler.IsTemporaryErr(err)

				if err := d.Reject(shouldRetry); err != nil {
					log.Printf("Reject error %s\n", d.Body)
					// span.RecordError(ctx, err)
				}
				log.Printf("Error '%s' in delivery '%s'", err, d.Body)
				// span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
			} else {
				if err := d.Ack(false); err != nil {
					log.Printf("Ack error %s\n", d.Body)

					// span.RecordError(ctx, err)
				}
				log.Printf("Done crawling: %s\n", d.Body)
			}
		}
	}
}

func (w *Worker) startWorkers(ctx context.Context, c *crawler.Crawler, deliveries <-chan samqp.Delivery, workers uint) {
	var i uint
	for i = 0; i < workers; i++ {
		go w.startWorker(ctx, c, deliveries)
	}
}

// TODO: This would prevent us from passing crawlers around and it would make Start not return errors.
// func (w *Worker) Initialize() error {
// 	w.crawler, err := w.getCrawler(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	w.consumeChans, err := w.getConsumeChans(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (w *Worker) Start(ctx context.Context) error {
	c, err := w.getCrawler(ctx)
	if err != nil {
		return err
	}

	consumeChans, err := w.getConsumeChans(ctx)
	if err != nil {
		return err
	}

	w.startWorkers(ctx, c, consumeChans.Files, w.config.Workers.FileWorkers)
	w.startWorkers(ctx, c, consumeChans.Hashes, w.config.Workers.HashWorkers)
	w.startWorkers(ctx, c, consumeChans.Directories, w.config.Workers.DirectoryWorkers)

	return nil
}

type consumeChans struct {
	Files       <-chan samqp.Delivery
	Directories <-chan samqp.Delivery
	Hashes      <-chan samqp.Delivery
}

func (w *Worker) getConsumeChans(ctx context.Context) (*consumeChans, error) {
	var c consumeChans

	queues, err := w.getQueues(ctx)
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
