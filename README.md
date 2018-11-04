# [ipfs-search](http://ipfs-search.com)
[![Build Status](https://travis-ci.org/ipfs-search/ipfs-search.svg?branch=travis)](https://travis-ci.org/ipfs-search/ipfs-search)
[![Maintainability](https://api.codeclimate.com/v1/badges/1c25261992991d72137c/maintainability)](https://codeclimate.com/github/ipfs-search/ipfs-search/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/1c25261992991d72137c/test_coverage)](https://codeclimate.com/github/ipfs-search/ipfs-search/test_coverage)
[![GoDoc](https://godoc.org/github.com/ipfs-search/ipfs-search?status.svg)](https://godoc.org/github.com/ipfs-search/ipfs-search)

Search engine for the [Interplanetary Filesystem](https://ipfs.io). Sniffs the DHT gossip and indexes file and directory hashes.

Metadata and contents are extracted using [ipfs-tika](https://github.com/dokterbob/ipfs-tika), searching is done using ElasticSearch 5, queueing is done using RabbitMQ. The crawler is implemented in Go, the API and frontend are built using Node.js.

## Docs
A preliminary start at providing a minimal amount of documentation can be found in the [docs](docs/) folder.

## Contributors wanted
Building a search engine like this takes a considerable amount of resources (money _and_ TLC).
If you are able to help out with either of them, mail us at info@ipfs-search.com or find us at #ipfssearch on Freenode (or #ipfs-search:chat.weho.st on Matrix).

## Roadmap
For discussing and suggesting features, look at the [project planning](https://github.com/ipfs-search/ipfs-search/projects).

## Requirements

* Go 1.11
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
