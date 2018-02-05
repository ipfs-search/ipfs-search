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

	CrawlerConfig *crawler.Config
}
