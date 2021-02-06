package sniffer

import "time"

// Config holds configuration for a Sniffer.
type Config struct {
	LastSeenExpiration time.Duration // Expiration time for the last-seen resources
	LastSeenPruneLen   int           // Cleanup expired resources from the last-seen
	LoggerTimeout      time.Duration // Throw timeout error when no log messages arrive
	BufferSize         uint          // Size of the channels buffering between yielder, filter and adder
}

// DefaultConfig returns the default configuration for a Sniffer.
func DefaultConfig() *Config {
	return &Config{
		LastSeenExpiration: 60 * time.Duration(time.Minute),
		LastSeenPruneLen:   32768,
		LoggerTimeout:      60 * time.Duration(time.Second),
		BufferSize:         512,
	}
}
