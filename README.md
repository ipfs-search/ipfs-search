# [ipfs-search](http://ipfs-search.com)
[![Build Status](https://travis-ci.org/ipfs-search/ipfs-search.svg?branch=travis)](https://travis-ci.org/ipfs-search/ipfs-search)
[![Maintainability](https://api.codeclimate.com/v1/badges/1c25261992991d72137c/maintainability)](https://codeclimate.com/github/ipfs-search/ipfs-search/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/1c25261992991d72137c/test_coverage)](https://codeclimate.com/github/ipfs-search/ipfs-search/test_coverage)
[![GoDoc](https://godoc.org/github.com/ipfs-search/ipfs-search?status.svg)](https://godoc.org/github.com/ipfs-search/ipfs-search)
[![Backers on Open Collective](https://opencollective.com/ipfs-search/backers/badge.svg)](#backers)
 [![Sponsors on Open Collective](https://opencollective.com/ipfs-search/sponsors/badge.svg)](#sponsors)

Search engine for the [Interplanetary Filesystem](https://ipfs.io). Sniffs the DHT gossip and indexes file and directory hashes.

Metadata and contents are extracted using [ipfs-tika](https://github.com/dokterbob/ipfs-tika), searching is done using ElasticSearch 5, queueing is done using RabbitMQ. The crawler is implemented in Go, the API and frontend are built using Node.js.

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

## Contributors wanted
Building a search engine like this takes a considerable amount of resources (money _and_ TLC).
If you are able to help out with either of them, mail us at info@ipfs-search.com or find us at #ipfssearch on Freenode (or #ipfs-search:chat.weho.st on Matrix).

Please read the Contributing.md file before contributing.

## Roadmap
For discussing and suggesting features, look at the [project planning](https://github.com/ipfs-search/ipfs-search/projects).

## Requirements

* Go 1.12
* Elasticsearch 5.x
* RabbitMQ / AMQP server
* NodeJS 9.x

## Configuration
Configuration can be done using a YAML configuration file, see [`example_config.yml`](example_config.yml).

The following configuration options can be overridden by environment variables:
* `IPFS_TIKA_URL`
* `IPFS_API_URL`
* `ELASTICSEARCH_URL`
* `AMQP_URL`

or by using environment variables.

## Building
```bash
$ go get ./...
$ make
```

## Running

### Docker
The most convenient way to run the crawler is through Docker. Simply run:

```bash
compose up
```

This will start the crawler and all its dependencies but will not (yet) launch the sniffer or search API. Hashes can be queued for crawling manually by running `ipfs-search a <hash>` from within the running container. For example:

```bash
compose exec ipfs-search ipfs-search add QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv
```

### Local setup
Local installation is done using vagrant:

```bash
git clone https://github.com/ipfs-search/ipfs-search.git ipfs-search
cd ipfs-search
vagrant up
```

This starts up the API on port 9615, Elasticsearch on 9200 and RabbitMQ on 15672.

Vagrant setup does not currently start up the frontend.

### Ansible deployment
Automated deployment can be done on any (virtual) Ubuntu 16.04 machine. The full production stack is automated and can be found [here](deployment/).

## Contributors

This project exists thanks to all the people who contribute.
<a href="https://github.com/ipfs-search/ipfs-search/graphs/contributors"><img src="https://opencollective.com/ipfs-search/contributors.svg?width=890&button=false" /></a>


## Backers

Thank you to all our backers! 🙏 [[Become a backer](https://opencollective.com/ipfs-search#backer)]

<a href="https://opencollective.com/ipfs-search#backers" target="_blank"><img src="https://opencollective.com/ipfs-search/backers.svg?width=890"></a>


## Sponsors

<a href="https://nlnet.nl/project/IPFS-search/"><img src="https://nlnet.nl/logo/banner.png"><img src="https://nlnet.nl/image/logos/NGI0_tag.png"></a>
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


