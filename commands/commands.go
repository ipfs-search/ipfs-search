package commands

import (
	"fmt"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/indexer"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs/go-ipfs-api"
	"time"
)

// TODO: Read this from configuration file.
const (
	ipfsAPI     = "localhost:5001"
	hashWorkers = 140
	fileWorkers = 120
	ipfsTimeout = 360 * time.Duration(time.Second)      // Timeout for IPFS gateway HTTPS requests
	hashWait    = time.Duration(100 * time.Millisecond) // Time to wait between creating hash workers
	fileWait    = hashWait                              // Time to wait between creating file workers
)

// AddHash queues a single IPFS hash for indexing
func AddHash(hash string) error {
	fmt.Printf("Adding hash '%s' to queue\n", hash)

	ch, err := queue.NewChannel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := queue.NewTaskQueue(ch, "hashes")
	if err != nil {
		return err
	}

	err = queue.AddTask(map[string]interface{}{
		"hash": hash,
	})

	return err
}

// Crawl starts crawling hashes and files.
// Returns an error if initialization fails or error channel with non-fatal
// crawling errors.
func Crawl() (chan error, error) {
	// For now, assume gateway running on default host:port
	sh := shell.NewShell(ipfsAPI)

	// Set 1 minute timeout on IPFS requests
	sh.SetTimeout(ipfsTimeout)

	el, err := getElastic()
	if err != nil {
		return nil, err
	}

	addCh, err := queue.NewChannel()
	if err != nil {
		return nil, err
	}
	defer addCh.Close()

	hq, err := queue.NewTaskQueue(addCh, "hashes")
	if err != nil {
		return nil, err
	}

	fq, err := queue.NewTaskQueue(addCh, "files")
	if err != nil {
		return nil, err
	}

	id := indexer.NewIndexer(el)

	crawli := crawler.NewCrawler(sh, id, fq, hq)

	errc := make(chan error, 1)
	for i := 0; i < hashWorkers; i++ {
		// Now create queues and channel for workers
		ch, err := queue.NewChannel()
		if err != nil {
			return nil, err
		}
		defer ch.Close()

		hq, err := queue.NewTaskQueue(ch, "hashes")
		if err != nil {
			return nil, err
		}

		hq.StartConsumer(func(params interface{}) error {
			args := params.(*crawler.Args)

			return crawli.CrawlHash(
				args.Hash,
				args.Name,
				args.ParentHash,
				args.ParentName,
			)
		}, &crawler.Args{}, errc)

		// Start workers timeout/hash time apart
		time.Sleep(hashWait)
	}

	for i := 0; i < fileWorkers; i++ {
		ch, err := queue.NewChannel()
		if err != nil {
			return nil, err
		}
		defer ch.Close()

		fq, err := queue.NewTaskQueue(ch, "files")
		if err != nil {
			return nil, err
		}

		fq.StartConsumer(func(params interface{}) error {
			args := params.(*crawler.Args)

			return crawli.CrawlFile(
				args.Hash,
				args.Name,
				args.ParentHash,
				args.ParentName,
				args.Size,
			)
		}, &crawler.Args{}, errc)

		// Start workers timeout/hash time apart
		time.Sleep(fileWait)
	}

	return errc, nil
}
