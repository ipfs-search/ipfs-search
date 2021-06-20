package config

import (
	"github.com/ipfs-search/ipfs-search/components/queue/amqp"
	"time"
)

// AMQP contains configuration pertaining to AMQP.
type AMQP struct {
	URL           string        `yaml:"url" env:"AMQP_URL"`                 // URL of AMQP server.
	MaxReconnect  int           `yaml:"max_reconnect"`                      // The maximum number of reconnection attempts after the server connection is lost.
	ReconnectTime time.Duration `yaml:"reconnect_time"`                     // The time to wait in between reconnect attempts.
	MessageTTL    time.Duration `yaml:"message_ttl" env:"AMQP_MESSAGE_TTL"` // The expiration time for messages in the queue.
}

// AMQPConfig returns component-specific configuration from the canonical configuration.
func (c *Config) AMQPConfig() *amqp.Config {
	cfg := amqp.Config(c.AMQP)
	return &cfg
}

// AMQPDefaults returns the defaults for component configuration, based on the component-specific configuration.
func AMQPDefaults() AMQP {
	return AMQP(*amqp.DefaultConfig())
}
