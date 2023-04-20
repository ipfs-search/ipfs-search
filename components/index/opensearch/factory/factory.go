package factory

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchutil"
)

var (
	errAliasNotFound  = errors.New("alias not found")
	ErrMappingInvalid = errors.New("mapping invalid")
)

type Factory struct {
	client *opensearch.Client
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

func (f *Factory) resolveAlias(aliasName string) (string, error) {
	resp, err := f.client.API.Indices.GetAlias(
		f.client.API.Indices.GetAlias.WithName(aliasName),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 404:
		return "", errAliasNotFound
	case 200:
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return "", err
		}

		if len(result) != 1 {
			panic(fmt.Sprintf("unexpected result resolving alias '%s': %s", aliasName, result))
		}

		for indexName := range result {
			return indexName, nil
		}
	}

	panic("unexpected status code returned resolving alias")
}

func (f *Factory) getCurrentMapping(indexName string) (interface{}, error) {
	resp, err := f.client.API.Indices.GetMapping(
		f.client.API.Indices.GetMapping.WithIndex(indexName),
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
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 0,
			},
		},
		"mappings": mapping,
		"aliases": map[string]interface{}{
			aliasName: nil,
		},
	}

	resp, err := f.client.API.Indices.Create(
		indexName,
		f.client.API.Indices.Create.WithBody(opensearchutil.NewJSONReader(body)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unable to create index: %s", indexName)
	}

	return nil
}

func New(client *opensearch.Client) *Factory {
	return &Factory{
		client,
	}
}

func (f *Factory) EnsureMapping(aliasName string) error {
	indexName, err := f.resolveAlias(aliasName)
	if err == errAliasNotFound {
		indexName = aliasName + "_v1"

		if err := f.createIndex(aliasName, indexName); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return f.validateMapping(aliasName, indexName)
}
