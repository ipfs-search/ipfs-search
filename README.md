# ipfs-search
IPFS search engine POC, based on initial work in ipfs-crawler

It's using RabbitMQ for queueing and indexes using Elasticsearch.

## Bootstrap
$ ansible-playbook -i provisioning/hosts provisioning/bootstrap.yml --user root --ask-pass
