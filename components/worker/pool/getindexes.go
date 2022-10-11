package pool

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/ipfs-search/ipfs-search/components/crawler"
	"github.com/ipfs-search/ipfs-search/components/index/cache"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch"
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

func (w *Pool) getOpenSearchClient() (*opensearch.Client, error) {
	config := &opensearch.ClientConfig{
		URL:       w.config.OpenSearch.URL,
		Transport: utils.GetHTTPTransport(w.dialer.DialContext, 100),
		Debug:     false,

		BulkIndexerWorkers:     w.config.OpenSearch.BulkIndexerWorkers,
		BulkIndexerFlushBytes:  int(w.config.OpenSearch.BulkIndexerFlushBytes),
		BulkGetterBatchSize:    w.config.OpenSearch.BulkGetterBatchSize,
		BulkGetterBatchTimeout: w.config.OpenSearch.BulkGetterBatchTimeout,
	}

	return opensearch.NewClient(config, w.Instrumentation)
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

func (w *Pool) getIndexes(ctx context.Context) (*crawler.Indexes, error) {
	osClient, err := w.getOpenSearchClient()
	if err != nil {
		return nil, err
	}

	redisClient, err := w.getRedisClient()
	if err != nil {
		return nil, err
	}

	go osWorkLoop(ctx, osClient.Work)

	redisClient.Start(ctx)
	go func() {
		<-ctx.Done()
		redisClient.Close(ctx)
	}()

	cacheCfg := &cache.Config{}

	return &crawler.Indexes{
		Files: cache.New(opensearch.New(
			osClient,
			&opensearch.Config{Name: w.config.Indexes.Files.Name},
		), redis.New(
			redisClient,
			&redis.Config{Name: w.config.Indexes.Files.Name, Prefix: w.config.Indexes.Files.Prefix},
		), indexTypes.Update{}, cacheCfg, w.Instrumentation),
		Directories: cache.New(opensearch.New(
			osClient,
			&opensearch.Config{Name: w.config.Indexes.Directories.Name},
		), redis.New(
			redisClient,
			&redis.Config{Name: w.config.Indexes.Directories.Name, Prefix: w.config.Indexes.Directories.Prefix},
		), indexTypes.Update{}, cacheCfg, w.Instrumentation),
		Invalids: cache.New(opensearch.New(
			osClient,
			&opensearch.Config{Name: w.config.Indexes.Invalids.Name},
		), redis.NewExistsIndex(
			redisClient,
			&redis.Config{Name: w.config.Indexes.Invalids.Name, Prefix: w.config.Indexes.Invalids.Prefix},
		), struct{}{}, cacheCfg, w.Instrumentation),
		Partials: cache.New(opensearch.New(
			osClient,
			&opensearch.Config{Name: w.config.Indexes.Partials.Name},
		), redis.NewExistsIndex(
			redisClient,
			&redis.Config{Name: w.config.Indexes.Partials.Name, Prefix: w.config.Indexes.Partials.Prefix},
		), struct{}{}, cacheCfg, w.Instrumentation),
	}, nil
}
