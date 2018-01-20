package factory

import (
	"context"
	"encoding/json"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/indexer"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs-search/ipfs-search/worker"
	"github.com/ipfs/go-ipfs-api"
)

type Factory struct {
	config        *Config
	pubConnection *queue.Connection
	conConnection *queue.Connection
	errChan       chan<- error
}

func New(config *Config, errc chan<- error) (*Factory, error) {
	pubConnection, err := queue.NewConnection(config.AMQPURL)
	if err != nil {
		return nil, err
	}

	conConnection, err := queue.NewConnection(config.AMQPURL)
	if err != nil {
		return nil, err
	}

	return &Factory{
		config:        config,
		pubConnection: pubConnection,
		conConnection: conConnection,
		errChan:       errc,
	}, nil
}

func (f *Factory) NewCrawler() (*crawler.Crawler, error) {
	// Setup channels
	filePubChannel, err := f.pubConnection.NewChannel()
	fileQueue, err := filePubChannel.NewQueue("files")

	hashPubChannel, err := f.pubConnection.NewChannel()
	hashQueue, err := hashPubChannel.NewQueue("hashes")

	// This is where we need config

	// Create and configure Ipfs shell
	sh := shell.NewShell(f.config.IpfsAPI)
	sh.SetTimeout(f.config.IpfsTimeout)

	el, err := getElastic(f.config.ElasticSearchURL)
	if err != nil {
		return nil, err
	}

	// Create elasticsearch indexer
	id := &indexer.Indexer{
		ElasticSearch: el,
	}

	return &crawler.Crawler{
		Config:    f.config.CrawlerConfig,
		Shell:     sh,
		Indexer:   id,
		FileQueue: fileQueue,
		HashQueue: hashQueue,
	}, nil
}

func (f *Factory) NewHashWorker() (worker.Worker, error) {
	conChannel, err := f.conConnection.NewChannel()
	hashConQueue, err := conChannel.NewQueue("hashes")
	if err != nil {
		return nil, err
	}

	c, err := f.NewCrawler()
	if err != nil {
		return nil, err
	}

	var hashFunc = func(ctx context.Context, msg *queue.WorkerMessage) error {
		// Unmarshall into
		args := &crawler.Args{}
		err := json.Unmarshal(msg.Delivery.Body, args)
		if err != nil {
			return err
		}

		i := crawler.Indexable{
			Args:    args,
			Crawler: c,
		}

		return i.CrawlHash()
	}

	return &queue.Worker{
		ErrChan: f.errChan,
		Func:    hashFunc,
		Queue:   hashConQueue,
	}, nil
}

func (f *Factory) NewFileWorker() (worker.Worker, error) {
	conChannel, err := f.conConnection.NewChannel()
	fileConQueue, err := conChannel.NewQueue("filees")
	if err != nil {
		return nil, err
	}

	c, err := f.NewCrawler()
	if err != nil {
		return nil, err
	}

	var fileFunc = func(ctx context.Context, msg *queue.WorkerMessage) error {
		// Unmarshall into
		args := &crawler.Args{}
		err := json.Unmarshal(msg.Delivery.Body, args)
		if err != nil {
			return err
		}

		i := crawler.Indexable{
			Args:    args,
			Crawler: c,
		}

		return i.CrawlHash()
	}

	return &queue.Worker{
		ErrChan: f.errChan,
		Func:    fileFunc,
		Queue:   fileConQueue,
	}, nil
}
