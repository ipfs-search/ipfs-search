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
	name    string
	crawler *crawler.Crawler

	*instr.Instrumentation
}

// New returns a new worker.
func New(name string, crawler *crawler.Crawler, i *instr.Instrumentation) *Worker {
	return &Worker{
		name, crawler, i,
	}
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
