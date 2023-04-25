package pool

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/ipfs-search/ipfs-search/components/crawler"
	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/components/index/cache"
	os_client "github.com/ipfs-search/ipfs-search/components/index/opensearch/client"
	os_factory "github.com/ipfs-search/ipfs-search/components/index/opensearch/factory"
	os_index "github.com/ipfs-search/ipfs-search/components/index/opensearch/index"
	"github.com/ipfs-search/ipfs-search/components/index/redis"
	indexTypes "github.com/ipfs-search/ipfs-search/components/index/types"
	"github.com/ipfs-search/ipfs-search/utils"
)

func osWorkLoop(ctx context.Context, workFunc func(context.Context) error) {
	for {
		// Keep starting worker unless context is done.
		select {
		case <-ctx.Done():
			return
		default:
			if err := workFunc(ctx); err != nil {
				log.Printf("Error in worker: %s, restarting.", err)
				// Prevent overly tight restart loop
				time.Sleep(time.Second)
			}
		}
	}
}

func (w *Pool) getOpenSearchClient() (*os_client.Client, error) {
	config := &os_client.Config{
		URL:       w.config.OpenSearch.URL,
		Transport: utils.GetHTTPTransport(w.dialer.DialContext, 100),
		Debug:     false,

		BulkIndexerWorkers:      w.config.OpenSearch.BulkIndexerWorkers,
		BulkIndexerFlushBytes:   int(w.config.OpenSearch.BulkIndexerFlushBytes),
		BulkIndexerFlushTimeout: w.config.OpenSearch.BulkIndexerFlushTimeout,
		BulkGetterBatchSize:     w.config.OpenSearch.BulkGetterBatchSize,
		BulkGetterBatchTimeout:  w.config.OpenSearch.BulkGetterBatchTimeout,
	}

	return os_client.New(config, w.Instrumentation)
}

func (w *Pool) getRedisClient() (*redis.Client, error) {
	config := redis.ClientConfig{
		Addrs: w.config.Redis.Addresses,
	}

	return redis.NewClient(&config, w.Instrumentation)
}

// getCachingFields() returns fields for caching based on fields in the indexTypes.Update struct.
func getCachingFields() []string {
	updateFields := reflect.VisibleFields(reflect.TypeOf(indexTypes.Update{}))
	cachingFields := make([]string, len(updateFields))
	for i, field := range updateFields {
		cachingFields[i] = field.Name
	}

	return cachingFields
}

func getOsIndex(c *os_client.Client, name string) index.Index {
	return os_index.New(
		c,
		&os_index.Config{Name: name},
	)
}

func (w *Pool) getIndexes(ctx context.Context) (*crawler.Indexes, error) {
	os, err := w.getOpenSearchClient()
	if err != nil {
		return nil, err
	}

	redis, err := w.getRedisClient()
	if err != nil {
		return nil, err
	}

	go osWorkLoop(ctx, os.Work)

	if err := redis.Start(ctx); err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		redis.Close(ctx)
	}()

	cfg := w.config.Indexes

	// TODO: Refactor/cleanup.
	osFactory := os_factory.New(os)
	fileOSIndex, err := osFactory.NewIndex(ctx, cfg.Files.Name)
	if err != nil {
		return nil, err
	}
	dirOSIndex, err := osFactory.NewIndex(ctx, cfg.Directories.Name)
	if err != nil {
		return nil, err
	}
	invalidOSIndex, err := osFactory.NewIndex(ctx, cfg.Invalids.Name)
	if err != nil {
		return nil, err
	}
	partialOSIndex, err := osFactory.NewIndex(ctx, cfg.Partials.Name)
	if err != nil {
		return nil, err
	}

	return &crawler.Indexes{
		Files: cache.New(
			fileOSIndex,
			redis.NewIndex(cfg.Files.Name, cfg.Files.Prefix, false),
			indexTypes.Update{},
			w.Instrumentation,
		),
		Directories: cache.New(
			dirOSIndex,
			redis.NewIndex(cfg.Directories.Name, cfg.Directories.Prefix, false),
			indexTypes.Update{},
			w.Instrumentation,
		),
		Invalids: cache.New(
			invalidOSIndex,
			redis.NewIndex(cfg.Invalids.Name, cfg.Invalids.Prefix, true),
			struct{}{},
			w.Instrumentation,
		),
		Partials: cache.New(
			partialOSIndex,
			redis.NewIndex(cfg.Partials.Name, cfg.Partials.Prefix, true),
			struct{}{},
			w.Instrumentation,
		),
	}, nil
}
