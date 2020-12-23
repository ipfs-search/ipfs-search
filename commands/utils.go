package commands

import (
	"net/http"

	"github.com/olivere/elastic/v7"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func getHttpClient() *http.Client {
	// TODO: Get more advanced client with circuit breaking etc. over manual
	// retrying get etc.
	// Ref: https://github.com/gojek/heimdall#creating-a-hystrix-like-circuit-breaker
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func getElasticClient(esURL string) (*elastic.Client, error) {
	httpClient := getHttpClient()

	return elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(esURL),
		elastic.SetHttpClient(httpClient),
	)
}
