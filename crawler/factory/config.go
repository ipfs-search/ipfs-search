package factory

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"time"
)

type Config struct {
	IpfsAPI          string
	ElasticSearchURL string
	AMQPURL          string
	IpfsTimeout      time.Duration // Timeout for IPFS gateway HTTPS requests

	CrawlerConfig *crawler.Config
}
