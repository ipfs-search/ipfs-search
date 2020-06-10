package factory

import (
	"context"
	"github.com/ipfs-search/ipfs-search/index"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"
	"github.com/olivere/elastic/v6"
	"log"
	// "os"
)

func getElastic(url string) (*elastic.Client, error) {
	// logger := log.New(os.Stderr, "es", log.LstdFlags)
	// elastic.SetTraceLog(logger)
	el, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(url))
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to ElasticSearch")

	return el, nil
}

func getIndex(ctx context.Context, el *elastic.Client, config *IndexConfig) (index.Index, error) {
	if config == nil {
		panic("configuration for index nil")
	}

	i := &elasticsearch.Index{
		Client: el,
		Name:   config.Name,
	}

	// Create index if it doesn't already exists, update if it is different (last parameter, true).
	if err := index.Assert(ctx, i, config.Settings, config.Mapping, true); err != nil {
		return nil, err
	}

	return i, nil
}
