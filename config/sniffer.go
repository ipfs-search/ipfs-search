package config

import (
	"github.com/ipfs-search/ipfs-search/sniffer"
	"time"
)

// Sniffer is configuration pertaining to the sniffer
type Sniffer struct {
	LastSeenExpiration time.Duration `yaml:"lastseen_expiration"`
	LastSeenPruneLen   int           `yaml:"lastseen_prunelen"`
	LoggerTimeout      time.Duration `yaml:"logger_timeout"`
}

func (c *Config) SnifferConfig() *sniffer.Config {
	return &sniffer.Config{
		LastSeenExpiration: c.Sniffer.LastSeenExpiration,
		LastSeenPruneLen:   c.Sniffer.LastSeenPruneLen,
		LoggerTimeout:      c.Sniffer.LoggerTimeout,
	}
}

func SnifferDefaults() Sniffer {
	return Sniffer(*sniffer.DefaultConfig())
}
