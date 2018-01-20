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
	crawlerConfig *crawler.Config
	pubConnection *queue.Connection
	conConnection *queue.Connection
	errChan       chan<- error
	indexer       *indexer.Indexer
	shell         *shell.Shell
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

func (f *Factory) NewCrawler() (*crawler.Crawler, error) {
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

func (f *Factory) NewHashWorker() (worker.Worker, error) {
	hashConQueue, err := f.conConnection.NewChannelQueue("hashes")
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
	fileConQueue, err := f.conConnection.NewChannelQueue("files")
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

		return i.CrawlFile()
	}

	return &queue.Worker{
		ErrChan: f.errChan,
		Func:    fileFunc,
		Queue:   fileConQueue,
	}, nil
}
