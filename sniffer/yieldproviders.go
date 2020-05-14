package sniffer

import (
	"context"
	"errors"
	t "github.com/ipfs-search/ipfs-search/types"
	"time"
)

// ErrorLoggerTimeout represents a timeout from the IPFS shell's logger.
var ErrorLoggerTimeout = errors.New("Timeout waiting for log messages")

// The default IPFS logger is a blocking function without a context, hence
// we wrap it in a goroutine to allow for timeouts.
// TODO: Upgrade to well-designed `go-ipfs-http-api` if and when Logger is
// implemented there and/or to use the generic `Request()` from there.
func loggerToChannel(ctx context.Context, l Logger, msgs chan<- map[string]interface{}, errc chan<- error) {
	for {
		select {
		case <-ctx.Done():
			errc <- ctx.Err()
			return
		default:
			msg, err := l.Next()
			if err != nil {
				errc <- err
			}

			msgs <- msg
		}
	}
}

type providerYielder struct {
	e       Extractor
	timeout time.Duration
}

func (y *providerYielder) yield(ctx context.Context, l Logger, providers chan<- t.Provider) error {
	msgs := make(chan map[string]interface{})
	errc := make(chan error, 1)

	go loggerToChannel(ctx, l, msgs, errc)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(y.timeout):
			return ErrorLoggerTimeout
		case err := <-errc:
			return err
		case msg := <-msgs:
			provider, err := y.e.Extract(msg)
			if err != nil {
				return err
			}

			if provider != nil {
				providers <- *provider
			}
		}
	}
}
