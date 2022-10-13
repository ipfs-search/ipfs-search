package redis

import (
	"context"
	"log"
	"strings"

	radix "github.com/mediocregopher/radix/v4"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/instr"
)

// ClientConfig contains configuration for the Redis client.
type ClientConfig struct {
	Addrs  []string // Address or addresses of a Redis node/cluster.
	Prefix string   // Prefix for storing entries.
}

// Client represents a Redis client.
type Client struct {
	cfg *ClientConfig
	*instr.Instrumentation

	radixClient radix.MultiClient
}

// NewClient instantiates a new Redis client.
func NewClient(cfg *ClientConfig, i *instr.Instrumentation) (*Client, error) {
	if cfg == nil {
		panic("NewClient ClientConfig cannot be nil.")
	}

	if len(cfg.Addrs) == 0 {
		panic("No Redis addresses specified.")
	}

	if i == nil {
		panic("NewCLient Instrumentation cannot be nil.")
	}

	return &Client{
		cfg:             cfg,
		Instrumentation: i,
	}, nil
}

func isClusterNotSupportedError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "cluster support disabled") || strings.Contains(s, "unknown command")
}

// Start starts radix connections and closes them when done.
func (c *Client) Start(ctx context.Context) error {
	var err error
	if c.radixClient, err = (radix.ClusterConfig{}).New(ctx, c.cfg.Addrs); err != nil {
		if isClusterNotSupportedError(err) && len(c.cfg.Addrs) == 1 {
			log.Printf("Redis not a cluster, attempting single connection.")
			singleClient, err := (radix.PoolConfig{}).New(ctx, "tcp", c.cfg.Addrs[0])
			if err != nil {
				return err
			}

			c.radixClient = radix.NewMultiClient(radix.ReplicaSet{
				Primary: singleClient,
			})
		} else {
			return err
		}
	}

	return nil
}

// NewIndex returns a new index given with given name and prefix. When existsIndex is true, an ExistsIndex will be returned.
func (c *Client) NewIndex(name, prefix string, existsIndex bool) index.Index {
	if existsIndex {
		return NewExistsIndex(
			c,
			&Config{Name: name, Prefix: prefix},
		)
	}

	return New(
		c,
		&Config{Name: name, Prefix: prefix},
	)
}

// Close closes the Redis client connection.
func (c *Client) Close(ctx context.Context) error {
	return c.radixClient.Close()
}
