package redis

import (
	"github.com/rueian/rueidis"
	"github.com/rueian/rueidis/rueidisotel"

	"github.com/ipfs-search/ipfs-search/instr"
)

// Client represents a Redis client.
type Client struct {
	cfg *ClientConfig

	rueidis.Client
	*instr.Instrumentation
}

// ClientConfig contains configuration for the Redis client.
type ClientConfig struct {
	Addrs  []string // Address or addresses of a Redis node/cluster.
	Prefix string   // Prefix for storing entries.
}

// NewClient instantiates a new Redis client.
func NewClient(cfg *ClientConfig, i *instr.Instrumentation) (*Client, error) {
	if cfg == nil {
		panic("NewClient ClientConfig cannot be nil.")
	}

	if i == nil {
		panic("NewCLient Instrumentation cannot be nil.")
	}

	c, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: cfg.Addrs,
		ShuffleInit: true, // Recommended for cluster.
		// To connect to sentinels, specify the required master set name:
		// Sentinel: rueidis.SentinelOption{
		//     MasterSet: "my_master",
		// },
	})

	// Enable OpenTelemetry Tracing
	c = rueidisotel.WithClient(c)

	if err != nil {
		return nil, err
	}

	return &Client{
		cfg,
		c,
		i,
	}, nil
}
