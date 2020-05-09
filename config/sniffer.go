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
	BufferSize         uint          `yaml:"buffer_size"`
}

func (c *Config) SnifferConfig() *sniffer.Config {
	cfg := sniffer.Config(c.Sniffer)
	return &cfg
}

func SnifferDefaults() Sniffer {
	return Sniffer(*sniffer.DefaultConfig())
}
