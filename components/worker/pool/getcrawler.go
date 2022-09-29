package pool

import (
	"context"
	"log"

	"github.com/ipfs-search/ipfs-search/components/crawler"
)

func (p *Pool) getCrawler(ctx context.Context) (*crawler.Crawler, error) {
	var (
		queues  *crawler.Queues
		indexes *crawler.Indexes
		err     error
	)

	log.Println("Getting publish queues.")
	if queues, err = p.getQueues(ctx); err != nil {
		return nil, err
	}

	log.Println("Getting indexes.")
	if indexes, err = p.getIndexes(ctx); err != nil {
		return nil, err
	}

	protocol := p.getProtocol()
	extractors := p.getExtractors(protocol)
	config := p.config.CrawlerConfig()

	return crawler.New(config, indexes, queues, protocol, extractors, p.Instrumentation), nil
}
