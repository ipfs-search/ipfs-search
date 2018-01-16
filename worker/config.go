package worker

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"time"
)

// Config contains user configurable options for a worker
type Config struct {
	IpfsAPI          string
	ElasticSearchURL string
	HashWorkers      int
	FileWorkers      int
	IpfsTimeout      time.Duration // Timeout for IPFS gateway HTTPS requests
	HashWait         time.Duration // Time to wait between creating hash workers
	FileWait         time.Duration // Time to wait between creating file workers
	CrawlerConfig    *crawler.Config
}
