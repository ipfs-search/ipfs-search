package bulkgetter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/opensearch-project/opensearch-go"
)

// BulkGetter allows batching/bulk gets.
type BulkGetter struct {
	cfg   Config
	queue chan reqresp
}

// New returns a new BulkGetter, setting sensible defaults for the configuration.
func New(cfg Config) *BulkGetter {
	if cfg.Client == nil {
		cfg.Client, _ = opensearch.NewDefaultClient()
	}

	if cfg.BatchSize == 0 {
		cfg.BatchSize = 100
	}

	if cfg.BatchTimeout == 0 {
		cfg.BatchTimeout = 100 * time.Millisecond
	}

	bg := BulkGetter{
		cfg:   cfg,
		queue: make(chan reqresp, 5*cfg.BatchSize),
	}

	return &bg
}

// Get queues a single Get() for a batching get.
func (bg *BulkGetter) Get(ctx context.Context, req *GetRequest, dst interface{}) <-chan GetResponse {
	resp := make(chan GetResponse, 1)

	bg.queue <- reqresp{req, resp, dst}

	return resp
}

// Work starts a single worker processing batched Get() requests. It will terminate on errors.
func (bg *BulkGetter) Work(ctx context.Context) error {
	var err error

	log.Println("Starting worker for BulkGetter.")

	for err == nil {
		err = bg.processBatch(ctx)
	}

	log.Printf("BulkGetter worker exiting, error: %s", err)

	return err
}

func (bg *BulkGetter) processBatch(ctx context.Context) error {
	var (
		b   *bulkRequest
		err error
	)

	if b, err = bg.populateBatch(ctx, bg.queue); err != nil {
		err = fmt.Errorf("error populating batch: %w", err)
		return err
	}

	if len(b.rrs) == 0 {
		return nil
	}

	return b.execute()
}

func (bg *BulkGetter) populateBatch(ctx context.Context, queue <-chan reqresp) (*bulkRequest, error) {
	// log.Println("Populating BulkGetter batch.")

	b := newBulkRequest(ctx, bg.cfg.Client, bg.cfg.BatchSize)

	for i := 0; i < bg.cfg.BatchSize; i++ {
		select {
		case <-ctx.Done():
			return b, ctx.Err()
		case <-time.After(bg.cfg.BatchTimeout):
			// log.Printf("Batch timeout, %d elements", len(b.rrs))

			return b, nil
		case rr := <-queue:
			// log.Printf("Batch add, %d elements", len(b.rrs))

			if err := b.add(rr); err != nil {
				return b, err
			}
		}
	}

	return b, nil
}

// Compile-time assurance that implementation satisfies interface.
var _ AsyncGetter = New(Config{})
