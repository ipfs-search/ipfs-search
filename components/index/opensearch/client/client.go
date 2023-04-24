package client

import (
	"context"
	"log"
	"time"

	"github.com/jpillora/backoff"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	opensearchtransport "github.com/opensearch-project/opensearch-go/v2/opensearchtransport"
	opensearchutil "github.com/opensearch-project/opensearch-go/v2/opensearchutil"
	"go.opentelemetry.io/otel/trace"

	"github.com/ipfs-search/ipfs-search/components/index/opensearch/aliasresolver"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch/bulkgetter"
	"github.com/ipfs-search/ipfs-search/instr"
)

// Client for search index.
type Client struct {
	SearchClient  *opensearch.Client
	AliasResolver aliasresolver.AliasResolver
	BulkIndexer   opensearchutil.BulkIndexer
	BulkGetter    bulkgetter.AsyncGetter

	*instr.Instrumentation
}

// New returns a configured search index, or an error.
func New(cfg *Config, i *instr.Instrumentation) (*Client, error) {
	if cfg == nil {
		panic("Config cannot be nil.")
	}

	if i == nil {
		panic("Instrumentation cannot be nil.")
	}

	c := Client{
		Instrumentation: i,
	}

	initFuncs := []func(*Config) error{
		c.setSearchClient,
		c.setAliasResolver,
		c.setBulkIndexer,
		c.setBulkGetter,
	}

	for _, initFunc := range initFuncs {
		if err := initFunc(cfg); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// Work starts (and closes) a client worker.
func (c *Client) Work(ctx context.Context) error {
	// Flush indexing buffers on context close.
	// Use background context because current context might already closed.
	defer c.BulkIndexer.Close(context.Background())

	return c.BulkGetter.Work(ctx)
}

func (c *Client) setSearchClient(cfg *Config) error {
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
		// Spread queries/load; discover nodes on start and do it again every 5 minutes.
		DiscoverNodesOnStart:  true,
		DiscoverNodesInterval: 5 * time.Minute,
	}

	if cfg.Debug {
		clientConfig.Logger = &opensearchtransport.TextLogger{
			Output:             log.Default().Writer(),
			EnableRequestBody:  cfg.Debug,
			EnableResponseBody: cfg.Debug,
		}
	}

	searchClient, err := opensearch.NewClient(clientConfig)
	if err != nil {
		return err
	}

	c.SearchClient = searchClient
	return nil
}

func (c *Client) setBulkIndexer(cfg *Config) error {
	if c.SearchClient == nil {
		panic("SearchClient is nil.")
	}

	iCfg := opensearchutil.BulkIndexerConfig{
		Client:        c.SearchClient,
		NumWorkers:    cfg.BulkIndexerWorkers,
		FlushBytes:    cfg.BulkIndexerFlushBytes,
		FlushInterval: cfg.BulkIndexerFlushTimeout,
		OnFlushStart: func(ctx context.Context) context.Context {
			newCtx, _ := c.Tracer.Start(ctx, "index.opensearch.BulkIndexerFlush")
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

	bulkIndexer, err := opensearchutil.NewBulkIndexer(iCfg)
	if err != nil {
		return err
	}

	c.BulkIndexer = bulkIndexer
	return nil
}

func (c *Client) setAliasResolver(cfg *Config) error {
	if c.SearchClient == nil {
		panic("SearchClient is nil.")
	}

	c.AliasResolver = aliasresolver.NewAliasResolver(c.SearchClient)
	return nil
}

func (c *Client) setBulkGetter(cfg *Config) error {
	if c.SearchClient == nil {
		panic("SearchClient is nil.")
	}

	if c.AliasResolver == nil {
		panic("AliasResolver is nil.")
	}

	bgCfg := bulkgetter.Config{
		Client:        c.SearchClient,
		BatchSize:     cfg.BulkGetterBatchSize,
		BatchTimeout:  cfg.BulkGetterBatchTimeout,
		AliasResolver: c.AliasResolver,
	}

	c.BulkGetter = bulkgetter.New(bgCfg)
	return nil
}
