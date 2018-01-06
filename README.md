# [ipfs-search](http://ipfs-search.com)
Search engine for the [Interplanetary Filesystem](https://ipfs.io). Sniffs the DHT gossip and indexes file and directory hashes.

Metadata and contents are extracted using [ipfs-tika](https://github.com/dokterbob/ipfs-tika), searching is done using ElasticSearch 5, queueing is done using RabbitMQ. The crawler is implemented in Go, the API and frontend are built using Node.js.

## Contributors wanted
Building a search engine like this takes a considerable amount of resources (money _and_ TLC).
If you are able to help out with either of them, mail us at info@ipfs-search.com or find us at #ipfssearch on Freenode (or #ipfs-search:chat.weho.st on Matrix).

## Roadmap
For discussing and suggesting features, look at the [project planning](https://github.com/ipfs-search/ipfs-search/projects).

## Running
First of all, make sure Ansible 2.2 is installed:

```bash
$ pip2 install 'ansible<2.3'
```

### Local setup
Local installation is done using vagrant:

```bash
git clone https://github.com/ipfs-search/ipfs-search.git $GOPATH/src/github.com/ipfs-search/ipfs-search
cd $GOPATH/src/github.com/ipfs-search/ipfs-search
vagrant up
```

This starts up the API on port 9615, Elasticsearch on 9200 and RabbitMQ on 15672.

Vagrant setup does not currently start up the frontend.

### Manual provisioning
```bash
$ ansible-playbook provisioning/bootstrap.yml --user root --ask-pass
$ ansible-playbook provisioning/backend.yml
$ ansible-playbook provisioning/frontend.yml
```
