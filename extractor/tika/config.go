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
