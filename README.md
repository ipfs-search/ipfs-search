# [ipfs-search](http://ipfs-search.com)
Search engine for the [Interplanetary Filesystem](https://ipfs.io). Sniffs the DHT gossip and indexes file and directory hashes.

Metadata and contents are extracted using [ipfs-tika](https://github.com/dokterbob/ipfs-tika), searching is done using ElasticSearch 5, queueing is done using RabbitMQ. The crawler is implemented in Go, the API and frontend are built using Node.js.

## Maintainer requested
So terribly sorry, but hosting a search engine like this takes a considerable amount of resources (money _and_ TLC).

As this moment, the founders of ipfs-search, moved on to bigger and better things and had to cut hosting.</p>

If you are able to help out with either of them, >mail us at info@ipfs-search.com or find us at #ipfssearch on Freenode (or #ipfs-search:chat.weho.st on Matrix).

## Roadmap
For discussing and suggesting features, look at the [project planning](https://github.com/ipfs-search/ipfs-search/projects).

## Vagrant
```bash
$ vagrant up
```
The search engine should now listen on port 8881 of your local machine, with the API directly exposed on 9615, ES on 9200 and RabbitMQ on 15672.

## Manual provisioning
```bash
$ ansible-playbook provisioning/bootstrap.yml --user root --ask-pass
$ ansible-playbook provisioning/ipfs-search.yml
```
