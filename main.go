package main

import (
	"fmt"
	"github.com/dokterbob/ipfs-search/crawler"
	"github.com/dokterbob/ipfs-search/indexer"
	"github.com/dokterbob/ipfs-search/queue"
	"gopkg.in/ipfs/go-ipfs-api.v1"
	"gopkg.in/olivere/elastic.v3"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"time"
)

const (
	IPFS_API     = "localhost:5001"
	HASH_WORKERS = 40
	FILE_WORKERS = 0
	TIMEOUT      = 60 * time.Duration(time.Second)
	HASH_WAIT    = time.Duration(time.Second)
	FILE_WAIT    = HASH_WAIT
)

func main() {
	app := cli.NewApp()
	app.Name = "ipfs-search"
	app.Usage = "IPFS search engine."

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add `HASH` to crawler queue",
			Action:  add,
		},
		{
			Name:    "crawl",
			Aliases: []string{"c"},
			Usage:   "start crawler",
			Action:  crawl,
		},
	}

	app.Run(os.Args)
}

func get_elastic() (*elastic.Client, error) {
	el, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}
	exists, err := el.IndexExists("ipfs").Do()
	if err != nil {
		return nil, err
	}
	if !exists {
		// Index does not exist yet, create
		el.CreateIndex("ipfs")
	}

	return el, nil
}

func add(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("Please supply one hash as argument.", 1)
	}

	hash := c.Args().Get(0)

	fmt.Printf("Adding hash '%s' to queue\n", hash)

	ch, err := queue.NewChannel()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer ch.Close()

	queue, err := queue.NewTaskQueue(ch, "hashes")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err = queue.AddTask(map[string]interface{}{
		"hash": hash,
	})

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func crawl(c *cli.Context) error {
	// For now, assume gateway running on default host:port
	sh := shell.NewShell(IPFS_API)

	// Set 1 minute timeout on IPFS requests
	sh.SetTimeout(TIMEOUT)

	el, err := get_elastic()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	add_ch, err := queue.NewChannel()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer add_ch.Close()

	hq, err := queue.NewTaskQueue(add_ch, "hashes")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fq, err := queue.NewTaskQueue(add_ch, "files")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	id := indexer.NewIndexer(el)

	crawli := crawler.NewCrawler(sh, id, fq, hq)

	errc := make(chan error, 1)

	for i := 0; i < HASH_WORKERS; i++ {
		// Now create queues and channel for workers
		ch, err := queue.NewChannel()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		defer ch.Close()

		hq, err := queue.NewTaskQueue(ch, "hashes")
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		hq.StartConsumer(func(params interface{}) error {
			args := params.(*crawler.CrawlerArgs)

			return crawli.CrawlHash(
				args.Hash,
				args.Name,
				args.ParentHash,
				args.ParentName,
			)
		}, &crawler.CrawlerArgs{}, errc, true, add_ch)

		// Start workers timeout/hash time apart
		time.Sleep(HASH_WAIT)
	}

	for i := 0; i < FILE_WORKERS; i++ {
		ch, err := queue.NewChannel()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		defer ch.Close()

		fq, err := queue.NewTaskQueue(ch, "files")
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fq.StartConsumer(func(params interface{}) error {
			args := params.(*crawler.CrawlerArgs)

			return crawli.CrawlFile(
				args.Hash,
				args.Name,
				args.ParentHash,
				args.ParentName,
				args.Size,
			)
		}, &crawler.CrawlerArgs{}, errc, true, add_ch)

		// Start workers timeout/hash time apart
		time.Sleep(FILE_WAIT)
	}

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGQUIT)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for {
		select {
		case err = <-errc:
			log.Printf("%T: %v", err, err)
		}
	}

	// No error
	return nil
}
