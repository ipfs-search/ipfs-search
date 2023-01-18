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

const debug = true

// By default, do not retry.
const shouldRetry = false

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
func (w *Worker) Start(ctx context.Context, deliveries <-chan samqp.Delivery) error {
	ctx, span := w.Tracer.Start(ctx, "crawler.pool.startWorker")
	defer span.End()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-deliveries:
			if !ok {
				// This is a fatal error; it should never happen - crash the program!
				panic("unexpected channel close")
			}

			if err := w.crawlDelivery(ctx, d); err != nil {
				span.RecordError(err)
				return err
			}
		}
	}
}

func getResource(d samqp.Delivery) (*t.AnnotatedResource, error) {
	var err error

	r := &t.AnnotatedResource{
		Resource: &t.Resource{},
	}

	if err = json.Unmarshal(d.Body, r); err != nil {
		err = fmt.Errorf("Error unmarshalling delivery: %w", err)
	}

	if !r.IsValid() {
		err = fmt.Errorf("Invalid resource: %v", r)
	}

	return r, err
}

// ackOrRejec based on crawler result: returns nil on crawler errors, rejects relevant deliveries instead.
func ackOrReject(err error, d samqp.Delivery) error {
	if err != nil {
		// Crawler error
		if err := d.Reject(shouldRetry); err != nil {
			// Error rejecting delivery,
			return err
		}

		// Crawl errors are expected: do not propagate error.
		return nil
	}

	if err := d.Ack(false); err != nil {
		return err
	}

	return nil
}

func (w *Worker) crawlDelivery(ctx context.Context, d samqp.Delivery) error {
	ctx, span := w.Tracer.Start(ctx, "crawler.pool.crawlDelivery", trace.WithNewRoot())
	defer span.End()

	// Errors in load limiter are unexpected; propagate them.
	if err := w.ll.LoadLimit(); err != nil {
		return err
	}

	// Errors in resource getter are unexpected; propagate them.
	r, err := getResource(d)
	if err != nil {
		return err
	}

	log.Printf("worker: Crawling '%s'", r)

	err = w.crawler.Crawl(ctx, r)
	if err != nil {
		if debug {
			log.Printf("worker: Error crawling '%s': %v", r, err)
		}
		span.RecordError(err)
	}

	err = ackOrReject(err, d)
	// BUG! This error *should* propagate through the worker to the error group and cause the crawler to crash.
	// But it doens't, leaving the crawler in an unrecoverable state.
	if err != nil {
		if debug {
			log.Printf("worker: Delivery Ack/Reject error '%s': %v", r, err)
		}
		span.RecordError(err)
	}

	if debug {
		log.Printf("worker: Done crawling '%s'", r)
	}

	return err
}
