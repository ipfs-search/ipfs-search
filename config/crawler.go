package config

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"time"
)

// Crawler contains configuration for a Crawler.
type Crawler struct {
	DirEntryBufferSize uint          // Size of buffer for processing directory entry channels.
	MinUpdateAge       time.Duration // The minimum age for items to be updated.
	StatTimeout        time.Duration // Timeout for Stat() calls.
	DirEntryTimeout    time.Duration // Timeout *between* directory entries.
}

func (c *Config) CrawlerConfig() *crawler.Config {
	cfg := crawler.Config(c.Crawler)
	return &cfg
}

func CrawlerDefaults() Crawler {
	return Crawler(*crawler.DefaultConfig())
}
