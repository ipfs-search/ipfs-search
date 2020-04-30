package sniffer

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs-api"
	"log"
	"time"
)

func getProviders(ctx context.Context, l shell.Logger, providers chan<- Provider, timeout time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(timeout):
			return fmt.Errorf("Timeout (%s) waiting for log messages", timeout)
		default:
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
