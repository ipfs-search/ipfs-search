package factory

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"
	"time"
)

// Config defines configuration for a crawler factory
type Config struct {
	IpfsAPI          string
	ElasticSearchURL string
	AMQPURL          string
	IpfsTimeout      time.Duration // Timeout for IPFS gateway HTTPS requests
	Indexes          map[string]*elasticsearch.Config

	CrawlerConfig *crawler.Config
}
