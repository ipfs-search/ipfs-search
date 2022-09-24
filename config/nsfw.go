package config

import (
	"time"

	"github.com/c2h5oh/datasize"

	"github.com/ipfs-search/ipfs-search/components/extractor/nsfw"
)

// NSFW is configuration pertaining to the sniffer
type NSFW struct {
	NSFWServerURL  string            `yaml:"url" env:"NSFW_URL"`
	RequestTimeout time.Duration     `yaml:"timeout"`
	MaxFileSize    datasize.ByteSize `yaml:"max_file_size"`
}

// NSFWConfig returns component-specific configuration from the canonical central configuration.
func (c *Config) NSFWConfig() *nsfw.Config {
	cfg := nsfw.Config(c.NSFW)
	return &cfg
}

// NSFWDefaults returns the defaults for component configuration, based on the component-specific configuration.
func NSFWDefaults() NSFW {
	return NSFW(*nsfw.DefaultConfig())
}
