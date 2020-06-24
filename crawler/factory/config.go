package factory

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/index"
	"time"
)

// Config defines configuration for a crawler factory
type Config struct {
	IpfsAPI          string
	ElasticSearchURL string
	AMQPURL          string
	IpfsTimeout      time.Duration // Timeout for IPFS gateway HTTPS requests
	Indexes          map[string]*index.Config

	CrawlerConfig *crawler.Config
}
