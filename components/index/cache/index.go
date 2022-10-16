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
	backingIndex index.Index
	cachingIndex index.Index
	cachingType  reflect.Type

	*instr.Instrumentation
}

// New returns a new index.
func New(backing index.Index, caching index.Index, cachingType interface{}, instr *instr.Instrumentation) index.Index {
	t := reflect.TypeOf(cachingType)
	if t.Kind() == reflect.Pointer {
		// Dereference pointer
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic("caching type should be a struct")
	}

	index := &Index{
		backingIndex:    backing,
		cachingIndex:    caching,
		cachingType:     t,
		Instrumentation: instr,
	}

	return index
}

// String returns the name of the index, for convenient logging.
func (i *Index) String() string {
	return fmt.Sprintf("'%s' through '%s'", i.backingIndex, i.cachingIndex)
}

func matchDstKind(src, dst reflect.Value) reflect.Value {
	dKind, sKind := dst.Kind(), src.Kind()

	if dKind == reflect.Pointer && sKind != reflect.Pointer {
		// dst is pointer, src is not.

		if !src.CanAddr() {
			panic(fmt.Sprintf("cannot address val %v for field %v", src, dst))
		}

		return src.Addr()
	}

	if dKind != reflect.Pointer && sKind == reflect.Pointer {
		// dst is value, src is pointer.
		return src.Elem()
	}

	return src
}

func setFieldVal(src, dst reflect.Value, dstField reflect.StructField) {
	// Set dst field to corresponding src value.
	// Note: this will panic when a dst field is not present in the src struct.
	srcVal := src.FieldByName(dstField.Name)
	dstVal := dst.FieldByIndex(dstField.Index)

	srcVal = matchDstKind(srcVal, dstVal)

	dstVal.Set(srcVal)
}

func (i *Index) makeCachingProperties(props interface{}) interface{} {
	src := GetStructElem(props)

	// Create pointer cache struct
	dstPtr := reflect.New(i.cachingType)
	dstFields := reflect.VisibleFields(i.cachingType)

	// Get the underlying struct
	dst := dstPtr.Elem()

	// Iterate fields of destination
	for _, dstField := range dstFields {
		setFieldVal(src, dst, dstField)
	}

	// if debug {
	// 	log.Printf("makeCachingProperties - src: %s: %v", src.Type(), src)
	// 	log.Printf("makeCachingProperties - dst: %s: %v", dst.Type(), dst)
	// }

	return dstPtr.Interface()
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
		log.Printf("cache %s: write %s", i, id)
	}

	if err := f(ctx, id, cachingProperties); err != nil {
		return ErrCache{err, fmt.Sprintf("cache error in writing %+v to %s: %s", cachingProperties, id, err.Error())}
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
		log.Printf("cache %s: delete %s", i, id)
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
		log.Printf("cache %s: hit %s", i.cachingIndex, id)
		// }

		return found, err
	}

	if debug {
		log.Printf("cache %s: miss %s", i.cachingIndex, id)
	}

	var backingErr error
	if found, backingErr = i.backingIndex.Get(ctx, id, dst, fields...); backingErr != nil {
		// Backing errors overwrite cache errors.
		err = backingErr
	}

	if found {
		if debug {
			log.Printf("backing %s: hit %s", i.backingIndex, id)
		}

		if indexErr := i.cacheWrite(ctx, id, dst, i.cachingIndex.Index); indexErr != nil {
			err = indexErr
		}
	}

	if debug {
		log.Printf("backing %s: miss %s", i.backingIndex, id)
	}

	return found, err
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Index{}
