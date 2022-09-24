package redis

import (
	"context"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/rueian/rueidis"
)

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
	return i.c.cfg.Prefix + i.cfg.Name + ":" + key
}

func (i *Index) set(ctx context.Context, id string, properties interface{}) error {
	key := i.getKey(id)
	val := rueidis.JSON(properties)

	cmd := i.c.B().Set().Key(key).Value(val).Build()
	res := i.c.Do(ctx, cmd)

	return res.Error()
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
	cmd := i.c.B().Del().Key(key).Build()
	res := i.c.Do(ctx, cmd)

	return res.Error()
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
	cmd := i.c.B().Get().Key(key).Build()
	res := i.c.Do(ctx, cmd)

	if err := res.Error(); err != nil {
		return false, err
	}

	if res.RedisError().IsNil() {
		return false, nil
	}

	return true, res.DecodeJSON(dst)
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Index{}
