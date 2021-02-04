package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/extractor/tika"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/protocol/ipfs"
	"github.com/ipfs-search/ipfs-search/queue/amqp"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs-search/ipfs-search/utils"

	samqp "github.com/streadway/amqp"
)

// Pool represents a pool of workers.
type Pool struct {
	config       *config.Config
	httpClient   *http.Client
	instr        *instr.Instrumentation
	dialer       *utils.RetryingDialer
	consumeChans struct {
		Files       <-chan samqp.Delivery
		Directories <-chan samqp.Delivery
		Hashes      <-chan samqp.Delivery
	}
	crawler *crawler.Crawler
}

func (w *Pool) makeCrawler(ctx context.Context) error {
	var (
		queues  *crawler.Queues
		indexes *crawler.Indexes
		err     error
	)

	log.Println("Getting publish queues.")
	if queues, err = w.getQueues(ctx); err != nil {
		return err
	}

	log.Println("Getting indexes.")
	if indexes, err = w.getIndexes(ctx); err != nil {
		return err
	}

	protocol := ipfs.New(w.config.IPFSConfig(), w.httpClient, w.instr)
	extractor := tika.New(w.config.TikaConfig(), w.httpClient, protocol, w.instr)

	w.crawler = crawler.New(w.config.CrawlerConfig(), indexes, queues, protocol, extractor)

	return nil
}

func (w *Pool) getElasticClient() (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(w.config.ElasticSearch.URL),
		elastic.SetHttpClient(w.httpClient),
	)
}

func (w *Pool) getIndexes(ctx context.Context) (*crawler.Indexes, error) {
	esClient, err := w.getElasticClient()
	if err != nil {
		return nil, err
	}

	return &crawler.Indexes{
		Files: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: w.config.Indexes.Files.Name},
		),
		Directories: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: w.config.Indexes.Directories.Name},
		),
		Invalids: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: w.config.Indexes.Invalids.Name},
		),
	}, nil
}

func (w *Pool) getQueues(ctx context.Context) (*crawler.Queues, error) {
	amqpConfig := &samqp.Config{
		Dial: w.dialer.Dial,
	}

	log.Println("Connecting to AMQP.")
	amqpConnection, err := amqp.NewConnection(ctx, w.config.AMQPConfig(), amqpConfig, w.instr)
	if err != nil {
		return nil, err
	}

	log.Println("Creating AMQP channels.")
	fq, err := amqpConnection.NewChannelQueue(ctx, w.config.Queues.Files.Name)
	if err != nil {
		return nil, err
	}

	dq, err := amqpConnection.NewChannelQueue(ctx, w.config.Queues.Directories.Name)
	if err != nil {
		return nil, err
	}

	hq, err := amqpConnection.NewChannelQueue(ctx, w.config.Queues.Hashes.Name)
	if err != nil {
		return nil, err
	}

	return &crawler.Queues{
		Files:       fq,
		Directories: dq,
		Hashes:      hq,
	}, nil
}

func (w *Pool) crawlDelivery(ctx context.Context, d samqp.Delivery) error {
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

	return w.crawler.Crawl(ctx, r)
}

func (w *Pool) startWorker(ctx context.Context, deliveries <-chan samqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-deliveries:
			if !ok {
				// This is a fatal error; it should never happen - crash the program!
				panic("unexpected channel close")
			}
			if err := w.crawlDelivery(ctx, d); err != nil {
				// By default, retry.
				shouldRetry := true

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

func (w *Pool) startPool(ctx context.Context, deliveries <-chan samqp.Delivery, workers int) {
	for i := 0; i < workers; i++ {
		log.Println("Starting worker.")
		go w.startWorker(ctx, deliveries)
	}
}

// Start launches the workerpool.
func (w *Pool) Start(ctx context.Context) {
	log.Println("Starting workers.")
	// TODO: Clean this up, generating the queue when initializing a worker, perhaps even give workers an 'identity'
	// for better debugging.
	w.startPool(ctx, w.consumeChans.Files, w.config.Workers.FileWorkers)
	w.startPool(ctx, w.consumeChans.Hashes, w.config.Workers.HashWorkers)
	w.startPool(ctx, w.consumeChans.Directories, w.config.Workers.DirectoryWorkers)
}

func (w *Pool) makeConsumeChans(ctx context.Context) error {
	var (
		queues *crawler.Queues
		err    error
	)

	if queues, err = w.getQueues(ctx); err != nil {
		return err
	}

	if w.consumeChans.Files, err = queues.Files.Consume(ctx); err != nil {
		return err
	}

	if w.consumeChans.Directories, err = queues.Directories.Consume(ctx); err != nil {
		return err
	}

	if w.consumeChans.Hashes, err = queues.Hashes.Consume(ctx); err != nil {
		return err
	}

	return nil
}

func (w *Pool) init(ctx context.Context) error {
	w.dialer = &utils.RetryingDialer{
		Dialer: net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
		},
		Context: ctx,
	}
	w.httpClient = utils.GetHTTPClient(w.dialer.DialContext)

	log.Println("Initializing crawler.")
	if err := w.makeCrawler(ctx); err != nil {
		return err
	}

	log.Println("Initializing consuming channels.")
	return w.makeConsumeChans(ctx)
}

// NewPool initializes and returns a new worker pool.
func NewPool(ctx context.Context, c *config.Config, i *instr.Instrumentation) (*Pool, error) {
	w := &Pool{
		config: c,
		instr:  i,
	}

	err := w.init(ctx)

	return w, err
}
