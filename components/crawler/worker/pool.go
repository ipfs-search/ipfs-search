package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	samqp "github.com/streadway/amqp"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"

	"github.com/ipfs-search/ipfs-search/components/crawler"
	"github.com/ipfs-search/ipfs-search/components/extractor/tika"
	"github.com/ipfs-search/ipfs-search/components/index/elasticsearch"
	"github.com/ipfs-search/ipfs-search/components/protocol/ipfs"
	"github.com/ipfs-search/ipfs-search/components/queue/amqp"

	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs-search/ipfs-search/utils"
)

// Pool represents a pool of workers.
type Pool struct {
	config       *config.Config
	dialer       *utils.RetryingDialer
	consumeChans struct {
		Files       <-chan samqp.Delivery
		Directories <-chan samqp.Delivery
		Hashes      <-chan samqp.Delivery
	}
	crawler *crawler.Crawler

	*instr.Instrumentation
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

	// Many stat/ls connections
	// TODO: Make this configurable
	ipfsTransport := utils.GetHTTPTransport(w.dialer.DialContext, 1000)
	ipfsClient := &http.Client{Transport: ipfsTransport}
	protocol := ipfs.New(w.config.IPFSConfig(), ipfsClient, w.Instrumentation)

	// Limited Tika connections (as resources are generally known to be available by now)
	tikaTransport := utils.GetHTTPTransport(w.dialer.DialContext, 100)
	tikaClient := &http.Client{Transport: tikaTransport}
	extractor := tika.New(w.config.TikaConfig(), tikaClient, protocol, w.Instrumentation)

	w.crawler = crawler.New(w.config.CrawlerConfig(), indexes, queues, protocol, extractor, w.Instrumentation)

	return nil
}

func (w *Pool) getSearchClient() (*elasticsearch.Client, error) {
	clientConfig := &elasticsearch.ClientConfig{
		URL:       w.config.ElasticSearch.URL,
		Transport: utils.GetHTTPTransport(w.dialer.DialContext, 20),
		Debug:     false,

		// TODO: Make configurable.
		BulkIndexerWorkers:    2,
		BulkIndexerFlushBytes: 5 * 1024 * 1024, // 5 MB

		BulkGetterBatchSize:    24,
		BulkGetterBatchTimeout: 2 * time.Second,
	}

	return elasticsearch.NewClient(clientConfig, w.Instrumentation)
}

func startSearchWorker(ctx context.Context, esClient *elasticsearch.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := esClient.Work(ctx); err != nil {
				log.Printf("Error in ES client worker, restarting worker: %s", err)
				// Prevent overly tight restart loop
				time.Sleep(time.Second)
			}
		}
	}
}

func (w *Pool) getIndexes(ctx context.Context) (*crawler.Indexes, error) {
	esClient, err := w.getSearchClient()
	if err != nil {
		return nil, err
	}

	// Start 4 ES workers
	go startSearchWorker(ctx, esClient)
	go startSearchWorker(ctx, esClient)
	go startSearchWorker(ctx, esClient)
	go startSearchWorker(ctx, esClient)

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
		Partials: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: w.config.Indexes.Partials.Name},
		),
	}, nil
}

func (w *Pool) getQueues(ctx context.Context) (*crawler.Queues, error) {
	amqpConfig := &samqp.Config{
		Dial: w.dialer.Dial,
	}

	log.Println("Connecting to AMQP.")
	amqpConnection, err := amqp.NewConnection(ctx, w.config.AMQPConfig(), amqpConfig, w.Instrumentation)
	if err != nil {
		return nil, err
	}

	log.Println("Creating AMQP channels.")
	fq, err := amqpConnection.NewChannelQueue(ctx, w.config.Queues.Files.Name, w.config.Workers.FileWorkers)
	if err != nil {
		return nil, err
	}

	dq, err := amqpConnection.NewChannelQueue(ctx, w.config.Queues.Directories.Name, w.config.Workers.DirectoryWorkers)
	if err != nil {
		return nil, err
	}

	hq, err := amqpConnection.NewChannelQueue(ctx, w.config.Queues.Hashes.Name, w.config.Workers.HashWorkers)
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
	// TODO: Get SpanContext from Delivery.
	// ctx = trace.ContextWithRemoteSpanContext(ctx, p.SpanContext)
	ctx, span := w.Tracer.Start(ctx, "crawler.worker.crawlDelivery", trace.WithNewRoot())
	defer span.End()

	r := &t.AnnotatedResource{
		Resource: &t.Resource{},
	}

	if err := json.Unmarshal(d.Body, r); err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}

	if !r.IsValid() {
		err := fmt.Errorf("Invalid resource: %v", r)
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}

	log.Printf("Crawling '%s'", r)
	err := w.crawler.Crawl(ctx, r)
	log.Printf("Done crawling '%s', result: %v", r, err)

	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
	}

	return err
}

func (w *Pool) startWorker(ctx context.Context, deliveries <-chan samqp.Delivery, name string) {
	ctx, span := w.Tracer.Start(ctx, "crawler.worker.startWorker")
	defer span.End()

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
				// By default, do not retry.
				shouldRetry := false

				span.RecordError(ctx, err)

				if err := d.Reject(shouldRetry); err != nil {
					span.RecordError(ctx, err)
				}
			} else {
				if err := d.Ack(false); err != nil {
					span.RecordError(ctx, err)
				}
			}
		}
	}
}

func (w *Pool) startPool(ctx context.Context, deliveries <-chan samqp.Delivery, workers int, poolName string) {
	ctx, span := w.Tracer.Start(ctx, "crawler.worker.startPool")
	defer span.End()

	for i := 0; i < workers; i++ {
		name := fmt.Sprintf("%s-%d", poolName, i)
		go w.startWorker(ctx, deliveries, name)
	}
}

// Start launches the workerpool.
func (w *Pool) Start(ctx context.Context) {
	ctx, span := w.Tracer.Start(ctx, "crawler.worker.Start")
	defer span.End()

	log.Printf("Starting %d workers for files", w.config.Workers.FileWorkers)
	w.startPool(ctx, w.consumeChans.Files, w.config.Workers.FileWorkers, "files")

	log.Printf("Starting %d workers for hashes", w.config.Workers.HashWorkers)
	w.startPool(ctx, w.consumeChans.Hashes, w.config.Workers.HashWorkers, "hashes")

	log.Printf("Starting %d workers for directories", w.config.Workers.DirectoryWorkers)
	w.startPool(ctx, w.consumeChans.Directories, w.config.Workers.DirectoryWorkers, "directories")
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
		config:          c,
		Instrumentation: i,
	}

	err := w.init(ctx)

	return w, err
}
