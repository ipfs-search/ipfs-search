package config

import (
	"time"
)

// Default() returns default configuration
func Default() *Config {
	return &Config{
		IPFS{
			IpfsAPI:     "localhost:5001",
			IpfsTimeout: 360 * time.Duration(time.Second),
		},
		ElasticSearch{
			ElasticSearchURL: "http://localhost:9200",
		},
		AMQP{
			AMQPURL: "amqp://guest:guest@localhost:5672/",
		},
		CrawlerDefaults(),
		SnifferDefaults(),
		IndexesDefaults(),
		ExtractorDefaults(),
	}
}
