package crawlworker

import (
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func getElastic(url string) (*elastic.Client, error) {
	el, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(url))
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
	log.Printf("Connected to ElasticSearch")

	return el, nil
}
