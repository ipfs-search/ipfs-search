package config

import (
	"github.com/ipfs-search/ipfs-search/extractor/tika"
	"time"
)

// Extractor is configuration pertaining to the sniffer
type Extractor struct {
	TikaServerURL  string        `yaml:"url" env:"IPFS_TIKA_URL"`
	RequestTimeout time.Duration `yaml:"timeout"`
	RetryWait      time.Duration `yaml:"retry_wait"`
}

func (c *Config) ExtractorConfig() *tika.Config {
	cfg := tika.Config(c.Extractor)
	return &cfg
}

func ExtractorDefaults() Extractor {
	return Extractor(*tika.DefaultConfig())
}
