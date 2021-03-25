# [ipfs-search.com](http://ipfs-search.com)
[![Build Status](https://travis-ci.org/ipfs-search/ipfs-search.svg?branch=master)](https://travis-ci.org/ipfs-search/ipfs-search)
[![Docker Build Status](https://img.shields.io/docker/build/ipfssearch/ipfs-search)](https://hub.docker.com/repository/docker/ipfssearch/ipfs-search)
[![Maintainability](https://api.codeclimate.com/v1/badges/1c25261992991d72137c/maintainability)](https://codeclimate.com/github/ipfs-search/ipfs-search/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/1c25261992991d72137c/test_coverage)](https://codeclimate.com/github/ipfs-search/ipfs-search/test_coverage)
[![Go Reference](https://pkg.go.dev/badge/github.com/ipfs-search/ipfs-search.svg)](https://pkg.go.dev/github.com/ipfs-search/ipfs-search)
[![Backers on Open Collective](https://opencollective.com/ipfs-search/backers/badge.svg)](#backers)
[![Sponsors on Open Collective](https://opencollective.com/ipfs-search/sponsors/badge.svg)](#sponsors)

Search engine for the [Interplanetary Filesystem](https://ipfs.io). Sniffs the DHT gossip and indexes file and directory hashes.

Metadata and contents are extracted using [ipfs-tika](https://github.com/dokterbob/ipfs-tika), searching is done using ElasticSearch 7, queueing is done using RabbitMQ. The crawler is implemented in Go, the API and frontend are built using Node.js.

The ipfs-search command consists of two components: the crawler and the sniffer. The sniffer extracts hashes from the gossip between nodes. The crawler extracts data from the hashes and indexes them.

## Docs
A preliminary start at providing a minimal amount of documentation can be found in the [docs](docs/) folder.

## Contact
Please find us on our Freenode/[Riot/Matrix](https://riot.im/app/#/room/#ipfssearch:matrix.org) channel #ipfssearch.

## Snapshots
ipfs-search provides the daily snapshot for all of the indexed data using
[elasticsearch snapshots](https://www.elastic.co/guide/en/elasticsearch/reference/current/modules-snapshots.html).
To learn more about downloading and restoring snapshots, read [docs](docs/snapshots.md)

## Related repo's
* [frontend](https://github.com/ipfs-search/ipfs-search-frontend)
* [metadata API](https://github.com/ipfs-search/ipfs-metadata-api)
* [search API](https://github.com/ipfs-search/ipfs-search-api)
* [deployment](https://github.com/ipfs-search/ipfs-search-deployment)

## Contributors wanted
Building a search engine like this takes a considerable amount of resources (money _and_ TLC).
If you are able to help out with either of them, mail us at info@ipfs-search.com or find us at #ipfssearch on Freenode (or #ipfs-search:chat.weho.st on Matrix).

Please read the Contributing.md file before contributing.

## Roadmap
For discussing and suggesting features, look at the [issues](https://github.com/ipfs-search/ipfs-search/issues).

## Dependencies

* Go 1.13
* Elasticsearch 7.x
* RabbitMQ / AMQP server
* NodeJS 9.x
* IPFS 0.7

## Configuration
Configuration can be done using a YAML configuration file, or by specifying the following environment variables:
* `IPFS_TIKA_URL`
* `IPFS_API_URL`
* `ELASTICSEARCH_URL`
* `AMQP_URL`

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

## Contributors

This project exists thanks to all the people who contribute.
<a href="https://github.com/ipfs-search/ipfs-search/graphs/contributors"><img src="https://opencollective.com/ipfs-search/contributors.svg?width=890&button=false" /></a>


## Backers

Thank you to all our backers! 🙏 [[Become a backer](https://opencollective.com/ipfs-search#backer)]

<a href="https://opencollective.com/ipfs-search#backers" target="_blank"><img src="https://opencollective.com/ipfs-search/backers.svg?width=890"></a>


## Sponsors

<a href="https://nlnet.nl/project/IPFS-search/"><img width="200pt" src="https://nlnet.nl/logo/banner.png"></a> <a href="https://nlnet.nl/project/IPFS-search/"><img width="200pt" src="https://nlnet.nl/image/logos/NGI0_tag.png"></a><br>
ipfs-search is supported by NLNet through the EU's Next Generation Internet (NGI0) programme.

Support this project by becoming a sponsor. Your logo will show up here with a link to your website. [[Become a sponsor](https://opencollective.com/ipfs-search#sponsor)]

<a href="https://opencollective.com/ipfs-search/sponsor/0/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/ipfs-search/sponsor/1/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/ipfs-search/sponsor/2/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/ipfs-search/sponsor/3/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/ipfs-search/sponsor/4/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/ipfs-search/sponsor/5/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/ipfs-search/sponsor/6/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/ipfs-search/sponsor/7/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/ipfs-search/sponsor/8/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/ipfs-search/sponsor/9/website" target="_blank"><img src="https://opencollective.com/ipfs-search/sponsor/9/avatar.svg"></a>


