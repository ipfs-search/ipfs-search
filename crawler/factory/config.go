package factory

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"time"
)

// Config defines configuration for a crawler factory
type Config struct {
	IpfsAPI          string
	ElasticSearchURL string
	AMQPURL          string
	IpfsTimeout      time.Duration // Timeout for IPFS gateway HTTPS requests
	Indexes          map[string]*IndexConfig

	CrawlerConfig *crawler.Config
}

// IndexConfig represents the configuration for a specific index.
type IndexConfig struct {
	Name     string
	Settings map[string]interface{}
	Mapping  map[string]interface{}
}
