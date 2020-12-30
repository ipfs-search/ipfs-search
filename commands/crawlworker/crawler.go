package crawlworker

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

	samqp "github.com/streadway/amqp"
)

type Worker struct {
	config       *config.Config
	httpClient   *http.Client
	instr        *instr.Instrumentation
	dialer       *RetryingDialer
	consumeChans struct {
		Files       <-chan samqp.Delivery
		Directories <-chan samqp.Delivery
		Hashes      <-chan samqp.Delivery
	}
	crawler *crawler.Crawler
}

func (w *Worker) makeCrawler(ctx context.Context) error {
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

func New(c *config.Config, i *instr.Instrumentation) *Worker {
	return &Worker{
		config: c,
		instr:  i,
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
	amqpConfig := &samqp.Config{
		Dial: f.dialer.Dial,
	}

	log.Println("Connecting to AMQP.")
	amqpConnection, err := amqp.NewConnection(ctx, f.config.AMQPConfig(), amqpConfig, f.instr)
	if err != nil {
		return nil, err
	}

	log.Println("Creating AMQP channels.")
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

func (w *Worker) crawlDelivery(ctx context.Context, d samqp.Delivery) error {
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

func (w *Worker) startWorker(ctx context.Context, deliveries <-chan samqp.Delivery) {
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

func (w *Worker) startWorkers(ctx context.Context, deliveries <-chan samqp.Delivery, workers uint) {
	var i uint
	for i = 0; i < workers; i++ {
		log.Println("Starting worker.")
		go w.startWorker(ctx, deliveries)
	}
}

func (w *Worker) Initialize(ctx context.Context) error {
	w.dialer = &RetryingDialer{
		Dialer: net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
		},
		Context: ctx,
	}
	w.httpClient = getHTTPClient(w.dialer.DialContext)

	log.Println("Initializing crawler.")
	if err := w.makeCrawler(ctx); err != nil {
		return err
	}

	log.Println("Initializing consuming channels.")
	return w.makeConsumeChans(ctx)
}

func (w *Worker) Start(ctx context.Context) {
	if w.crawler == nil {
		panic("Must call Initialize() before Start()")
	}

	log.Println("Starting workers.")
	// TODO: Clean this up, generating the queue when initializing a worker, perhaps even give workers an 'identity'
	// for better debugging.
	w.startWorkers(ctx, w.consumeChans.Files, w.config.Workers.FileWorkers)
	w.startWorkers(ctx, w.consumeChans.Hashes, w.config.Workers.HashWorkers)
	w.startWorkers(ctx, w.consumeChans.Directories, w.config.Workers.DirectoryWorkers)
}

func (w *Worker) makeConsumeChans(ctx context.Context) error {
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
