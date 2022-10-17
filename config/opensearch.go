package config

import (
	"runtime"
	"time"

	"github.com/c2h5oh/datasize"
)

// OpenSearch holds configuration for OpenSearch.
type OpenSearch struct {
	URL                    string            `yaml:"url" env:"OPENSEARCH_URL"`
	BulkIndexerWorkers     int               `yaml:"bulk_indexer_workers"`      // Amount of workers to user for indexer.
	BulkIndexerFlushBytes  datasize.ByteSize `yaml:"bulk_flush_bytes"`          // Flush index buffer after this many bytes.
	BulkIndexerFlushTimeout time.Duration     `yaml:"bulk_flush_timeout"`        // Flush index buffer after this much time.
	BulkGetterBatchSize    int               `yaml:"bulk_getter_batch_size"`    // Maximum batch size for bulk gets.
	BulkGetterBatchTimeout time.Duration     `yaml:"bulk_getter_batch_timeout"` // Maximum time to wait until executing batch.
}

// OpenSearchDefaults returns the defaults for OpenSearch.
func OpenSearchDefaults() OpenSearch {
	return OpenSearch{
		URL:                    "http://localhost:9200",
		BulkIndexerWorkers:     runtime.NumCPU(),
		BulkIndexerFlushTimeout: 5 * time.Minute,
		BulkIndexerFlushBytes:  5e+6, // 5MB
		BulkGetterBatchSize:    48,
		BulkGetterBatchTimeout: 150 * time.Millisecond,
	}
}
