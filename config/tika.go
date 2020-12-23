package config

import (
	"github.com/ipfs-search/ipfs-search/extractor/tika"
	"time"
)

// Tika is configuration pertaining to the sniffer
type Tika struct {
	TikaServerURL  string        `yaml:"url" env:"IPFS_TIKA_URL"`
	RequestTimeout time.Duration `yaml:"timeout"`
	RetryWait      time.Duration `yaml:"retry_wait"`
}

func (c *Config) TikaConfig() *tika.Config {
	cfg := tika.Config(c.Tika)
	return &cfg
}

func TikaDefaults() Tika {
	return Tika(*tika.DefaultConfig())
}
