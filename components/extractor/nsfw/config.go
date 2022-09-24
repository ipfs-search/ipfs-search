package nsfw

import (
	"github.com/c2h5oh/datasize"
	"time"
)

// Config specifies the configuration for NSFW server.
type Config struct {
	NSFWServerURL  string            // URL of nsfw-server
	RequestTimeout time.Duration     // Timeout for metadata requests for the server.
	MaxFileSize    datasize.ByteSize // Don't attempt to get metadata for files over this size.
}

// DefaultConfig returns the default configuration for a Sniffer.
func DefaultConfig() *Config {
	return &Config{
		NSFWServerURL:  "http://localhost:3000",
		RequestTimeout: 300 * time.Duration(time.Second),
		MaxFileSize:    1 * datasize.GB,
	}
}
