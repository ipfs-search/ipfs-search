package config

import (
	"time"

	"github.com/c2h5oh/datasize"

	"github.com/ipfs-search/ipfs-search/extractor/tika"
)

// Tika is configuration pertaining to the sniffer
type Tika struct {
	TikaServerURL  string            `yaml:"url" env:"IPFS_TIKA_URL"`
	RequestTimeout time.Duration     `yaml:"timeout"`
	MaxFileSize    datasize.ByteSize `yaml:"max_file_size"`
}

// TikaConfig returns component-specific configuration from the canonical central configuration.
func (c *Config) TikaConfig() *tika.Config {
	cfg := tika.Config(c.Tika)
	return &cfg
}

// TikaDefaults returns the defaults for component configuration, based on the component-specific configuration.
func TikaDefaults() Tika {
	return Tika(*tika.DefaultConfig())
}
