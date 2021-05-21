# Quick-start

## Dependencies

* Go 1.13
* Elasticsearch 7.x
* RabbitMQ / AMQP server
* NodeJS 9.x
* IPFS 0.7

## Building
```bash
$ go get ./...
$ make
```

## Running

### Docker
The most convenient way to run the crawler is through Docker. Simply run:

```bash
docker-compose up
```

This will start the crawler, the sniffer and all its dependencies. Hashes can also be queued for crawling manually by running `ipfs-search a <hash>` from within the running container. For example:

```bash
docker-compose exec ipfs-crawler ipfs-search add QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv
```

### Ansible deployment
Automated deployment can be done on any (virtual) Ubuntu 16.04 machine. The full production stack is automated and can be found in it's own [repository](https://github.com/ipfs-search/ipfs-search-deployment).
