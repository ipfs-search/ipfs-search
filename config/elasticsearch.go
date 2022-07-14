package config

import (
	"runtime"
	"time"

	"github.com/c2h5oh/datasize"
)

// ElasticSearch holds configuration for ElasticSearch.
type ElasticSearch struct {
	URL                    string            `yaml:"url" env:"ELASTICSEARCH_URL"`
	BulkIndexerWorkers     int               `yaml:"bulk_indexer_workers"`      // Amount of workers to user for indexer.
	BulkIndexerFlushBytes  datasize.ByteSize `yaml:"bulk_flush_bytes"`          // Flush index buffer after this many bytes.
	BulkGetterBatchSize    int               `yaml:"bulk_getter_batch_size"`    // Maximum batch size for bulk gets.
	BulkGetterBatchTimeout time.Duration     `yaml:"bulk_getter_batch_timeout"` // Maximum time to wait until executing batch.
}

// ElasticSearchDefaults returns the defaults for ElasticSearch.
func ElasticSearchDefaults() ElasticSearch {
	return ElasticSearch{
		URL:                    "http://localhost:9200",
		BulkIndexerWorkers:     runtime.NumCPU(),
		BulkIndexerFlushBytes:  5e+6, // 5MB
		BulkGetterBatchSize:    48,
		BulkGetterBatchTimeout: 150 * time.Millisecond,
	}
}
