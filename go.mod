module github.com/ipfs-search/ipfs-search

require (
	github.com/Netflix/go-env v0.0.0-20210116210345-8f74e74141f7
	github.com/alanshaw/ipfs-hookds v0.3.0
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/dankinder/httpmock v1.0.1
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-ipfs-api v0.3.0
	github.com/ipfs/go-unixfs v0.2.4
	github.com/jpillora/backoff v1.0.0
	github.com/kr/text v0.2.0 // indirect
	github.com/libp2p/go-eventbus v0.2.1
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-kad-dht v0.10.0
	github.com/multiformats/go-base32 v0.0.3
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opensearch-project/opensearch-go/v2 v2.0.0
	github.com/rabbitmq/amqp091-go v1.3.4
	github.com/stretchr/testify v1.8.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.36.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/jaeger v1.10.0
	go.opentelemetry.io/otel/metric v0.32.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/urfave/cli.v1 v1.20.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/stretchr/testify => github.com/ipfs-search/testify v1.8.1-0.20220714120938-9ebebef47942

// +heroku goVersion go1.16
go 1.16
