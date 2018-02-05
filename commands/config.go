package commands

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/crawler/factory"
	"time"
)

// Config defines configuration for commands
type Config struct {
	Crawler *crawler.Config
	Factory *factory.Config

	HashWait time.Duration // Time to wait between creating hash workers
	FileWait time.Duration // Time to wait between creating file workers

	HashWorkers uint // Amount of workers for the hash queue
	FileWorkers uint // Amount of workers for the file queue
}

func getConfig() (*Config, error) {
	crawlerConfig := &crawler.Config{
		IpfsTikaURL:     "http://localhost:8081",
		IpfsTikaTimeout: 300 * time.Duration(time.Second),
		RetryWait:       2 * time.Duration(time.Second),
		MetadataMaxSize: 50 * 1024 * 1024,
		PartialSize:     262144,
	}

	factoryConfig := &factory.Config{
		IpfsAPI:          "localhost:5001",
		ElasticSearchURL: "http://localhost:9200",
		AMQPURL:          "amqp://guest:guest@localhost:5672/",
		IpfsTimeout:      360 * time.Duration(time.Second),
		CrawlerConfig:    crawlerConfig,
	}

	config := &Config{
		Crawler:     crawlerConfig,
		Factory:     factoryConfig,
		HashWait:    time.Duration(100 * time.Millisecond),
		FileWait:    time.Duration(100 * time.Millisecond),
		HashWorkers: 140,
		FileWorkers: 120,
	}

	return config, nil
}
