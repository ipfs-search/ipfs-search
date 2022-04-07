package bulkgetter

import (
	"context"
	"log"
	"time"

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
	var err error

	for err != nil {
		err = bg.processBatch(ctx)
	}

	return err
}

func (bg *BatchingGetter) processBatch(ctx context.Context) error {
	b, err := bg.populateBatch(ctx, bg.queue)

	if err != nil {
		return err
	}

	return b.execute(ctx, bg.config.Client)
}

func (bg *BatchingGetter) populateBatch(ctx context.Context, queue <-chan reqresp) (batch, error) {
	b := newBatch()

	for i := 0; i < bg.config.BatchSize; i++ {
		log.Println(i)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(bg.config.BatchTimeout):
			return b, nil
		case rr := <-queue:
			b.add(rr)
		}
	}

	return b, nil
}

// Compile-time assurance that implementation satisfies interface.
var _ AsyncGetter = New(Config{})
