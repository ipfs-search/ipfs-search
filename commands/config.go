package commands

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/crawler/factory"
	"time"
)

func getConfig() (*factory.Config, error) {
	crawlerConfig := &crawler.Config{
		IpfsTikaURL:     "http://localhost:8081",
		IpfsTikaTimeout: 300 * time.Duration(time.Second),
		RetryWait:       2 * time.Duration(time.Second),
		MetadataMaxSize: 50 * 1024 * 1024,
		PartialSize:     262144,
	}

	config := &factory.Config{
		IpfsAPI:          "localhost:5001",
		ElasticSearchURL: "http://localhost:9200",
		AMQPURL:          "amqp://guest:guest@localhost:5672/",
		HashWorkers:      140,
		FileWorkers:      120,
		IpfsTimeout:      360 * time.Duration(time.Second),
		// HashWait:         time.Duration(100 * time.Millisecond),
		// FileWait:         time.Duration(100 * time.Millisecond),
		CrawlerConfig: crawlerConfig,
	}

	return config, nil
}
