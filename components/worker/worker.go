package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	samqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/trace"

	"github.com/ipfs-search/ipfs-search/components/crawler"
	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

// Worker crawls deliveries from a queue.
type Worker struct {
	cfg     *Config
	crawler *crawler.Crawler
	ll      LoadLimiter

	*instr.Instrumentation
}

// New returns a new worker.
func New(cfg *Config, crawler *crawler.Crawler, i *instr.Instrumentation) *Worker {
	ll := NewLoadLimiter(cfg.Name, cfg.MaxLoadRatio, cfg.ThrottleMin, cfg.ThrottleMax)

	//# 0.8, 10*time.Second, 5*time.Minute)

	return &Worker{
		cfg, crawler, ll, i,
	}
}

// String returns the name of the worker.
func (w *Worker) String() string {
	return w.cfg.Name
}

// Start crawling deliveries, synchronously.
func (w *Worker) Start(ctx context.Context, deliveries <-chan samqp.Delivery) {
	ctx, span := w.Tracer.Start(ctx, "crawler.pool.startWorker")
	defer span.End()

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-deliveries:
			if !ok {
				// This is a fatal error; it should never happen - crash the program!
				panic("unexpected channel close")
			}

			if err := w.ll.LoadLimit(); err != nil {
				log.Printf("load limit exception: %s", err)
				span.RecordError(err)
			}

			if err := w.crawlDelivery(ctx, d); err != nil {
				// By default, do not retry.
				shouldRetry := false

				span.RecordError(err)

				if err := d.Reject(shouldRetry); err != nil {
					span.RecordError(err)
				}
			} else {
				if err := d.Ack(false); err != nil {
					span.RecordError(err)
				}
			}
		}
	}
}

func (w *Worker) crawlDelivery(ctx context.Context, d samqp.Delivery) error {
	ctx, span := w.Tracer.Start(ctx, "crawler.pool.crawlDelivery", trace.WithNewRoot())
	defer span.End()

	r := &t.AnnotatedResource{
		Resource: &t.Resource{},
	}

	if err := json.Unmarshal(d.Body, r); err != nil {
		span.RecordError(err)
		return err
	}

	if !r.IsValid() {
		err := fmt.Errorf("Invalid resource: %v", r)
		span.RecordError(err)
		return err
	}

	log.Printf("Crawling '%s'", r)
	err := w.crawler.Crawl(ctx, r)
	log.Printf("Done crawling '%s', result: %v", r, err)

	if err != nil {
		span.RecordError(err)
	}

	return err
}
