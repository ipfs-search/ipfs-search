module github.com/ipfs-search/ipfs-search

require (
	github.com/alanshaw/ipfs-hookds v0.3.0
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/dankinder/httpmock v1.0.1
	github.com/ipfs-search/go-env v0.0.0-20220928152343-588b5d46eac9
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-ipfs-api v0.3.0
	github.com/ipfs/go-unixfs v0.2.4
	github.com/jpillora/backoff v1.0.0
	github.com/libp2p/go-eventbus v0.2.1
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-kad-dht v0.10.0
	github.com/multiformats/go-base32 v0.0.3
	github.com/opensearch-project/opensearch-go/v2 v2.0.0
	github.com/rabbitmq/amqp091-go v1.3.4
	github.com/rueian/rueidis v0.0.77
	github.com/stretchr/testify v1.8.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.36.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/jaeger v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gopkg.in/urfave/cli.v1 v1.20.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/crackcomm/go-gitignore v0.0.0-20170627025303-887ab5e44cc3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/ipfs/bbloom v0.0.1 // indirect
	github.com/ipfs/go-block-format v0.0.2 // indirect
	github.com/ipfs/go-blockservice v0.1.0 // indirect
	github.com/ipfs/go-ipfs-blockstore v0.0.1 // indirect
	github.com/ipfs/go-ipfs-ds-help v0.0.1 // indirect
	github.com/ipfs/go-ipfs-exchange-interface v0.0.1 // indirect
	github.com/ipfs/go-ipfs-files v0.0.9 // indirect
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-ipld-cbor v0.0.2 // indirect
	github.com/ipfs/go-ipld-format v0.0.2 // indirect
	github.com/ipfs/go-log v1.0.4 // indirect
	github.com/ipfs/go-log/v2 v2.1.1 // indirect
	github.com/ipfs/go-merkledag v0.2.3 // indirect
	github.com/ipfs/go-metrics-interface v0.0.1 // indirect
	github.com/ipfs/go-verifcid v0.0.1 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/libp2p/go-buffer-pool v0.0.2 // indirect
	github.com/libp2p/go-flow-metrics v0.0.3 // indirect
	github.com/libp2p/go-openssl v0.0.7 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multiaddr v0.3.1 // indirect
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/multiformats/go-multihash v0.0.14 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polydawn/refmt v0.0.0-20190408063855-01bf1e26dd14 // indirect
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/whyrusleeping/tar-utils v0.0.0-20180509141711-8c6c8ba81d5c // indirect
	go.opencensus.io v0.22.4 // indirect
	go.opentelemetry.io/otel/metric v0.32.0 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20220704084225-05e143d24a9e // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/stretchr/testify => github.com/ipfs-search/testify v1.8.1-0.20220714120938-9ebebef47942

// +heroku goVersion go1.19
go 1.19
