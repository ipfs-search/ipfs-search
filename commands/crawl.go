package commands

import (
	"encoding/json"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/queue"
)

type WorkerFactory struct {
	PubConnection *queue.Connection
	ConConnection *queue.Connection
	ErrChan       chan<- error
}

func (f *WorkerFactory) NewCrawler() (*crawler.Crawler, error) {
	// Setup channels
	filePubChannel, err := f.PubConnection.NewChannel()
	fileQueue, err := filePubChannel.NewQueue("files")

	hashPubChannel, err := f.PubConnection.NewChannel()
	hashQueue, err := filePubChannel.NewQueue("hashes")

	// This is where we need config

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

	return &crawler.Crawler{
		Config:    config.CrawlerConfig,
		Shell:     sh,
		Indexer:   id,
		FileQueue: fileQueue,
		HashQueue: hashQueue,
	}, nil
}

func (f *WorkerFactory) NewHashWorker() (*queue.Worker, error) {
	conChannel, err := f.ConConnection.NewChannel()
	hashConQueue, err := filePubChannel.NewQueue("hashes")

	crawler, err := f.NewCrawler()

	var hashFunc = func(msg queue.Message) error {
		// Unmarshall into
		args := make(crawler.Args)
		err := json.Unmarshal(msg.Delivery.Body, args)
		if err != nil {
			return err
		}

		return c.CrawlHash(args)
	}

	return queue.Worker{
		ErrChan: errc,
		Func:    fileFunc,
		Queue:   hashConQueue,
	}, nil
}

// Crawl configures and initializes crawling
func Crawl() error {
	// This is common for workers
	errc := make(chan error, 1)

	pubConnection, err := queue.NewConnection("<AMQP URL>")
	if err != nil {
		return err
	}

	conConnection, err := queue.NewConnection("<AMQP URL>")
	if err != nil {
		return err
	}

	factory, err = WorkerFactory{
		PubConnection: pubConnection,
		ConConnection: conConnection,
		ErrChan:       errc,
	}

	hashGroup, err = worker.Group{
		Count: 100,
		factory.NewHashWorker,
	}

	// Start work loop
	hashGroup.Work()
}
