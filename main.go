package main

import (
	"context"
	"fmt"
	"log"

	shell "github.com/ipfs/go-ipfs-api"
)

var ipfsURL = "localhost:5001"

func printlogs(ctx context.Context, sh *shell.Shell) error {
	log.Printf("Opening logger")
	logger, err := sh.GetLogs(ctx)
	if err != nil {
		return err
	}
	defer func() {
		log.Printf("Closing logger")

		err := logger.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("Printing log messages")
	for {
		msg, err := logger.Next()
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", msg)
	}

	return nil
}

func main() {
	// Open shell
	sh := shell.NewShell(ipfsURL)

	// Create context
	ctx := context.Background()

	err := printlogs(ctx, sh)
	if err != nil {
		log.Fatal(err)
	}
}
