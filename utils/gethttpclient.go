package utils

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// GetHTTPClient initializes a HTTP client with OpenTelemetry transport for tracing.
func GetHTTPClient(dialcontext func(ctx context.Context, network, address string) (net.Conn, error)) *http.Client {
	transport := otelhttp.NewTransport(&http.Transport{
		Proxy:               nil,
		DialContext:         dialcontext,
		ForceAttemptHTTP2:   false,
		MaxIdleConns:        100, // Defaut
		MaxIdleConnsPerHost: 10,  // 10 connections to IPFS/Tika/ES (each).
		IdleConnTimeout:     90 * time.Second,
	})

	return &http.Client{
		Transport: transport,
	}
}
