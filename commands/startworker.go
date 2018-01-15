package commands

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/worker"
	"log"
	"time"
)

// getWorkgerConfig sets up configuration for worker
func getWorkerConfig() (*worker.Config, error) {
	el, err := getElastic()
	if err != nil {
		return nil, err
	}

	crawlerConfig := &crawler.Config{
		IpfsTikaURL:     "http://localhost:8081",
		IpfsTikaTimeout: 300 * time.Duration(time.Second),
		RetryWait:       2 * time.Duration(time.Second),
		MetadataMaxSize: 50 * 1024 * 1024,
		PartialSize:     262144,
	}

	config := &worker.Config{
		IpfsAPI:       "localhost:5001",
		ElasticSearch: el,
		HashWorkers:   140,
		FileWorkers:   120,
		IpfsTimeout:   360 * time.Duration(time.Second),
		HashWait:      time.Duration(100 * time.Millisecond),
		FileWait:      time.Duration(100 * time.Millisecond),
		CrawlerConfig: crawlerConfig,
	}

	return config, nil
}

func startWorker(config *worker.Config, errc chan<- error) (*worker.Worker, error) {
	w, err := worker.New(config, errc)
	if err != nil {
		return nil, err
	}

	err = w.Start()
	if err != nil {
		return nil, err
	}

	return w, nil
}

// StartWorker configures and initializes a worker
func StartWorker() (err error) {
	config, err := getWorkerConfig()
	if err != nil {
		return
	}

	errc := make(chan error, 1)

	worker, err := startWorker(config, errc)
	if err != nil {
		return
	}
	defer worker.Close()

	// TODO: Catch QUIT signal here and create shutdown channel to properly
	// exit crawler. This would involve implementing a stop channel all the
	// way down to the queue consumers.

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// Start block on errc, logging messages
	for {
		select {
		case err = <-errc:
			log.Printf("%T: %v", err, err)
		}
	}
}
