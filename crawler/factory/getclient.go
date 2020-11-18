package factory

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

func getClient() *http.Client {
	// TODO: Get more advanced client with circuit breaking etc. over manual
	// retrying get etc.
	// Ref: https://github.com/gojek/heimdall#creating-a-hystrix-like-circuit-breaker
	return http.Client{
		Timeout:   config.RequestTimeout,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}
