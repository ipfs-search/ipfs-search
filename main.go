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
	sh := shell.NewShell("localhost:5001")

	el, err := get_elastic()

	ch, err := queue.NewChannel()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer ch.Close()

	hq, err := queue.NewTaskQueue(ch, "hashes")
	fq, err := queue.NewTaskQueue(ch, "files")
	if err != nil {
		// do something with the error
		return cli.NewExitError(err.Error(), 1)
	}

	id := indexer.NewIndexer(el)

	crawler := crawler.NewCrawler(sh, id, fq, hq)

	errc := make(chan error, 1)

	hq.StartConsumer(func(params map[string]interface{}) error {
		// TODO: Assert hash in map, ideally by using custom type
		return crawler.CrawlHash(params["hash"].(string))
	}, errc)

	fq.StartConsumer(func(params map[string]interface{}) error {
		// TODO: Assert hash in map, ideally by using custom type
		return crawler.CrawlFile(params["hash"].(string))
	}, errc)

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGQUIT)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for {
		select {
		case err = <-errc:
			log.Println(err.Error())
		}
	}

	// No error
	return nil
}
