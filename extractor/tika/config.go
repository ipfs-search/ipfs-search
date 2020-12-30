package tika

import (
	// "github.com/c2h5oh/datasize"
	"time"
)

// Config specifies the configuration for a Tika extractor.
type Config struct {
	TikaServerURL  string        // TikaServer is the URL of the ipfs-tika server.
	RequestTimeout time.Duration // Timeout for metadata requests for the server.
	// MetadataMaxSize datasize.ByteSize `yaml:"metadata_max_size"` // TODO
}

// DefaultConfig returns the default configuration for a Sniffer.
func DefaultConfig() *Config {
	return &Config{
		TikaServerURL:  "http://localhost:8081",
		RequestTimeout: 300 * time.Duration(time.Second),
		// MetadataMaxSize: 50 * 1024 * 1024,
	}
}
