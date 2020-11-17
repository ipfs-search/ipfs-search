package tika

import (
	"time"
)

// Config specifies the configuration for a Tika extractor.
type Config struct {
	TikaServerURL  string        // TikaServer is the URL of the ipfs-tika server.
	RequestTimeout time.Duration // Timeout for metadata requests for the server.
	RetryWait      time.Duration // Wait time between retries of failed requests.
}

// DefaultConfig returns the default configuration for a Sniffer.
func DefaultConfig() *Config {
	return &Config{
		TikaServerURL:  "http://localhost:8081",
		RequestTimeout: 300 * time.Duration(time.Second),
		RetryWait:      2 * time.Duration(time.Second),
	}
}
