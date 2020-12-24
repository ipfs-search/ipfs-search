/*

Search engine for IPFS using Elasticsearch, RabbitMQ and Tika.
*/
package main

import (
	"context"
	"fmt"
	"github.com/ipfs-search/ipfs-search/commands"
	"github.com/ipfs-search/ipfs-search/config"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"os/signal"
	"syscall"
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
		{
			Name:    "config",
			Aliases: []string{},
			Usage:   "configuration",
			Subcommands: []cli.Command{
				{
					Name:   "generate",
					Usage:  "generate default configuration",
					Action: generateConfig,
				},
				{
					Name:   "check",
					Usage:  "check configuration",
					Action: checkConfig,
				},
				{
					Name:   "dump",
					Usage:  "dump current configuration to stdout",
					Action: dumpConfig,
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getConfig(c *cli.Context) (*config.Config, error) {
	configFile := c.GlobalString("config")

	cfg, err := config.Get(configFile)
	if err != nil {
		return nil, err
	}

	err = cfg.Check()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func checkConfig(c *cli.Context) error {
	_, err := getConfig(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("Configuration checked.")

	return nil
}

func generateConfig(c *cli.Context) error {
	cfg := config.Default()

	configFile := c.GlobalString("config")
	if configFile == "" {
		return cli.NewExitError("Configuration file not specified. Use the \"-c\" option.", 1)
	}

	fmt.Printf("Writing default configuration to: %s\n", configFile)
	return cfg.Write(configFile)
}

func dumpConfig(c *cli.Context) error {
	cfg, err := getConfig(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return cfg.Dump()
}

func add(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())

	// Allow SIGTERM / Control-C quit through context
	onSigTerm(cancel)

	if c.NArg() != 1 {
		return cli.NewExitError("Please supply one hash as argument.", 1)
	}
	hash := c.Args().Get(0)

	cfg, err := getConfig(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("Adding hash '%s' to queue\n", hash)

	err = commands.AddHash(ctx, cfg, hash)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

// onSigTerm calls f() when SIGTERM (control-C) is received
func onSigTerm(f func()) {
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var fail = func() {
		<-sigChan
		os.Exit(1)
	}

	var quit = func() {
		<-sigChan

		go fail()

		fmt.Println("Received SIGTERM, quitting... One more SIGTERM and we'll abort!")
		f()
	}

	go quit()
}

func crawl(c *cli.Context) error {
	fmt.Println("Starting worker")

	ctx, cancel := context.WithCancel(context.Background())

	// Allow SIGTERM / Control-C quit through context
	onSigTerm(cancel)

	cfg, err := getConfig(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err = commands.Crawl(ctx, cfg)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
