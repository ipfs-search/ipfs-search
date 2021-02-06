package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/instr"
)

// Index wraps an Elasticsearch index to store documents
type Index struct {
	es  *elastic.Client
	cfg *Config

	*instr.Instrumentation
}

// New returns a new index.
func New(es *elastic.Client, cfg *Config, i *instr.Instrumentation) index.Index {
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

	_, err := i.es.Index().
		Index(i.cfg.Name).
		Id(id).
		BodyJson(properties).
		Do(ctx)

	if err != nil {
		// Handle error
		return err
	}

	return nil

}

// Update a document's properties, given id
func (i *Index) Update(ctx context.Context, id string, properties interface{}) error {
	ctx, span := i.Tracer.Start(ctx, "index.elasticsearch.Update")
	defer span.End()

	_, err := i.es.Update().
		Index(i.cfg.Name).
		Id(id).
		Doc(properties).
		Do(ctx)

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

	fsc := elastic.NewFetchSourceContext(true)
	fsc.Include(fields...)

	result, err := i.es.
		Get().
		Index(i.cfg.Name).
		FetchSourceContext(fsc).
		Id(id).
		Do(ctx)

	switch {
	case err == nil:
		// Found

		// Decode resulting field json into `dst`
		err = json.Unmarshal(result.Source, dst)

		if err != nil {
			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		}

		return true, err
	case elastic.IsNotFound(err):
		// 404
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Ok))

		return false, nil

	default:
		// Unknown error, propagate
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return false, err
	}
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Index{}
