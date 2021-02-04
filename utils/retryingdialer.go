package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

const (
	retryWait = 2 * time.Second
	maxTries  = 60
)

// ErrRetriesExhausted signifies that the maximum amount of connection attempts have been exhausted while dialing.
var ErrRetriesExhausted = errors.New("Dial retries exhausted")

// RetryingDialer is a dialer which returns Dial and DialContext wrapped in a retrier when the requested connection is refused
// (e.g. the service is unavailable/still starting).
type RetryingDialer struct {
	net.Dialer
	context.Context
}

func (d *RetryingDialer) retrier(ctx context.Context, dial func() (net.Conn, error)) (net.Conn, error) {
	var (
		err error
		c   net.Conn
	)

	for tryCnt := 0; tryCnt < maxTries; tryCnt++ {
		c, err = dial()

		if err == nil {
			// Connected!
			return c, nil
		}

		if !errors.Is(err, syscall.ECONNREFUSED) {
			// Propagate any non-connection-refused errors.
			return c, err
		}

		log.Printf("Connection error (try %d of %d): %v, sleeping %s", tryCnt, maxTries, err, retryWait)
		select {
		case <-time.After(retryWait):
			// Wait or cancel when context is canceled.
		case <-ctx.Done():
			// TODO: Find out why this context is not canceled appropriately for ES client connection.
			return c, ctx.Err()
		}
	}

	return c, fmt.Errorf("%w:%T %v", ErrRetriesExhausted, err, err)
}

// Dial wraps net.Dialer.Dial so that it retries dials in case a connection is refused.
func (d *RetryingDialer) Dial(network, address string) (net.Conn, error) {
	return d.retrier(d.Context, func() (net.Conn, error) {
		return d.Dialer.Dial(network, address)
	})
}

// DialContext wraps net.Dialer.DialContext so that it retries dials in case a connection is refused.
func (d *RetryingDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return d.retrier(ctx, func() (net.Conn, error) {
		return d.Dialer.DialContext(ctx, network, address)
	})
}
