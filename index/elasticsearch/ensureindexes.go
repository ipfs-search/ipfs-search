package elasticsearch

import (
	"context"
	"github.com/ipfs-search/ipfs-search/index"
	"github.com/olivere/elastic/v6"
	"log"
)

func EnsureIndexes(ctx context.Context, esURL string, configs map[string]*index.Config) (indexes map[string]index.Index, err error) {
	es, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(esURL))
	if err != nil {
		return
	}

	log.Printf("Connected to ElasticSearch")

	indexes = make(map[string]index.Index, len(configs))

	for n, c := range configs {
		i := &Index{
			Client: es,
			Config: c,
		}

		if err = index.EnsureUpdated(ctx, i, c); err != nil {
			return
		}

		indexes[n] = i
	}

	return
}
