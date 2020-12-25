package config

import (
	"github.com/ipfs-search/ipfs-search/queue/amqp"
	"time"
)

// AMQP contains configuration pertaining to AMQP.
type AMQP struct {
	URL           string        `yaml:"url" env:"AMQP_URL"`
	MaxReconnect  int           `yaml:"max_reconnect"`
	ReconnectTime time.Duration `yaml:"reconnect_time"`
}

func (c *Config) AMQPConfig() *amqp.Config {
	cfg := amqp.Config(c.AMQP)
	return &cfg
}

func AMQPDefaults() AMQP {
	return AMQP(*amqp.DefaultConfig())
}
