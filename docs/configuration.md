# Configuration

Configuration can be done using a YAML configuration file, or by specifying the following environment variables:
* `IPFS_API_URL`
* `IPFS_GATEWAY_URL`
* `ELASTICSEARCH_URL`
* `AMQP_URL`
* `TIKA_EXTRACTOR`
* `OTEL_TRACE_SAMPLER_ARG`
* `OTEL_EXPORTER_JAEGER_ENDPOINT`
* `HASH_WORKERS`
* `FILE_WORKERS`
* `DIRECTORY_WORKERS`
* `SNIFFER_LASTSEEN_EXPIRATION`
* `SNIFFER_LASTSEEN_PRUNELEN`
* `SNIFFER_BUFFER_SIZE`

A default configuration can be generated with:
```bash
ipfs-search -c config.yml config generate
```
(substitute `config.yml` with the configuration file you'd like to use.)

To use a configuration file, it is necessary to specify the `-c` option, as in:
```bash
ipfs-search -c config.yml crawl
```

The configuration can be (rudimentarily) checked with:
```bash
ipfs-search -c config.yml config check
```


## Annotated default configuration
```yaml
ipfs:
  api_url: http://localhost:5001                      # IPFS API endpoint, also IPFS_API_URL in env
  gateway_url: http://localhost:8080                  # IPFS gateway, also IPFS_GATEWAY_URL in env
  partial_size: 256KB                                 # Size of items considered to be partial (when unreferenced)
elasticsearch:
  url: http://localhost:9200                          # Also ELASTICSEARCH_URL in env
amqp:
  url: amqp://guest:guest@localhost:5672/             # Also AMQP_URL in env.
  max_reconnect: 100                                  # Maximum number of reconnect attempts
  reconnect_time: 2s                                  # Time to wait between reconnects
tika:
  url: http://localhost:8081                          # tika-extractor endpoint URL, also TIKA_EXTRACTOR in environment.
  timeout: 5m                                         # Timeout for requests to tika-extractor.
  max_file_size: 4GB                                  # Don't attempt to extract metadata for resources larger than this.
instrumentation:
  sampling_ratio: 0.01                                # Ratio of requests to sample for tracing. OTEL_TRACE_SAMPLER_ARG in env.
  jaeger_endpoint: http://localhost:14268/api/traces  # HTTP jaeger.thrift endpoint for tracing. OTEL_EXPORTER_JAEGER_ENDPOINT in env.
crawler:
  direntry_buffer_size: 8192                          # Buffer this many directory entries between listing and queue'ing
  min_update_age: 1h                                  # Minimum time between updating `last-seen` on objects.
  stat_timeout: 1m                                    # Request timeout for Stat() calls.
  direntry_timeout: 1m                                # Request timeout for Ls() calls.
  max_dirsize: 32768                                  # Don't index directories larger than this (contained items will be queue'd nonetheless).
sniffer:
  lastseen_expiration: 1h                             # Expire items in lastseen/dedup buffer after this time. SNIFFER_LASTSEEN_EXPIRATION in env.
  lastseen_prunelen: 32768                            # Expire lastseen buffer when size exceeds this. SNIFFER_LASTSEEN_PRUNELEN in env.
  logger_timeout: 1m                                  # Throw timeout error when no log messages arrive
  buffer_size: 512                                    # Size of the channels buffering between yielder, filter and adder. SNIFFER_BUFFER_SIZE in env.
indexes:
  files:
    name: ipfs_files                                  # Name of ES index to use.
  directories:
    name: ipfs_directories
  invalids:
    name: ipfs_invalids
queues:
  files:
    name: files                                       # Name of RabbitMQ queue to use.
  directories:
    name: directories
  hashes:
    name: hashes
workers:
  hash_workers: 70                                    # Amount of workers for various resources. Also HASH_WORKERS in env.
  file_workers: 120                                   # Also FILE_WORKERS in env.
  directory_workers: 70                               # Also DIRECTORY in env.
```
