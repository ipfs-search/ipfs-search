# [ipfs-search.com](https://ipfs-search.com)
[![pipeline status](https://gitlab.com/ipfs-search.com/ipfs-search/badges/master/pipeline.svg)](https://gitlab.com/ipfs-search.com/ipfs-search/-/commits/master)
[![Maintainability](https://api.codeclimate.com/v1/badges/1c25261992991d72137c/maintainability)](https://codeclimate.com/github/ipfs-search/ipfs-search/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/1c25261992991d72137c/test_coverage)](https://codeclimate.com/github/ipfs-search/ipfs-search/test_coverage)
[![Documentation Status](https://readthedocs.org/projects/ipfs-search/badge/?version=latest)](https://ipfs-search.readthedocs.io/en/latest/?badge=latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/ipfs-search/ipfs-search.svg)](https://pkg.go.dev/github.com/ipfs-search/ipfs-search)
[![Backers on Open Collective](https://opencollective.com/ipfs-search/backers/badge.svg)](#backers)
[![Sponsors on Open Collective](https://opencollective.com/ipfs-search/sponsors/badge.svg)](#sponsors)

Search engine for the [Interplanetary Filesystem](https://ipfs.io). Sniffs the DHT gossip and indexes file and directory hashes.

Metadata and contents are extracted using [ipfs-tika](https://github.com/ipfs-search/ipfs-tika), searching is done using OpenSearch, queueing is done using RabbitMQ. The crawler is implemented in Go, the API and frontend are built using Node.js.

The ipfs-search command consists of two components: the crawler and the sniffer. The sniffer extracts hashes from the gossip between nodes. The crawler extracts data from the hashes and indexes them.

## Docs
Documentation is hosted on [Read the Docs](https://ipfs-search.readthedocs.io/en/latest/), based on files contained in the [docs](https://github.com/ipfs-search/ipfs-search/tree/master/docs) folder. In addition, there's extensive [Go docs](https://pkg.go.dev/github.com/ipfs-search/ipfs-search) for the internal API as well as [SwaggerHub OpenAPI documentation](https://app.swaggerhub.com/apis-docs/ipfs-search/ipfs-search/) for the REST API.

## Contact
Please find us on our Freenode/[Riot/Matrix](https://riot.im/app/#/room/#ipfs-search:matrix.org) channel [#ipfs-search:matrix.org](https://matrix.to/#/#ipfs-search:matrix.org).

## Snapshots
ipfs-search provides the daily snapshot for all of the indexed data using
[snapshots](https://opensearch.org/docs/latest/opensearch/rest-api/snapshots/index/).
To learn more about downloading and restoring snapshots please refer to the [relevant section](https://ipfs-search.readthedocs.io/en/latest/snapshots.html) in our documentation.

## Related repo's
* [frontend](https://github.com/ipfs-search/ipfs-search-frontend)
* [search API](https://github.com/ipfs-search/ipfs-search-api)
* [deployment](https://github.com/ipfs-search/ipfs-search-deployment)
* [nsfw-server](https://github.com/ipfs-search/nsfw-server)
* [ipfs-tika](https://github.com/ipfs-search/ipfs-tika)

## Contributors wanted
Building a search engine like this takes a considerable amount of resources (money _and_ TLC).
If you are able to help out with either of them, do reach out (see the contact section in this file).

Please read the Contributing.md file before contributing.

## Roadmap
For discussing and suggesting features, look at the [issues](https://github.com/ipfs-search/ipfs-search/issues).

## External dependencies

* Go 1.19
* OpenSearch 2.3.x
* RabbitMQ / AMQP server
* NodeJS 9.x
* IPFS 0.7
* Redis

## Internal dependencies

* [nsfw-server](https://github.com/ipfs-search/nsfw-server)
* [ipfs-tika](https://github.com/ipfs-search/ipfs-tika)

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

<a href="https://nlnet.nl/project/IPFS-search/"><img width="200pt" src="https://nlnet.nl/logo/banner.png"></a> <a href="https://nlnet.nl/project/IPFS-search/"><img width="200pt" src="https://nlnet.nl/image/logos/NGI0_tag.png"></a>
<br>
ipfs-search is supported by NLNet through the EU's Next Generation Internet (NGI0) programme.

<a href="https://redpencil.io/projects/"><img width="270pt" src="https://raw.githubusercontent.com/redpencilio/frontend-redpencil.io/327318b84ffb396d8af6776f19b9f36212596082/public/assets/vector/rpio-logo.svg"> </a><br>
RedPencil is supporting the hosting of ipfs-search.com.

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


