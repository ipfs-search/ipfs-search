package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
	opensearchutil "github.com/opensearch-project/opensearch-go/opensearchutil"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/instr"
)

// Index wraps an Elasticsearch index to store documents
type Index struct {
	es  *opensearch.Client
	cfg *Config

	*instr.Instrumentation
}

// New returns a new index.
func New(es *opensearch.Client, cfg *Config, i *instr.Instrumentation) index.Index {
	return &Index{
		es:              es,
		cfg:             cfg,
		Instrumentation: i,
	}
}

// String returns the name of the index, for convenient logging.
func (i *Index) String() string {
	return i.cfg.Name
}

// Index a document's properties, identified by id
func (i *Index) Index(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.Tracer.Start(ctx, "index.elasticsearch.Index")
	defer span.End()

	req := opensearchapi.IndexRequest{
		Index:      i.cfg.Name,
		Body:       opensearchutil.NewJSONReader(properties),
		DocumentID: id,
	}

	res, err := req.Do(ctx, i.es)
	defer res.Body.Close()

	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
	}

	return err
}

// Update a document's properties, given id
func (i *Index) Update(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.Tracer.Start(ctx, "index.elasticsearch.Update")
	defer span.End()

	req := opensearchapi.UpdateRequest{
		Index:      i.cfg.Name,
		Body:       opensearchutil.NewJSONReader(properties),
		DocumentID: id,
	}

	res, err := req.Do(ctx, i.es)
	defer res.Body.Close()

	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
	}

	return err
}

// Get retreives `fields` from document with `id` from the index, returning:
// - (true, decoding_error) if found (decoding error set when errors in json)
// - (false, nil) when not found
// - (false, error) otherwise
func (i *Index) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	ctx, span := i.Tracer.Start(ctx, "index.elasticsearch.Get")
	defer span.End()

	req := opensearchapi.GetRequest{
		Index:          i.cfg.Name,
		DocumentID:     id,
		SourceIncludes: fields,
		Realtime:       &[]bool{true}[0],
		Preference:     "_local",
	}

	res, err := req.Do(ctx, i.es)

	// Handle connection errors
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return false, err
	}

	defer res.Body.Close()

	// Decode body
	response := struct {
		Found  bool            `json:"found"`
		Source json.RawMessage `json:"_source"`
	}{}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&response)
	if err != nil {
		err = fmt.Errorf("error decoding body: %w", err)
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return false, err
	}

	if response.Found {
		// Decode source
		err = json.Unmarshal(response.Source, dst)
		if err != nil {
			err = fmt.Errorf("error decoding source: %w", err)
			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
			return false, err
		}
	}

	if res.StatusCode != 404 {
		// 404's do not signify an error, other status codes do.
		err = fmt.Errorf("unexpected status from backend: %s", res.Status())
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
	}

	return false, err
}

// Delete item from index
func (i *Index) Delete(ctx context.Context, id string) error {
	ctx, span := i.Tracer.Start(ctx, "index.elasticsearch.Delete")
	defer span.End()

	req := opensearchapi.DeleteRequest{
		Index:      i.cfg.Name,
		DocumentID: id,
	}

	res, err := req.Do(ctx, i.es)
	defer res.Body.Close()

	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
	}

	return err
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Index{}
