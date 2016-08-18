package indexer

import (
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

// Add parent references to indexed item
func (i Indexer) IndexReference(doctype string, hash string, name string, parent_hash string) error {
	/*
		'<hash>': {
			'references': {
				'<parent_hash>': {
					'name': '<name>'
				}
			}
		}

		if (document_exists) {
			if (references_exists) {
				add_parent_hash to references
			} else {
				add references to document
			}
		} else {
			create document with references as only information
		}
	*/
	reference := map[string]interface{}{
		parent_hash: map[string]interface{}{
			"name": name,
		},
	}

	properties := map[string]interface{}{
		"references": reference,
	}

	// TODO: Use smart scripted update to allow for multiple files per object
	_, err := i.el.Update().
		Index("ipfs").
		Type(doctype).
		Id(hash).
		Doc(properties).
		DocAsUpsert(true).
		// Script(NewScript("if (ctx._source.references) {ctx._source.references += reference } else { references = [reference, ] }").
		// Params(reference).
		// Upsert(properties).
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
