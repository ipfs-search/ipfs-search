package factory

import (
	"time"

	"github.com/ipfs-search/ipfs-search/crawler"
	tika "github.com/ipfs-search/ipfs-search/extractor/tika"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"
)

// Config defines configuration for a crawler factory
type Config struct {
	IpfsAPI          string
	ElasticSearchURL string
	AMQPURL          string
	IpfsTimeout      time.Duration // Timeout for IPFS gateway HTTPS requests
	Indexes          map[string]*elasticsearch.Config

	ExtractorConfig *tika.Config
	CrawlerConfig   *crawler.Config
}
