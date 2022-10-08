package pool

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/ipfs-search/ipfs-search/components/crawler"
	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/components/index/cache"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch"
	"github.com/ipfs-search/ipfs-search/components/index/redis"
	indexTypes "github.com/ipfs-search/ipfs-search/components/index/types"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/utils"
)

func startOpenSearchWorker(ctx context.Context, osClient *opensearch.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := osClient.Work(ctx); err != nil {
				log.Printf("Error in ES client pool, restarting pool: %s", err)
				// Prevent overly tight restart loop
				time.Sleep(time.Second)
			}
		}
	}
}

func ensureRedisClose(ctx context.Context, client *redis.Client) {
	<-ctx.Done()
	client.Close()
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

type indexFactory struct {
	osClient    *opensearch.Client
	redisClient *redis.Client
	cacheCfg    *cache.Config
	*instr.Instrumentation
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

func (f *indexFactory) getIndex(name string) index.Index {
	osIndex := opensearch.New(
		f.osClient,
		&opensearch.Config{Name: name},
	)

	redisIndex := redis.New(
		f.redisClient,
		&redis.Config{Name: name},
	)

	return cache.New(osIndex, redisIndex, indexTypes.Update{}, f.cacheCfg, f.Instrumentation)
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

	go ensureRedisClose(ctx, redisClient)
	go startOpenSearchWorker(ctx, osClient)

	cacheCfg := &cache.Config{}

	iFactory := indexFactory{
		osClient, redisClient, cacheCfg, w.Instrumentation,
	}

	// Note: Manually adjust indexCount whenever the amount of indexes change
	const indexCount = 4
	var (
		indexNames = [indexCount]string{w.config.Indexes.Files.Name, w.config.Indexes.Directories.Name, w.config.Indexes.Invalids.Name, w.config.Indexes.Partials.Name}
		indexes    [indexCount]index.Index
	)

	for i, name := range indexNames {
		indexes[i] = iFactory.getIndex(name)
	}

	// Note: Manually adjust order here!
	return &crawler.Indexes{
		Files:       indexes[0],
		Directories: indexes[1],
		Invalids:    indexes[2],
		Partials:    indexes[3],
	}, nil
}
