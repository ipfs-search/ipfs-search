package sniffer

import (
	"context"
	"github.com/ipfs/go-ipfs-api"
	"log"
)

func getProviders(ctx context.Context, l shell.Logger, providers chan<- Provider) error {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return context error
			return ctx.Err()
		default:
			// Note: this one is blocking, and might stall. We should have a timeout on this
			// or something!
			log.Printf("Waiting for next messaage")
			msg, err := l.Next()
			if err != nil {
				return err
			}

			provider, err := Message(msg).ResourceProvider()
			if err != nil {
				return err
			}

			if provider != nil {
				providers <- *provider
			}
		}
	}
}
