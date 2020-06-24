package elasticsearch

import (
	"context"
	"github.com/ipfs-search/ipfs-search/index"
	"github.com/olivere/elastic/v6"
	"log"
)

func getElastic(url string) (*elastic.Client, error) {
	el, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(url))
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to ElasticSearch")

	return el, nil
}

func getIndex(ctx context.Context, es *elastic.Client, config *index.Config) (index.Index, error) {
	if config == nil {
		panic("configuration for index nil")
	}

	i := &Index{
		Client: es,
		Name:   config.Name,
	}

	// Create index if it doesn't already exists, update if it is different (last parameter, true).
	if err := index.EnsureUpdated(ctx, i, config); err != nil {
		return nil, err
	}

	return i, nil
}

func EnsureIndexes(ctx context.Context, esURL string, configs map[string]*index.Config) (indexes map[string]index.Index, err error) {
	es, err := getElastic(esURL)

	if err != nil {
		return
	}

	indexes = make(map[string]index.Index, len(configs))

	for n, c := range configs {
		indexes[n], err = getIndex(ctx, es, c)
		if err != nil {
			return
		}
	}

	return
}
