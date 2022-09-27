package cache

import (
	"context"
	"fmt"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/instr"
)

// Index wraps a backing index and caches it using another index.
type Index struct {
	cfg          *Config
	backingIndex index.Index
	cachingIndex index.Index

	*instr.Instrumentation
}

// New returns a new index.
func New(backing index.Index, caching index.Index, cfg *Config) index.Index {

	if cfg == nil {
		panic("Index.New Config cannot be nil.")
	}

	index := &Index{
		backingIndex: backing,
		cachingIndex: caching,
		cfg:          cfg,
	}

	return index
}

// String returns the name of the index, for convenient logging.
func (i *Index) String() string {
	return fmt.Sprintf("cache for %s through %s", i.backingIndex, i.cachingIndex)
}

func makeCachingProperties(properties interface{}) interface{} {
	// Take care to allocate map for caching properties on the stack.
	return nil
}

// Index a document's properties, identified by id
func (i *Index) Index(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.Tracer.Start(ctx, "index.cache.Index")
	defer span.End()

	if err := i.backingIndex.Index(ctx, id, properties); err != nil {
		return err
	}

	cachingProperties := makeCachingProperties(properties)
	return i.cachingIndex.Index(ctx, id, cachingProperties)
}

// Update a document's properties, given id
func (i *Index) Update(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.Tracer.Start(ctx, "index.cache.Update")
	defer span.End()

	if err := i.backingIndex.Update(ctx, id, properties); err != nil {
		return err
	}

	cachingProperties := makeCachingProperties(properties)
	return i.cachingIndex.Update(ctx, id, cachingProperties)
}

// Delete item from index
func (i *Index) Delete(ctx context.Context, id string) error {
	ctx, span := i.Tracer.Start(ctx, "index.cache.Delete")
	defer span.End()

	if err := i.backingIndex.Delete(ctx, id); err != nil {
		return err
	}

	return i.cachingIndex.Delete(ctx, id)
}

// Get retreives `fields` from document with `id` from the index, returning:
func (i *Index) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	ctx, span := i.Tracer.Start(ctx, "index.opensearch.Get")
	defer span.End()

	// First, try caching index
	found, err := i.cachingIndex.Get(ctx, id, dst)
	if err != nil {
		return false, err
	}

	if found {
		return true, nil
	}

	// Secondly, try backind index
	found, err = i.backingIndex.Get(ctx, id, dst)
	if err != nil {
		return false, err
	}

	if found {
		// Add to cache
		cachingProperties := makeCachingProperties(dst)
		return true, i.cachingIndex.Index(ctx, id, cachingProperties)
	}

	return false, nil
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Index{}
