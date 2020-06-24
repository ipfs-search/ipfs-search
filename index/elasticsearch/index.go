package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/ipfs-search/ipfs-search/index"
	"github.com/olivere/elastic/v6"
)

// Index wraps an Elasticsearch index to store documents
type Index struct {
	es  *elastic.Client
	cfg *index.Config
}

// New returns a new index.
func New(es *elastic.Client, cfg *index.Config) *Index {
	return &Index{
		es:  es,
		cfg: cfg,
	}
}

// GetConfig returns the config for the index.
func (i *Index) GetConfig() *index.Config {
	return i.cfg
}

// String returns the name of the index, for convenient logging.
func (i *Index) String() string {
	return i.cfg.Name
}

// Index a document's properties, identified by id
func (i *Index) Index(ctx context.Context, id string, properties map[string]interface{}) error {
	_, err := i.es.Index().
		Index(i.cfg.Name).
		Type("_doc").
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
func (i *Index) Update(ctx context.Context, id string, properties map[string]interface{}) error {
	_, err := i.es.Update().
		Index(i.cfg.Name).
		Type("_doc").
		Id(id).
		Doc(properties).
		Do(ctx)

	if err != nil {
		// Handle error
		return err
	}

	return nil

}

// Get retreives `fields` from document with `id` from the index, returning:
// - (true, decoding_error) if found (decoding error set when errors in json)
// - (false, nil) when not found
// - (false, error) otherwise
func (i *Index) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	fsc := elastic.NewFetchSourceContext(true)
	fsc.Include(fields...)

	result, err := i.es.
		Get().
		Index(i.cfg.Name).
		Type("_doc").
		FetchSourceContext(fsc).
		Id(id).
		Do(ctx)

	switch {
	case err == nil:
		// Found

		// Decode resulting field json into `dst`
		err = json.Unmarshal(*result.Source, dst)

		return true, err
	case elastic.IsNotFound(err):
		// 404
		return false, nil

	default:
		// Unknown error, propagate
		return false, err
	}
}
