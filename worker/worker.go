package worker

import (
	"errors"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/indexer"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs/go-ipfs-api"
	"time"
)

// Worker crawls hashes and files in parallel; it consumes hashes and files
// from respective queues and starts the configured number fo goroutines
// processing each.
type Worker struct {
	crawler      *crawler.Crawler
	config       *Config
	openChannels []*queue.TaskChannel // Channels to be closed later
	errc         chan<- error         // Channel the worker sends errors through
}

func newCrawler(config *Config, addCh *queue.TaskChannel) (*crawler.Crawler, error) {
	// Create tasks queue's
	// As there's potential failure, execute this first to allow quick fail
	hq, err := queue.NewTaskQueue(addCh, "hashes")
	if err != nil {
		addCh.Close()
		return nil, err
	}

	fq, err := queue.NewTaskQueue(addCh, "files")
	if err != nil {
		addCh.Close()
		return nil, err
	}

	// Create and configure Ipfs shell
	sh := shell.NewShell(config.IpfsAPI)
	sh.SetTimeout(config.IpfsTimeout)

	el, err := getElastic(config.ElasticSearchURL)
	if err != nil {
		addCh.Close()
		return nil, err
	}

	// Create elasticsearch indexer
	id := &indexer.Indexer{
		ElasticSearch: el,
	}

	c := &crawler.Crawler{
		Config:    config.CrawlerConfig,
		Shell:     sh,
		Indexer:   id,
		FileQueue: fq,
		HashQueue: hq,
	}

	return c, nil
}

// New returns an initialized worker
func New(config *Config, errc chan<- error) (*Worker, error) {
	// These is the channel the crawler uses to add newly crawled hashes
	addCh, err := queue.NewChannel()
	if err != nil {
		return nil, err
	}

	c, err := newCrawler(config, addCh)
	if err != nil {
		return nil, err
	}

	return &Worker{
		crawler:      c,
		config:       config,
		openChannels: []*queue.TaskChannel{addCh},
		errc:         errc,
	}, nil
}

// workerQueue creates a channel and named queue for a worker to consume
func (w *Worker) workerQueue(name string) (*queue.TaskQueue, error) {
	ch, err := queue.NewChannel()
	if err != nil {
		return nil, err
	}
	w.openChannels = append(w.openChannels, ch)

	q, err := queue.NewTaskQueue(ch, name)
	if err != nil {
		return nil, err
	}

	return q, nil
}

// startWorkers starts count workers for q executing crawlFunc and waiting wait between starting them
func (w *Worker) startWorkers(count int, qName string, crawlFunc func(params interface{}) error, wait time.Duration) error {
	for i := 0; i < count; i++ {
		q, err := w.workerQueue(qName)
		if err != nil {
			return err
		}

		consumer := &queue.Consumer{
			Func:    crawlFunc,
			ErrChan: w.errc,
			Queue:   q,
			Params:  &crawler.Args{},
		}

		consumer.Start()

		// Start workers timeout/hash time apart
		time.Sleep(wait)
	}

	return nil
}

func (w *Worker) paramsToIndexable(params interface{}) (*crawler.Indexable, error) {
	args, ok := params.(*crawler.Args)
	if !ok {
		return nil, errors.New("could not assert params as crawler.Args")
	}

	return &crawler.Indexable{
		Args:    args,
		Crawler: w.crawler,
	}, nil

}

func (w *Worker) crawlHash(params interface{}) error {
	i, err := w.paramsToIndexable(params)
	if err != nil {
		return err
	}

	return i.CrawlHash()
}

func (w *Worker) crawlFile(params interface{}) error {
	i, err := w.paramsToIndexable(params)
	if err != nil {
		return err
	}

	return i.CrawlFile()
}

// Start initiates crawling of the worker (asynchronuously)
func (w *Worker) Start() error {
	// Start hash worker
	err := w.startWorkers(w.config.HashWorkers, "hashes", w.crawlHash, w.config.HashWait)
	if err != nil {
		w.Close()
		return err
	}

	// Start file workers
	err = w.startWorkers(w.config.FileWorkers, "files", w.crawlFile, w.config.FileWait)
	if err != nil {
		w.Close()
		return err
	}

	return nil
}

// Close destroy closes worker channels and frees resources
func (w *Worker) Close() error {
	for _, channel := range w.openChannels {
		err := channel.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
