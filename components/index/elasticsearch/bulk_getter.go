package elasticsearch

import (
	"context"

	"github.com/opensearch-project/opensearch-go"
	// "github.com/opensearch-project/opensearch-go/opensearchapi"
)

// BatchingGetter wraps a GET API call and an optimized batching implementation.
// It's inspired by BulkIndexer from https://github.com/opensearch-project/opensearch-go/blob/v1.0.0/opensearchutil/bulk_indexer.go
type BatchingGetter interface {
	// Get retrieves specific fields of a document from an index by id, returning whether or not the document was found.
	Get(context.Context, GetRequest) (bool, error)
}

// BatchingGetterConfig contains the configuration for a BatchingGetter.
type BatchingGetterConfig struct {
	// NumWorkers int // The number of workers. Defaults to runtime.NumCPU().
	// FlushInterval time.Duration // The flush threshold as duration. Defaults to 1sec.

	Client *opensearch.Client // The OpenSearch client.
	// Decoder     BulkResponseJSONDecoder // A custom JSON decoder.
	// DebugLogger BulkIndexerDebugLogger  // An optional logger for debugging.

	// OnError      func(context.Context, error)          // Called for indexer errors.
	// OnFlushStart func(context.Context) context.Context // Called when the flush starts.
	// OnFlushEnd   func(context.Context)                 // Called when the flush ends.

	// // Parameters of the Bulk API.
	// Index               string
	// ErrorTrace          bool
	// FilterPath          []string
	// Header              http.Header
	// Human               bool
	// Pipeline            string
	// Pretty              bool
	// Refresh             string
	// Routing             string
	// Source              []string
	// SourceExcludes      []string
	// SourceIncludes      []string
	// Timeout             time.Duration
	// WaitForActiveShards string
}

// GetRequest represents an item to GET.
type GetRequest struct {
	Index       string
	DocumentID  string
	Fields      []string
	Destination interface{}
}

// BatchingGetterStats represents statistics for a bulk getter.
// type BatchingGetterStats struct {
// 	NumFound    uint64
// 	NumNotFound uint64
// 	NumError    uint64
// }

// NewBatchingGetter returns a new BatchingGetter, setting sensible defaults for the configuration.
func NewBatchingGetter(cfg BatchingGetterConfig) BatchingGetter {
	if cfg.Client == nil {
		cfg.Client, _ = opensearch.NewDefaultClient()
	}

	// if cfg.Decoder == nil {
	// 	cfg.Decoder = defaultJSONDecoder{}
	// }

	// if cfg.NumWorkers == 0 {
	// 	cfg.NumWorkers = runtime.NumCPU()
	// }

	// if cfg.FlushInterval == 0 {
	// 	cfg.FlushInterval = 1 * time.Second
	// }

	bg := batchingGetter{
		config: cfg,
		done:   make(chan bool),
		// stats:  &batchingGetterStats{},
	}

	bg.init()

	return &bg
}

type batchingGetter struct {
}

func (bg *batchingGetter) init() {
	// Initialize a batchingGetter.
}
