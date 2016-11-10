# [ipfs-search](http://ipfs-search.com)
Search engine for the [Interplanetary Filesystem](https://ipfs.io). Sniffs the DHT gossip and indexes file and directory hashes.

Metadata and contents are extracted using [ipfs-tika](https://github.com/dokterbob/ipfs-tika), searching is done using ElasticSearch 5, queueing is done using RabbitMQ. The crawler is implemented in Go, the API frontend in Node.js.

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
