# ipfs-search
IPFS search engine POC, based on initial work in ipfs-crawler

It's using RabbitMQ for queueing and indexes using Elasticsearch.

## Provisioning
```bash
$ ansible-playbook provisioning/bootstrap.yml --user root --ask-pass
$ ansible-playbook provisioning/ipfs-search.yml
```
