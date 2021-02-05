package amqp

import (
	"time"
)

// Config specifies the configuration for an AMQP queue.
type Config struct {
	URL           string
	MaxReconnect  int
	ReconnectTime time.Duration
}

// DefaultConfig generates a default configuration for an AMQP queue.
func DefaultConfig() *Config {
	return &Config{
		URL:           "amqp://guest:guest@localhost:5672/",
		MaxReconnect:  100,
		ReconnectTime: 2 * time.Second,
	}
}
