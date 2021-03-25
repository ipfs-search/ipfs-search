package tika

import (
	"github.com/c2h5oh/datasize"
	"time"
)

// Config specifies the configuration for a Tika extractor.
type Config struct {
	TikaExtractorURL string            // TikaServer is the URL of the ipfs-tika server.
	RequestTimeout   time.Duration     // Timeout for metadata requests for the server.
	MaxFileSize      datasize.ByteSize // Don't attempt to get metadata for files over this size.
}

// DefaultConfig returns the default configuration for a Sniffer.
func DefaultConfig() *Config {
	return &Config{
		TikaExtractorURL: "http://localhost:8081",
		RequestTimeout:   300 * time.Duration(time.Second),
		MaxFileSize:      4 * 1024 * 1024 * 1024, // 4GB
	}
}
