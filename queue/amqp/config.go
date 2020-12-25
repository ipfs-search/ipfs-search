package amqp

import (
	"time"
)

type Config struct {
	URL           string
	MaxReconnect  int
	ReconnectTime time.Duration
}

// DefaultConfig generates a default configuration.
func DefaultConfig() *Config {
	return &Config{
		URL:           "amqp://guest:guest@localhost:5672/",
		MaxReconnect:  100,
		ReconnectTime: 2 * time.Second,
	}
}
