package pool

import (
	"context"
	"log"
	"time"

	"github.com/ipfs-search/ipfs-search/components/crawler"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch"
	"github.com/ipfs-search/ipfs-search/utils"
)

func startOpenSearchWorker(ctx context.Context, esClient *opensearch.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := esClient.Work(ctx); err != nil {
				log.Printf("Error in ES client pool, restarting pool: %s", err)
				// Prevent overly tight restart loop
				time.Sleep(time.Second)
			}
		}
	}
}

func (w *Pool) getOpenSearchClient() (*opensearch.Client, error) {
	clientConfig := &opensearch.ClientConfig{
		URL:       w.config.OpenSearch.URL,
		Transport: utils.GetHTTPTransport(w.dialer.DialContext, 100),
		Debug:     false,

		BulkIndexerWorkers:     w.config.OpenSearch.BulkIndexerWorkers,
		BulkIndexerFlushBytes:  int(w.config.OpenSearch.BulkIndexerFlushBytes),
		BulkGetterBatchSize:    w.config.OpenSearch.BulkGetterBatchSize,
		BulkGetterBatchTimeout: w.config.OpenSearch.BulkGetterBatchTimeout,
	}

	return opensearch.NewClient(clientConfig, w.Instrumentation)
}

func (w *Pool) getIndexes(ctx context.Context) (*crawler.Indexes, error) {
	esClient, err := w.getOpenSearchClient()
	if err != nil {
		return nil, err
	}

	// Start ES pools
	go startOpenSearchWorker(ctx, esClient)

	return &crawler.Indexes{
		Files: opensearch.New(
			esClient,
			&opensearch.Config{Name: w.config.Indexes.Files.Name},
		),
		Directories: opensearch.New(
			esClient,
			&opensearch.Config{Name: w.config.Indexes.Directories.Name},
		),
		Invalids: opensearch.New(
			esClient,
			&opensearch.Config{Name: w.config.Indexes.Invalids.Name},
		),
		Partials: opensearch.New(
			esClient,
			&opensearch.Config{Name: w.config.Indexes.Partials.Name},
		),
	}, nil
}
