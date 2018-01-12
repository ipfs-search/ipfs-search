package indexer

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

// Indexer performs indexing of items and its references using ElasticCloud
type Indexer struct {
	ElasticSearch *elastic.Client
}

// Reference to indexed item
type Reference struct {
	ParentHash string `json:"parent_hash"`
	Name       string `json:"name"`
}

// IndexItem adds or updates an IPFS item with arbitrary properties
func (i *Indexer) IndexItem(doctype string, hash string, properties map[string]interface{}) error {
	_, err := i.ElasticSearch.Update().
		Index("ipfs").
		Type(doctype).
		Id(hash).
		Doc(properties).
		DocAsUpsert(true).
		Do(context.TODO())

	if err != nil {
		// Handle error
		return err
	}

	return nil
}

// extractRefrences reads the refernces from the JSON response from ElasticSearch
func extractReferences(result *elastic.GetResult) ([]Reference, error) {
	var parsedResult map[string][]Reference

	err := json.Unmarshal(*result.Source, parsedResult)
	if err != nil {
		return nil, err
	}

	references, ok := parsedResult["references"]
	if !ok {
		return nil, errors.New("references not in output")
	}

	return references, nil
}

// GetReferences returns existing references and the type for an object, or nil.
// When no object is found nil is returned but no error is set.
// If no object is found, an empty list is returned.
func (i *Indexer) GetReferences(hash string) ([]Reference, string, error) {
	fsc := elastic.NewFetchSourceContext(true)
	fsc.Include("references")

	result, err := i.ElasticSearch.
		Get().
		Index("ipfs").Type("_all").
		FetchSourceContext(fsc).
		Id(hash).
		Do(context.TODO())

	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, "", nil
		}
		return nil, "", err
	}

	references, err := extractReferences(result)
	if err != nil {
		return nil, "", err
	}

	return references, result.Type, nil
}
