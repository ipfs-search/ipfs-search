package indexer

import (
	"gopkg.in/ipfs/go-ipfs-api.v1"
	"gopkg.in/olivere/elastic.v3"
)

type Indexer struct {
	el *elastic.Client
}

func NewIndexer(el *elastic.Client) *Indexer {
	i := new(Indexer)
	i.el = el
	return i
}

// Create directory index based on hash
func (i Indexer) IndexDirectory(list *shell.UnixLsObject) error {
	_, err := i.el.Index().
		Index("ipfs").
		Type("directory").
		Id(list.Hash).
		BodyJson(list.Links).
		Refresh(true).
		Do()
	if err != nil {
		// Handle error
		return err
	}

	return nil
}
