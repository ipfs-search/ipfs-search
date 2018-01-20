package crawlworker

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"time"
)

type Config struct {
	IpfsAPI          string
	ElasticSearchURL string
	AMQPURL          string
	HashWorkers      uint
	FileWorkers      uint
	IpfsTimeout      time.Duration // Timeout for IPFS gateway HTTPS requests

	CrawlerConfig *crawler.Config

	// Temporarily disabled
	// HashWait         time.Duration // Time to wait between creating hash workers
	// FileWait         time.Duration // Time to wait between creating file workers
}
