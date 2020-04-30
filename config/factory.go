package config

import (
	"github.com/ipfs-search/ipfs-search/crawler/factory"
)

func (c *Config) FactoryConfig() *factory.Config {
	return &factory.Config{
		IpfsAPI:          c.IPFS.IpfsAPI,
		IpfsTimeout:      c.IPFS.IpfsTimeout,
		ElasticSearchURL: c.ElasticSearch.ElasticSearchURL,
		AMQPURL:          c.AMQP.AMQPURL,
		CrawlerConfig:    c.CrawlerConfig(),
	}
}
