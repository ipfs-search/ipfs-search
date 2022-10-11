package redis

import (
	"context"
	"log"

	"github.com/ipfs-search/ipfs-search/components/index"

	radix "github.com/mediocregopher/radix/v4"
)

// ExistsIndex stores properties as JSON in Redis.
type ExistsIndex struct {
	cfg *Config
	c   *Client
	key string
}

// NewExistsIndex returns a new index.
func NewExistsIndex(client *Client, cfg *Config) index.Index {
	if client == nil {
		panic("NewExistsIndex Client cannot be nil.")
	}

	if cfg == nil {
		panic("NewExistsIndex Config cannot be nil.")
	}

	index := &ExistsIndex{
		c:   client,
		cfg: cfg,
		key: client.cfg.Prefix + "e:" + cfg.Prefix,
	}

	return index
}

func (i *ExistsIndex) set(ctx context.Context, id string, properties interface{}) error {
	if debug {
		log.Printf("redis exists %s: add %s to %s", i, id, i.key)
	}

	action := radix.Cmd(nil, "SADD", i.key, id)
	return i.c.radixClient.Do(ctx, action)
}

// String returns the name of the index, for convenient logging.
func (i *ExistsIndex) String() string {
	return i.cfg.Name
}

// Index a document's properties, identified by id
func (i *ExistsIndex) Index(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.c.Tracer.Start(ctx, "index.redis.Index")
	defer span.End()

	return i.set(ctx, id, properties)
}

// Update a document's properties, given id
func (i *ExistsIndex) Update(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.c.Tracer.Start(ctx, "index.redis.Update")
	defer span.End()

	return i.set(ctx, id, properties)
}

// Delete item from set.
func (i *ExistsIndex) Delete(ctx context.Context, id string) error {
	ctx, span := i.c.Tracer.Start(ctx, "index.redis.Delete")
	defer span.End()

	if debug {
		log.Printf("redis exists %s: delete %s from %s", i, id, i.key)
	}

	// Non-blocking DEL-equivalent
	action := radix.Cmd(nil, "SREM", i.key, id)
	return i.c.radixClient.Do(ctx, action)
}

// Get returns whether or not an item is found (but doesn't update its properties).
func (i *ExistsIndex) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	ctx, span := i.c.Tracer.Start(ctx, "index.redis.Get")
	defer span.End()

	var found bool
	action := radix.Cmd(&found, "SISMEMBER", i.key, id)
	err := i.c.radixClient.Do(ctx, action)

	if debug {
		log.Printf("redis exists %s: get %s from %s, res: %v", i, id, i.key, found)
	}

	return found, err
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &ExistsIndex{}
