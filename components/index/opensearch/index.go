package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	opensearchutil "github.com/opensearch-project/opensearch-go/v2/opensearchutil"
	"go.opentelemetry.io/otel/codes"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch/bulkgetter"
)

const debug bool = true

// Index wraps an OpenSearch index to store documents
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

// String returns the name of the index, for convenient logging.
func (i *Index) String() string {
	return i.cfg.Name
}

func getBody(v interface{}) (io.ReadSeeker, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

// index wraps BulkIndexer.Add().
func (i *Index) index(
	ctx context.Context,
	action string,
	id string,
	properties interface{},
) error {
	ctx, span := i.c.Tracer.Start(ctx, "index.opensearch.index")
	defer span.End()

	var (
		body io.ReadSeeker
		err  error
	)

	if properties != nil {
		if action == "update" {
			// For updates, the updated fields need to be wrapped in a `doc` field
			body, err = getBody(struct {
				Doc interface{} `json:"doc"`
			}{properties})
		} else {
			body, err = getBody(properties)
		}
		if err != nil {
			panic(err)
		}
	}

	item := opensearchutil.BulkIndexerItem{
		Index:      i.cfg.Name,
		Action:     action,
		Body:       body,
		DocumentID: id,
		Version:    nil,
		OnFailure: func(
			ctx context.Context,
			item opensearchutil.BulkIndexerItem,
			res opensearchutil.BulkIndexerResponseItem, err error,
		) {
			if err == nil {
				err = fmt.Errorf("Error flushing: %+v (%s)", res, id)
			}

			span.RecordError(err)
			log.Println(err)

		},
	}

	ctx, span = i.c.Tracer.Start(ctx, "index.opensearch.bulkIndexer.Add")
	defer span.End()

	err = i.c.bulkIndexer.Add(ctx, item)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error adding to BulkIndexer.")
	}

	return err
}

// Index a document's properties, identified by id
func (i *Index) Index(ctx context.Context, id string, properties interface{}) error {
	return i.index(ctx, "create", id, properties)
}

// Update a document's properties, given id
func (i *Index) Update(ctx context.Context, id string, properties interface{}) error {
	return i.index(ctx, "update", id, properties)
}

// Delete item from index
func (i *Index) Delete(ctx context.Context, id string) error {
	return i.index(ctx, "delete", id, nil)
}

// Get retreives `fields` from document with `id` from the index, returning:
// - (true, decoding_error) if found (decoding error set when errors in json)
// - (false, nil) when not found
// - (false, error) otherwise
func (i *Index) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	ctx, span := i.c.Tracer.Start(ctx, "index.opensearch.Get")
	defer span.End()

	req := bulkgetter.GetRequest{
		Index:      i.cfg.Name,
		DocumentID: id,
		Fields:     fields,
	}

	resp := <-i.c.bulkGetter.Get(ctx, &req, dst)

	if debug {
		if resp.Found {
			log.Printf("opensearch: found %s in %s", id, i)
		} else {
			if resp.Error != nil {
				log.Printf("opensearch: error getting %s in %s: %v", id, i, resp.Error)
			}
		}
	}

	return resp.Found, resp.Error
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Index{}
