package ipfs

import (
	"github.com/c2h5oh/datasize"
)

// Config specifies the configuration for the IPFS protocol.
type Config struct {
	APIURL      string            // URL of an IPFS API endpoint (for Ls and Stat calls).
	GatewayURL  string            // URL of an IPFS Gateway (to request content).
	PartialSize datasize.ByteSize // Filesize of items which are being considered partials (chunks).
}

// DefaultConfig returns the default configuration for a Sniffer.
func DefaultConfig() *Config {
	return &Config{
		APIURL:      "http://localhost:5001",
		GatewayURL:  "http://localhost:8080",
		PartialSize: 262144,
		// 256KB is the default chunker block size. Therefore, unreferenced files with exactly
		// this size are very likely to be chunks of files (partials) rather than full files.
	}
}
