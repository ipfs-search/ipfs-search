package indexer

import (
	"context"
	"encoding/json"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

// Indexer performs indexing of items and its references using ElasticCloud
type Indexer struct {
	ElasticSearch *elastic.Client
}

// IndexItem adds or updates an IPFS item with arbitrary properties
func (i *Indexer) IndexItem(ctx context.Context, doctype string, hash string, properties map[string]interface{}) error {
	_, err := i.ElasticSearch.Update().
		Index("ipfs").
		Type(doctype).
		Id(hash).
		Doc(properties).
		DocAsUpsert(true).
		Do(ctx)

	if err != nil {
		// Handle error
		return err
	}

	return nil
}

// extractRefrences reads the refernces from the JSON response from ElasticSearch
func extractReferences(result *elastic.GetResult) (References, error) {
	var parsedResult map[string]References

	err := json.Unmarshal(*result.Source, &parsedResult)
	if err != nil {
		log.Printf("can't unmarshal references JSON: %s", *result.Source)
		return nil, err
	}

	references := parsedResult["references"]

	return references, nil
}

// GetReferences returns existing references and the type for an object.
// When no object is found an empty list of references is returned, the
// type is "" and no error is set.
func (i *Indexer) GetReferences(ctx context.Context, hash string) (References, string, error) {
	fsc := elastic.NewFetchSourceContext(true)
	fsc.Include("references")

	result, err := i.ElasticSearch.
		Get().
		Index("ipfs").Type("_all").
		FetchSourceContext(fsc).
		Id(hash).
		Do(ctx)

	if err != nil {
		if elastic.IsNotFound(err) {
			// Initialize empty references when none have been found
			return []Reference{}, "", nil
		}
		return nil, "", err
	}

	references, err := extractReferences(result)
	if err != nil {
		return nil, "", err
	}

	return references, result.Type, nil
}
