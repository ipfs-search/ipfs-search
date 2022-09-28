package cache

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/instr"
)

const debug bool = false

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

	for _, field := range fields {
		if contains(i.cfg.CachingFields, field.Name) {
			value := valueof.FieldByName(field.Name).Interface()

			dst[field.Name] = value
		}
	}

	return dst
}

func (i *Index) cacheGet(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	var (
		found bool
		err   error
	)

	// Ignore fields for now; the OpenSearch API uses the json field names
	// Ref: https://github.com/ipfs-search/ipfs-search/issues/234

	if found, err = i.cachingIndex.Get(ctx, id, dst, fields...); err != nil {
		// Ignore context closed
		err = ErrCache{err, fmt.Sprintf("cache error in get: %s", err.Error())}
	}

	return found, err
}

type indexWrite func(context.Context, string, interface{}) error

func (i *Index) cacheWrite(ctx context.Context, id string, properties interface{}, f indexWrite) error {
	cachingProperties := i.makeCachingProperties(properties)

	if debug {
		log.Printf("cache: write %s", id)
	}

	if err := f(ctx, id, cachingProperties); err != nil {
		return ErrCache{err, fmt.Sprintf("cache error in '%v': %s", f, err.Error())}
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

	if debug {
		log.Printf("cache: delete %s", id)
	}

	// Delete cache first; maintain consistency as our backing index is the source of truth.
	if err := i.cachingIndex.Delete(ctx, id); err != nil {
		return ErrCache{
			err,
			fmt.Sprintf("error deleting cache: %s", err.Error()),
		}
	}

	if err := i.backingIndex.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// Get retreives *all* fields from document with `id` from the cache, falling back to the backing index.
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
		// if debug {
		log.Printf("cache: hit %s", id)
		// }

		return found, err
	}

	if debug {
		log.Printf("cache: miss %s", id)
	}

	var backingErr error
	if found, backingErr = i.backingIndex.Get(ctx, id, dst, fields...); backingErr != nil {
		// Backing errors overwrite cache errors.
		err = backingErr
	}

	if found {
		if debug {
			log.Printf("backing: hit %s", id)
		}

		if indexErr := i.cacheWrite(ctx, id, dst, i.cachingIndex.Index); indexErr != nil {
			err = indexErr
		}
	}

	if debug {
		log.Printf("backing: miss %s", id)
	}

	return found, err
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Index{}
