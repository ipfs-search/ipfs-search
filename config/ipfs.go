package config

import (
	"github.com/c2h5oh/datasize"
	"github.com/ipfs-search/ipfs-search/components/protocol/ipfs"
)

// IPFS specifies the configuration for the IPFS protocol.
type IPFS struct {
	APIURL      string            `yaml:"api_url" env:"IPFS_API_URL"`
	GatewayURL  string            `yaml:"gateway_url"`
	PartialSize datasize.ByteSize `yaml:"partial_size"`
}

// IPFSConfig returns component-specific configuration from the canonical central configuration.
func (c *Config) IPFSConfig() *ipfs.Config {
	cfg := ipfs.Config(c.IPFS)
	return &cfg
}

// IPFSDefaults returns the defaults for component configuration, based on the component-specific configuration.
func IPFSDefaults() IPFS {
	return IPFS(*ipfs.DefaultConfig())
}
