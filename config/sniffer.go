package config

import (
	"github.com/ipfs-search/ipfs-search/sniffer"
	"time"
)

type Sniffer struct {
	LastSeenExpiration time.Duration `yaml:"lastseen_expiration"`
	LastSeenPruneLen   int           `yaml:"lastseen_prunelen"`
	LoggerTimeout      time.Duration `yaml:"logger_timeout"`
}

func (c *Config) SnifferConfig() *sniffer.Config {
	return &sniffer.Config{
		IpfsAPI:            c.IPFS.IpfsAPI,
		AMQPURL:            c.AMQP.AMQPURL,
		LastSeenExpiration: c.Sniffer.LastSeenExpiration,
		LastSeenPruneLen:   c.Sniffer.LastSeenPruneLen,
		LoggerTimeout:      c.Sniffer.LoggerTimeout,
	}
}

func SnifferDefaults() Sniffer {
	return Sniffer{
		LastSeenExpiration: 60 * time.Duration(time.Minute),
		LastSeenPruneLen:   16383,
		LoggerTimeout:      60 * time.Duration(time.Second),
	}
}
