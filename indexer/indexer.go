package indexer

import (
	"encoding/json"
	"gopkg.in/olivere/elastic.v3"
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
		Refresh(false).
		Do()

	if err != nil {
		// Handle error
		return err
	}

	return nil
}

// Return existing references for an object, or nil, and the type
func (i Indexer) GetReferences(hash string) ([]Reference, string, error) {
	fsc := elastic.NewFetchSourceContext(true)
	fsc.Include("references")

	res, err := i.el.Get().
		Index("ipfs").Type("_all").FetchSourceContext(fsc).Id(hash).Do()

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
