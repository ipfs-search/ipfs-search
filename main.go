/*

Search engine for IPFS using Elasticsearch, RabbitMQ and Tika.
*/
package main

import (
	"fmt"
	"github.com/ipfs-search/ipfs-search/commands"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

func main() {
	// Prefix logging with filename and line number: "d.go:23"
	// log.SetFlags(log.Lshortfile)

	// Logging w/o prefix
	log.SetFlags(0)

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

func add(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("Please supply one hash as argument.", 1)
	}

	hash := c.Args().Get(0)

	fmt.Printf("Adding hash '%s' to queue\n", hash)

	err := commands.AddHash(hash)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func crawl(c *cli.Context) error {
	fmt.Printf("Starting worker\n")

	err := commands.StartWorker()

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
