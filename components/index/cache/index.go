package cache

import (
	"context"
	"fmt"
	"reflect"

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
func New(backing index.Index, caching index.Index,
	cfg *Config, instr *instr.Instrumentation) index.Index {

	if cfg == nil {
		panic("Index.New Config cannot be nil.")
	}

	index := &Index{
		backingIndex:    backing,
		cachingIndex:    caching,
		cfg:             cfg,
		Instrumentation: instr,
	}

	return index
}

// String returns the name of the index, for convenient logging.
func (i *Index) String() string {
	return fmt.Sprintf("cache for '%s' through '%s'", i.backingIndex, i.cachingIndex)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (i *Index) makeCachingProperties(properties interface{}) map[string]interface{} {
	// Take care to allocate map for caching properties on the stack.

	valueof := reflect.ValueOf(properties)

	if valueof.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("not called with pointer but %T", properties))
	}

	// Dereference pointer
	valueof = valueof.Elem()

	if valueof.Kind() != reflect.Struct {
		panic(fmt.Sprintf("not struct pointer but %T", properties))
	}

	dst := make(map[string]interface{}, len(i.cfg.CachingFields))
	fields := reflect.VisibleFields(valueof.Type())

	for k, field := range fields {
		if contains(i.cfg.CachingFields, field.Name) {
			dst[field.Name] = valueof.Field(k).Interface()
		}
	}

	if len(dst) == 0 {
		panic("no cachable properties found")
	}

	return dst
}

func (i *Index) allFieldsCachable(fields []string) bool {
	for _, field := range fields {
		exists := contains(i.cfg.CachingFields, field)
		if !exists {
			return false
		}
	}

	return true
}

func (i *Index) cacheGet(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	var (
		found bool
		err   error
	)

	if i.allFieldsCachable(fields) {
		if found, err = i.cachingIndex.Get(ctx, id, dst, fields...); err != nil {
			err = ErrCache{err, fmt.Sprintf("cache error deleting: %e", err)}
		}
	} else {
	}

	return found, err
}

type indexWrite func(context.Context, string, interface{}) error

func (i *Index) cacheWrite(ctx context.Context, id string, properties interface{}, f indexWrite) error {
	cachingProperties := i.makeCachingProperties(properties)

	if err := f(ctx, id, cachingProperties); err != nil {
		return ErrCache{err, fmt.Sprintf("cache error in '%v': %e", f, err)}
	}

	return nil
}

// Index a document's properties, identified by id.
// Returns error of type ErrCache if the caching index returned an error.
func (i *Index) Index(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.Tracer.Start(ctx, "index.cache.Index")
	defer span.End()

	// Backing index first.
	if err := i.backingIndex.Index(ctx, id, properties); err != nil {
		return err
	}

	return i.cacheWrite(ctx, id, properties, i.cachingIndex.Index)
}

// Update a document's properties, given id.
// Returns error of type ErrCache if the caching index returned an error.
func (i *Index) Update(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.Tracer.Start(ctx, "index.cache.Update")
	defer span.End()

	// Update cache first and backing index later.
	if err := i.cacheWrite(ctx, id, properties, i.cachingIndex.Update); err != nil {
		return err
	}

	if err := i.backingIndex.Update(ctx, id, properties); err != nil {
		return err
	}

	return nil
}

// Delete item from index.
// Returns error of type ErrCache if the caching index returned an error.
func (i *Index) Delete(ctx context.Context, id string) error {
	ctx, span := i.Tracer.Start(ctx, "index.cache.Delete")
	defer span.End()

	// Delete cache first; maintain consistency as our backing index is the source of truth.
	if err := i.cachingIndex.Delete(ctx, id); err != nil {
		return ErrCache{err, "error deleting cache"}
	}

	if err := i.backingIndex.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// Get retreives *all* fields from document with `id` from the cache, falling back to the backing index.
// `fields` parameter is used to determine whether cache can be used based on configured CachingFields.
// Returns: (exists, err) where err is of type ErrCache if there was (only) an error from the
// caching index.
func (i *Index) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	ctx, span := i.Tracer.Start(ctx, "index.opensearch.Get")
	defer span.End()

	var (
		found bool
		err   error
	)

	if found, err = i.cacheGet(ctx, id, dst, fields...); found {
		return found, err
	}

	var backingErr error
	if found, backingErr = i.backingIndex.Get(ctx, id, dst, fields...); backingErr != nil {
		// Backing errors overwrite cache errors.
		err = backingErr
	}

	if found {
		if indexErr := i.cacheWrite(ctx, id, dst, i.cachingIndex.Index); indexErr != nil {
			err = indexErr
		}
	}

	return found, err
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Index{}
