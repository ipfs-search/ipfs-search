package factory

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch/aliasresolver"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch/client"
	os_index "github.com/ipfs-search/ipfs-search/components/index/opensearch/index"
	"github.com/opensearch-project/opensearch-go/v2/opensearchutil"
)

var (
	ErrMappingInvalid = errors.New("mapping invalid")
)

type Factory struct {
	client *client.Client
}

//go:embed mappings/*
var mappings embed.FS

func (f *Factory) getDesiredMapping(aliasName string) (interface{}, error) {
	var mapping interface{}

	fileName := "mappings/" + aliasName + ".json"

	file, err := mappings.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&mapping)
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func (f *Factory) getCurrentMapping(indexName string) (interface{}, error) {
	resp, err := f.client.SearchClient.API.Indices.GetMapping(
		f.client.SearchClient.API.Indices.GetMapping.WithIndex(indexName),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Unexpected response on index mapping query for index '" + indexName + "'.")
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	mappings, ok := result[indexName].(map[string]interface{})["mappings"]
	if !ok {
		panic("Unable to parse mappings from response.")
	}

	return mappings, nil
}

func (f *Factory) validateMapping(aliasName, indexName string) error {
	desired, err := f.getDesiredMapping(aliasName)
	if err != nil {
		return err
	}

	current, err := f.getCurrentMapping(indexName)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(desired, current) {
		return fmt.Errorf("%w: %s", ErrMappingInvalid, aliasName)
	}

	return nil
}

func (f *Factory) createIndex(aliasName, indexName string) error {
	mapping, err := f.getDesiredMapping(aliasName)
	if err != nil {
		return err
	}

	body := map[string]interface{}{
		// For developing and testing, these are sensible defaults.
		// For production, one will want to manually create indexes.
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 0,
			},
		},
		"mappings": mapping,
		"aliases": map[string]interface{}{
			aliasName: aliasName,
		},
	}

	resp, err := f.client.SearchClient.API.Indices.Create(
		indexName,
		f.client.SearchClient.API.Indices.Create.WithBody(opensearchutil.NewJSONReader(body)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unable to create index: %s", indexName)
	}

	// Validate resulting mapping.
	return f.validateMapping(aliasName, indexName)
}

func New(client *client.Client) *Factory {
	return &Factory{
		client,
	}
}

func (f *Factory) generateIndexName(aliasName string) string {
	return aliasName + "_v1"
}

func (f *Factory) resolveOrCreateIndex(ctx context.Context, aliasName string) (string, error) {
	indexName, err := f.client.AliasResolver.GetIndex(ctx, aliasName)
	if errors.Is(err, aliasresolver.ErrNotFound) {
		// Alias not found, create index.
		indexName = f.generateIndexName(aliasName)
		return indexName, f.createIndex(aliasName, indexName)
	}

	return indexName, err
}

func (f *Factory) NewIndex(ctx context.Context, aliasName string) (index.Index, error) {
	indexName, err := f.resolveOrCreateIndex(ctx, aliasName)
	if err != nil {
		return nil, err
	}

	cfg := os_index.Config{
		Name: indexName,
	}

	return os_index.New(f.client, &cfg), nil
}
