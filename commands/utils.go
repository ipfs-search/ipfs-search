package commands

import (
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

func getElastic() (*elastic.Client, error) {
	el, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}
	exists, err := el.IndexExists("ipfs").Do(context.TODO())
	if err != nil {
		return nil, err
	}
	if !exists {
		// Index does not exist yet, create
		el.CreateIndex("ipfs")
	}

	return el, nil
}
