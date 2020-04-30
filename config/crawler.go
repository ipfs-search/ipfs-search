package config

import (
	"github.com/c2h5oh/datasize"
	"github.com/ipfs-search/ipfs-search/crawler"
	"time"
)

type Crawler struct {
	RetryWait   time.Duration     `yaml:"retry_wait"`
	HashWait    time.Duration     `yaml:"hash_wait"`
	FileWait    time.Duration     `yaml:"file_wait"`
	PartialSize datasize.ByteSize `yaml:"partial_size"`
	HashWorkers uint              `yaml:"hash_workers"`
	FileWorkers uint              `yaml:"file_workers"`
}

func (c *Config) CrawlerConfig() *crawler.Config {
	return &crawler.Config{
		IpfsTikaURL:     c.Tika.IpfsTikaURL,
		IpfsTikaTimeout: c.Tika.IpfsTikaTimeout,
		MetadataMaxSize: uint64(c.Tika.MetadataMaxSize),
		RetryWait:       c.Crawler.RetryWait,
		PartialSize:     uint64(c.Crawler.PartialSize),
	}
}

func CrawlerDefaults() Crawler {
	return Crawler{
		HashWait:    time.Duration(100 * time.Millisecond),
		FileWait:    time.Duration(100 * time.Millisecond),
		HashWorkers: 140,
		FileWorkers: 120,
		RetryWait:   2 * time.Duration(time.Second),
		PartialSize: 262144,
	}
}
