package utils

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// GetHTTPTransport initializes a HTTP transport with OpenTelemetry transport for tracing.
func GetHTTPTransport(dialcontext func(ctx context.Context, network, address string) (net.Conn, error), maxConns int) http.RoundTripper {
	return otelhttp.NewTransport(&http.Transport{
		Proxy:               nil,
		DialContext:         dialcontext,
		ForceAttemptHTTP2:   false,
		MaxIdleConns:        maxConns,
		MaxIdleConnsPerHost: maxConns,
		IdleConnTimeout:     90 * time.Second,
	})
}
