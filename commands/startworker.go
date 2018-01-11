package commands

import (
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

	config := &worker.Config{
		IpfsAPI:       "localhost:5001",
		ElasticSearch: el,
		HashWorkers:   140,
		FileWorkers:   120,
		IpfsTimeout:   360 * time.Duration(time.Second),
		HashWait:      time.Duration(100 * time.Millisecond),
		FileWait:      time.Duration(100 * time.Millisecond),
	}

	return config, nil
}

func startWorker(config *worker.Config) (errc chan error, err error) {
	worker, err := worker.New(config)
	if err != nil {
		return
	}

	errc, err = worker.Start()
	return
}

// StartWorker configures and initializes a worker
func StartWorker() (err error) {
	config, err := getWorkerConfig()
	if err != nil {
		return
	}

	errc, err := startWorker(config)
	if err != nil {
		return
	}

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
