package config

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/crawler/factory"
	"time"
)

type Tika struct {
	IpfsTikaURL     string        `yaml:"url"`
	IpfsTikaTimeout time.Duration `yaml:"timeout"`
	MetadataMaxSize uint64        `yaml:"max_size"`
}

type IPFS struct {
	IpfsAPI     string        `yaml:"api_url"`
	IpfsTimeout time.Duration `yaml:"timeout"`
}

type ElasticSearch struct {
	ElasticSearchURL string `yaml:"url"`
}

type AMPQ struct {
	AMQPURL string `yaml:"url"`
}

type Crawler struct {
	RetryWait   time.Duration `yaml:"retry_wait"`
	HashWait    time.Duration `yaml:"hash_wait"`
	FileWait    time.Duration `yaml:"file_wait"`
	PartialSize uint64        `yaml:"partial_size"`
	HashWorkers uint          `yaml:"hash_workers"`
	FileWorkers uint          `yaml:"file_workers"`
}

type Config struct {
	Tika          Tika          `yaml:"tika"`
	IPFS          IPFS          `yaml:"ipfs"`
	ElasticSearch ElasticSearch `yaml:"elasticsearch"`
	AMPQ          AMPQ          `yaml:"ampq"`
	Crawler       Crawler       `yaml:"crawler"`
}

func (c *Config) CrawlerConfig() *crawler.Config {
	return &crawler.Config{
		IpfsTikaURL:     c.Tika.IpfsTikaURL,
		IpfsTikaTimeout: c.Tika.IpfsTikaTimeout,
		MetadataMaxSize: c.Tika.MetadataMaxSize,
		RetryWait:       c.Crawler.RetryWait,
		PartialSize:     c.Crawler.PartialSize,
	}
}

func (c *Config) FactoryConfig() *factory.Config {
	return &factory.Config{
		IpfsAPI:          c.IPFS.IpfsAPI,
		IpfsTimeout:      c.IPFS.IpfsTimeout,
		ElasticSearchURL: c.ElasticSearch.ElasticSearchURL,
		AMQPURL:          c.AMPQ.AMQPURL,
		CrawlerConfig:    c.CrawlerConfig(),
	}
}

func Get() (*Config, error) {
	return Default(), nil
}
