package indexer

import (
	"encoding/json"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

type Indexer struct {
	el *elastic.Client
}

type Reference struct {
	ParentHash string `json:"parent_hash"`
	Name       string `json:"name"`
}

func NewIndexer(el *elastic.Client) *Indexer {
	return &Indexer{
		el: el,
	}
}

// Add file or directory to index
func (i Indexer) IndexItem(doctype string, hash string, properties map[string]interface{}) error {
	_, err := i.el.Update().
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

// Return existing references for an object, or nil, and the type.
// When no object is found nil is returned but no error is set.
// Otherwise, an empty list is returned.
func (i Indexer) GetReferences(hash string) ([]Reference, string, error) {
	fsc := elastic.NewFetchSourceContext(true)
	fsc.Include("references")

	res, err := i.el.Get().
		Index("ipfs").Type("_all").FetchSourceContext(fsc).Id(hash).Do(context.TODO())

	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, "", nil
		}
		return nil, "", err
	}

	var result map[string][]Reference
	err = json.Unmarshal(*res.Source, &result)
	if err != nil {
		return nil, "", err
	}

	return result["references"], res.Type, nil
}
