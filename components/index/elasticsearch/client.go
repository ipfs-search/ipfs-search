package elasticsearch

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jpillora/backoff"
	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchtransport "github.com/opensearch-project/opensearch-go/opensearchtransport"
	opensearchutil "github.com/opensearch-project/opensearch-go/opensearchutil"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"

	"github.com/ipfs-search/ipfs-search/instr"
)

// Client for search index.
type Client struct {
	searchClient *opensearch.Client
	bulkIndexer  opensearchutil.BulkIndexer

	*instr.Instrumentation
}

// ClientConfig configures search index.
type ClientConfig struct {
	URL             string
	Transport       http.RoundTripper
	Debug           bool
	IndexBufferSize int
}

// NewClient returns a configured search index, or an error.
func NewClient(cfg *ClientConfig, i *instr.Instrumentation) (*Client, error) {
	if cfg == nil {
		panic("NewClient ClientConfig cannot be nil.")
	}

	if i == nil {
		panic("NewCLient Instrumentation cannot be nil.")
	}

	c, err := getSearchClient(cfg, i)
	if err != nil {
		return nil, err
	}

	bi, err := getBulkIndexer(c, cfg, i)
	if err != nil {
		return nil, err
	}

	return &Client{
		searchClient:    c,
		bulkIndexer:     bi,
		Instrumentation: i,
	}, nil
}

// Close client connection and flush bulk indexer.
func (c *Client) Close(ctx context.Context) error {
	return c.bulkIndexer.Close(ctx)
}

func getSearchClient(cfg *ClientConfig, i *instr.Instrumentation) (*opensearch.Client, error) {

	// TODO: Re-enable
	// httpClient := utils.GetHTTPClient(w.dialer.DialContext, 5)

	b := backoff.Backoff{
		Factor: 2.0,
		Jitter: true,
	}

	// Ref: https://pkg.go.dev/github.com/opensearch-project/opensearch-go@v1.0.0#Config
	clientConfig := opensearch.Config{
		Addresses:    []string{cfg.URL},
		Transport:    cfg.Transport,
		DisableRetry: cfg.Debug,
		// Retry/backoff management
		// https://www.elastic.co/guide/en/elasticsearch/reference/master/tune-for-indexing-speed.html#multiple-workers-threads
		RetryOnStatus:        []int{429, 502, 503, 504},
		EnableRetryOnTimeout: true,
		RetryBackoff:         func(i int) time.Duration { return b.ForAttempt(float64(i)) },
	}

	if cfg.Debug {
		clientConfig.Logger = &opensearchtransport.TextLogger{
			Output:             log.Default().Writer(),
			EnableRequestBody:  cfg.Debug,
			EnableResponseBody: cfg.Debug,
		}
	}

	return opensearch.NewClient(clientConfig)
}

func getBulkIndexer(client *opensearch.Client, cfg *ClientConfig, i *instr.Instrumentation) (opensearchutil.BulkIndexer, error) {
	iCfg := opensearchutil.BulkIndexerConfig{
		Client:     client,
		NumWorkers: 1, // Start conservatively with 1 worker.
		FlushBytes: cfg.IndexBufferSize,
		OnFlushStart: func(ctx context.Context) context.Context {
			newCtx, _ := i.Tracer.Start(ctx, "index.elasticsearch.BulkIndexerFlush")
			return newCtx
		},
		OnError: func(ctx context.Context, err error) {
			span := trace.SpanFromContext(ctx)
			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
			log.Printf("Error flushing index buffer: %v", err)
		},
		OnFlushEnd: func(ctx context.Context) {
			span := trace.SpanFromContext(ctx)
			log.Println("Flushed index buffer")

			// log.Printf("ES stats: %+v", )
			span.End()
		},
	}

	if cfg.Debug {
		iCfg.FlushBytes = 1
		iCfg.FlushInterval = 0
	}

	return opensearchutil.NewBulkIndexer(iCfg)
}
