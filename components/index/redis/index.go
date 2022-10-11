package redis

import (
	"context"
	"log"

	"github.com/ipfs-search/ipfs-search/components/index"

	radix "github.com/mediocregopher/radix/v4"
	"github.com/mediocregopher/radix/v4/resp/resp3"
)

const debug bool = true

// Index stores properties as JSON in Redis.
type Index struct {
	cfg *Config
	c   *Client
}

// New returns a new index.
func New(client *Client, cfg *Config) index.Index {
	if client == nil {
		panic("Index.New Client cannot be nil.")
	}

	if cfg == nil {
		panic("Index.New Config cannot be nil.")
	}

	index := &Index{
		c:   client,
		cfg: cfg,
	}

	return index
}

func (i *Index) getKey(key string) string {
	return i.c.cfg.Prefix + i.cfg.Prefix + ":" + key
}

func (i *Index) set(ctx context.Context, id string, properties interface{}) error {
	key := i.getKey(id)
	args := []string{key}

	flattened, err := resp3.Flatten(properties, nil)
	if err != nil {
		return err
	}

	if len(flattened) == 0 {
		panic("Redis cannot index without properties.")
	}

	if debug {
		log.Printf("redis %s: writing %+v to %s", i, flattened, key)
	}

	args = append(args, flattened...)

	action := radix.Cmd(nil, "HSET", args...)
	return i.c.radixClient.Do(ctx, action)
}

// String returns the name of the index, for convenient logging.
func (i *Index) String() string {
	return i.cfg.Name
}

// Index a document's properties, identified by id
func (i *Index) Index(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.c.Tracer.Start(ctx, "index.redis.Index")
	defer span.End()

	return i.set(ctx, id, properties)
}

// Update a document's properties, given id
func (i *Index) Update(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.c.Tracer.Start(ctx, "index.redis.Update")
	defer span.End()

	return i.set(ctx, id, properties)
}

// Delete item from index
func (i *Index) Delete(ctx context.Context, id string) error {
	ctx, span := i.c.Tracer.Start(ctx, "index.redis.Delete")
	defer span.End()

	key := i.getKey(id)

	if debug {
		log.Printf("redis %s: delete %s", i, key)
	}

	// Non-blocking DEL-equivalent
	action := radix.Cmd(nil, "UNLINK", key)
	return i.c.radixClient.Do(ctx, action)
}

// Get *all* fields from document with `id` from the index, ignoring the 'fields' parameters.
//
// Returs:
// - (true, decoding_error) if found (decoding error set when errors in json)
// - (false, nil) when not found
// - (false, error) otherwise
func (i *Index) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	ctx, span := i.c.Tracer.Start(ctx, "index.redis.Get")
	defer span.End()

	key := i.getKey(id)

	// Wrap receiver so we can determine whether we're found or not.
	mb := &radix.Maybe{Rcv: dst}
	action := radix.Cmd(mb, "HGETALL", key)
	err := i.c.radixClient.Do(ctx, action)

	// if debug {
	// 	log.Printf("redis %s: get %s", i, key)
	// 	log.Printf("redis %s: maybe: %+v", i, mb)
	// 	log.Printf("redis %s: dst: %T: %v", i, dst, dst)
	// }

	return !(mb.Null || mb.Empty || err != nil), err
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Index{}
