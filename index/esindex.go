package index

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v6"
)

// ESIndex wraps an Elasticsearch index to store documents
type ESIndex struct {
	Client *elastic.Client
	Name   string
}

// Index a document's properties, identified by id
func (i *ESIndex) Index(ctx context.Context, id string, properties map[string]interface{}) error {
	_, err := i.Client.Index().
		Index(i.Name).
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
func (i *ESIndex) Update(ctx context.Context, id string, properties map[string]interface{}) error {
	_, err := i.Client.Update().
		Index(i.Name).
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
func (i *ESIndex) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	fsc := elastic.NewFetchSourceContext(true)
	fsc.Include(fields...)

	result, err := i.Client.
		Get().
		Index(i.Name).
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
