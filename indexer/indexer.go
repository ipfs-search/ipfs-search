package indexer

import (
	"gopkg.in/olivere/elastic.v3"
)

type Indexer struct {
	el *elastic.Client
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
		Refresh(true).
		Do()

	if err != nil {
		// Handle error
		return err
	}

	return nil
}

// Whether or not an object exists in index
func (i Indexer) IsIndexed(hash string) (bool, error) {
	return i.el.Exists().
		Index("ipfs").Type("_all").Id(hash).Do()
}
