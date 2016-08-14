package main

import (
	"fmt"
	machinery "github.com/RichardKnop/machinery/v1"
	machinery_config "github.com/RichardKnop/machinery/v1/config"
	signatures "github.com/RichardKnop/machinery/v1/signatures"
	"github.com/dokterbob/ipfs-search/crawler"
	"github.com/dokterbob/ipfs-search/indexer"
	"gopkg.in/ipfs/go-ipfs-api.v1"
	"gopkg.in/olivere/elastic.v3"
	"gopkg.in/urfave/cli.v1"
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

func get_machinery() (*machinery.Server, error) {
	cnf := machinery_config.Config{
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
		// 	Exchange:      *exchange,
		// 	ExchangeType:  *exchangeType,
		// 	DefaultQueue:  *defaultQueue,
		// 	BindingKey:    *bindingKey,
	}
	server, err := machinery.NewServer(&cnf)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func add(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("Please supply one hash as argument.", 1)
	}

	hash := c.Args().Get(0)

	fmt.Printf("Adding hash '%s' to queue\n", hash)

	server, err := get_machinery()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	task := signatures.TaskSignature{
		Name: "crawl",
		Args: []signatures.TaskArg{
			signatures.TaskArg{
				Type:  "string",
				Value: hash,
			},
		},
	}
	asyncResult, err := server.SendTask(&task)
	if err != nil {
		// failed to send the task
		return cli.NewExitError(err.Error(), 1)
	}

	taskState := asyncResult.GetState()
	fmt.Printf("Current state of %v task is:\n", taskState.TaskUUID)
	fmt.Println(taskState.State)

	return nil
}

func crawl(c *cli.Context) error {
	// For now, assume gateway running on default host:port
	sh := shell.NewShell("localhost:5001")

	el, err := get_elastic()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	server, err := get_machinery()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	id := indexer.NewIndexer(el)
	crawli := crawler.NewCrawler(sh, id, server)

	server.RegisterTask("crawl", func(hash string) error {
		return crawli.CrawlHash(hash)
	})

	worker := server.NewWorker("crawler")
	err = worker.Launch()
	if err != nil {
		// do something with the error
		return cli.NewExitError(err.Error(), 1)
	}

	// No error
	return nil
}
