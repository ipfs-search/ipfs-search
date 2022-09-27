package opensearch

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jpillora/backoff"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	opensearchtransport "github.com/opensearch-project/opensearch-go/v2/opensearchtransport"
	opensearchutil "github.com/opensearch-project/opensearch-go/v2/opensearchutil"
	"go.opentelemetry.io/otel/trace"

	"github.com/ipfs-search/ipfs-search/components/index/opensearch/bulkgetter"
	"github.com/ipfs-search/ipfs-search/instr"
)

// Client for search index.
type Client struct {
	searchClient *opensearch.Client
	bulkIndexer  opensearchutil.BulkIndexer
	bulkGetter   bulkgetter.AsyncGetter

	*instr.Instrumentation
}

// ClientConfig configures search index.
type ClientConfig struct {
	URL       string
	Transport http.RoundTripper
	Debug     bool

	BulkIndexerWorkers    int
	BulkIndexerFlushBytes int

	BulkGetterBatchSize    int
	BulkGetterBatchTimeout time.Duration
}

// NewClient returns a configured search index, or an error.
func NewClient(cfg *ClientConfig, i *instr.Instrumentation) (*Client, error) {
	var (
		c   *opensearch.Client
		bi  opensearchutil.BulkIndexer
		bg  bulkgetter.AsyncGetter
		err error
	)

	if cfg == nil {
		panic("NewClient ClientConfig cannot be nil.")
	}

	if i == nil {
		panic("NewCLient Instrumentation cannot be nil.")
	}

	if c, err = getSearchClient(cfg, i); err != nil {
		return nil, err
	}

	if bi, err = getBulkIndexer(c, cfg, i); err != nil {
		return nil, err
	}

	if bg, err = getBulkGetter(c, cfg, i); err != nil {
		return nil, err
	}

	return &Client{
		searchClient:    c,
		bulkIndexer:     bi,
		bulkGetter:      bg,
		Instrumentation: i,
	}, nil
}

// Work starts (and closes) a client worker.
func (c *Client) Work(ctx context.Context) error {
	// Flush indexing buffers on context close.
	// Use background context because current context is already closed.
	defer c.bulkIndexer.Close(context.Background())

	return c.bulkGetter.Work(ctx)
}

func getSearchClient(cfg *ClientConfig, i *instr.Instrumentation) (*opensearch.Client, error) {
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
		// https://www.elastic.co/guide/en/opensearch/reference/master/tune-for-indexing-speed.html#multiple-workers-threads
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
		NumWorkers: cfg.BulkIndexerWorkers,
		FlushBytes: cfg.BulkIndexerFlushBytes,
		OnFlushStart: func(ctx context.Context) context.Context {
			newCtx, _ := i.Tracer.Start(ctx, "index.opensearch.BulkIndexerFlush")
			return newCtx
		},
		OnError: func(ctx context.Context, err error) {
			span := trace.SpanFromContext(ctx)
			span.RecordError(err)
			log.Printf("Error flushing index buffer: %s", err)
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

func getBulkGetter(client *opensearch.Client, cfg *ClientConfig, i *instr.Instrumentation) (bulkgetter.AsyncGetter, error) {
	bgCfg := bulkgetter.Config{
		Client:       client,
		BatchSize:    cfg.BulkGetterBatchSize,
		BatchTimeout: cfg.BulkGetterBatchTimeout,
	}

	return bulkgetter.New(bgCfg), nil
}
