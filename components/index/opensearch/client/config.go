package client

import (
	"net/http"
	"time"
)

// Config configures search index.
type Config struct {
	URL       string
	Transport http.RoundTripper
	Debug     bool

	BulkIndexerWorkers      int
	BulkIndexerFlushBytes   int
	BulkIndexerFlushTimeout time.Duration

	BulkGetterBatchSize    int
	BulkGetterBatchTimeout time.Duration
}
