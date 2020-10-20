package factory

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/index"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs-search/ipfs-search/queue/amqp"
	"github.com/ipfs-search/ipfs-search/worker"
	"github.com/ipfs/go-ipfs-api"
	"github.com/olivere/elastic/v7"
	samqp "github.com/streadway/amqp"
	"log"
)

// Factory creates hash and file crawl workers
type Factory struct {
	crawlerConfig *crawler.Config
	pubConnection *amqp.Connection
	conConnection *amqp.Connection
	errChan       chan<- error

	fileIndex      index.Index
	directoryIndex index.Index
	invalidIndex   index.Index

	shell *shell.Shell
}

// New creates a new crawl worker factory
func New(ctx context.Context, config *Config, errc chan<- error) (*Factory, error) {
	pubConnection, err := amqp.NewConnection(config.AMQPURL)
	if err != nil {
		return nil, err
	}

	conConnection, err := amqp.NewConnection(config.AMQPURL)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to AMQP")

	// Create and configure Ipfs shell
	sh := shell.NewShell(config.IpfsAPI)
	sh.SetTimeout(config.IpfsTimeout)

	es, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(config.ElasticSearchURL))
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to ElasticSearch.")

	indexes := elasticsearch.NewMulti(es, config.Indexes["files"], config.Indexes["directories"], config.Indexes["invalids"])

	return &Factory{
		crawlerConfig:  config.CrawlerConfig,
		pubConnection:  pubConnection,
		conConnection:  conConnection,
		errChan:        errc,
		shell:          sh,
		fileIndex:      indexes[0],
		directoryIndex: indexes[1],
		invalidIndex:   indexes[2],
	}, nil
}

func (f *Factory) newCrawler() (*crawler.Crawler, error) {
	log.Printf("Initializing crawler")

	fileQueue, err := f.pubConnection.NewChannelQueue("files")
	if err != nil {
		return nil, err
	}

	hashQueue, err := f.pubConnection.NewChannelQueue("hashes")
	if err != nil {
		return nil, err
	}

	return &crawler.Crawler{
		Config: f.crawlerConfig,
		Shell:  f.shell,

		FileIndex:      f.fileIndex,
		DirectoryIndex: f.directoryIndex,
		InvalidIndex:   f.invalidIndex,

		FileQueue: fileQueue,
		HashQueue: hashQueue,
	}, nil
}

// newWorker generalizes creating new workers; it takes a queue name and a
// crawlFunc, which takes an Indexable and returns the function performing
// the actual crawling
func (f *Factory) newWorker(queueName string, crawl CrawlFunc) (worker.Worker, error) {
	conQueue, err := f.conConnection.NewChannelQueue(queueName)
	if err != nil {
		return nil, err
	}

	c, err := f.newCrawler()
	if err != nil {
		return nil, err
	}

	log.Printf("Crawler initialised")

	// A MessageWorkerFactory generates a worker for every message in a queue
	messageWorkerFactory := func(msg *samqp.Delivery) worker.Worker {
		log.Printf("Creating worker for message %s", msg.Body)

		return &Worker{
			Crawler:   c,
			Delivery:  msg,
			CrawlFunc: crawl,
		}
	}

	log.Printf("Creating worker for queue %s", queueName)
	return queue.NewWorker(f.errChan, conQueue, messageWorkerFactory), nil
}

// NewHashWorker returns a new hash crawl worker
func (f *Factory) NewHashWorker() (worker.Worker, error) {
	return f.newWorker("hashes", func(i *crawler.Indexable) func(context.Context) error {
		return i.CrawlHash
	})
}

// NewFileWorker returns a new file crawl worker
func (f *Factory) NewFileWorker() (worker.Worker, error) {
	return f.newWorker("files", func(i *crawler.Indexable) func(context.Context) error {
		return i.CrawlFile
	})
}
