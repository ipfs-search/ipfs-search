package batchinggetter

import (
	"context"
	"log"
	"strings"

	"github.com/opensearch-project/opensearch-go"
)

// BatchingGetter allows batching/bulk gets.
type BatchingGetter struct {
	config Config
	queue  chan reqresp
}

// New returns a new BatchingGetter, setting sensible defaults for the configuration.
func New(cfg Config) *BatchingGetter {
	if cfg.Client == nil {
		cfg.Client, _ = opensearch.NewDefaultClient()
	}

	if cfg.BatchSize == 0 {
		cfg.BatchSize = 100
	}

	bg := BatchingGetter{
		config: cfg,
		queue:  make(chan reqresp, cfg.BatchSize),
	}

	return &bg
}

// Get queues a single Get() for a batching get.
func (bg *BatchingGetter) Get(ctx context.Context, req *GetRequest, dst interface{}) <-chan GetResponse {
	resp := make(chan GetResponse, 1)

	bg.queue <- reqresp{req, resp, dst}

	return resp
}

// StartWorker starts a single worker processing batched Get() requests. It will terminate on errors.
func (bg *BatchingGetter) StartWorker(ctx context.Context) error {
	b, err := bg.populateBatch(ctx, bg.queue)

	if err != nil {
		return err
	}

	return bg.performBatch(ctx, b)
}

type reqresp struct {
	req  *GetRequest
	resp chan GetResponse
	dst  interface{}
}

type batch map[string]map[string]bulkRequest

func getFieldsKey(fields []string) string {
	return strings.Join(fields, "")
}

func (bg *BatchingGetter) populateBatch(ctx context.Context, queue <-chan reqresp) (batch, error) {
	var b batch

	for i := 0; i < bg.config.BatchSize; i++ {
		log.Println(i)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case rr := <-queue:
			// Add batch.Add(fields, index, documentid)
			b[getFieldsKey(rr.req.Fields)][rr.req.Index][rr.req.DocumentID] = rr
		}
	}

	return b, nil
}

func (bg *BatchingGetter) performBatch(ctx context.Context, b batch) error {
	for _, indexes := range b {
		for _, r := range indexes {
			err := r.performBulkRequest(ctx, bg.config.Client)
			if err != nil {
				// Note: this will terminate batch on errors in any requests.
				return err
			}
		}
	}

	return nil
}

// Compile-time assurance that implementation satisfies interface.
var _ AsyncGetter = New(Config{})
