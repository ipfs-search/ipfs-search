package factory

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/streadway/amqp"
)

type CrawlFunc func(i *crawler.Indexable) func(context.Context) error

type Worker struct {
	*crawler.Crawler

	Delivery  *amqp.Delivery
	CrawlFunc CrawlFunc
}

func (c *Worker) Work(ctx context.Context) error {
	// Create an Indexable from the message's body
	i, err := c.IndexableFromJSON(c.Delivery.Body)
	if err != nil {
		return err
	}

	// Call crawler function with context
	return c.CrawlFunc(i)(ctx)
}
