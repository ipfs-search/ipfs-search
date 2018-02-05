package factory

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/indexer"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs-search/ipfs-search/worker"
	"github.com/ipfs/go-ipfs-api"
	"github.com/streadway/amqp"
)

// Factory creates hash and file crawl workers
type Factory struct {
	crawlerConfig *crawler.Config
	pubConnection *queue.Connection
	conConnection *queue.Connection
	errChan       chan<- error
	indexer       *indexer.Indexer
	shell         *shell.Shell
}

// New creates a new crawl worker factory
func New(config *Config, errc chan<- error) (*Factory, error) {
	pubConnection, err := queue.NewConnection(config.AMQPURL)
	if err != nil {
		return nil, err
	}

	conConnection, err := queue.NewConnection(config.AMQPURL)
	if err != nil {
		return nil, err
	}

	// Create and configure Ipfs shell
	sh := shell.NewShell(config.IpfsAPI)
	sh.SetTimeout(config.IpfsTimeout)

	el, err := getElastic(config.ElasticSearchURL)
	if err != nil {
		return nil, err
	}

	// Create elasticsearch indexer
	id := &indexer.Indexer{
		ElasticSearch: el,
	}

	return &Factory{
		crawlerConfig: config.CrawlerConfig,
		pubConnection: pubConnection,
		conConnection: conConnection,
		errChan:       errc,
		shell:         sh,
		indexer:       id,
	}, nil
}

func (f *Factory) newCrawler() (*crawler.Crawler, error) {
	fileQueue, err := f.pubConnection.NewChannelQueue("files")
	if err != nil {
		return nil, err
	}

	hashQueue, err := f.pubConnection.NewChannelQueue("hashes")
	if err != nil {
		return nil, err
	}

	return &crawler.Crawler{
		Config:    f.crawlerConfig,
		Shell:     f.shell,
		Indexer:   f.indexer,
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

	// A MessageWorkerFactory generates a worker for every message in a queue
	messageWorkerFactory := func(msg *amqp.Delivery) worker.Worker {
		return &Worker{
			Crawler:   c,
			Delivery:  msg,
			CrawlFunc: crawl,
		}
	}

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
