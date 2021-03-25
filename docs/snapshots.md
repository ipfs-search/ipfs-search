# ipfs-search snapshots
ipfs-search makes daily [elasticsearch snapshots](https://www.elastic.co/guide/en/elasticsearch/reference/5.6/modules-snapshots.html) of the indexed data.

We are currently experimenting with automated publishing of these daily snapshots over IPFS. This should allow anyone to inspect our index and/or to fork or mirror our service.
As of the time of writing (April 5, 2020) the full index is about 425 GB.

## Cluster
We are running an [ipfs-cluster](https://cluster.ipfs.io/), automating the process of pinning the latest updates. The easiest way to do this, is throuhg [ipfs-cluster-follow](https://cluster.ipfs.io/documentation/collaborative/joining/):

1. [Run a local IPFS Node](https://docs.ipfs.io/how-to/command-line-quick-start/).
2. [Download](https://dist.ipfs.io/#ipfs-cluster-follow) for your platform and extract the archive.
3. Run: `ipfs-cluster-follow ipfs-search run --init cluster.ipfs-search.com`

Now `ipfs-cluster-follow` should download the cluster configuration, connect to other nodes and start pinning the latest snapshot, automatically updating every night.

## Manual pinning
The daily snapshots, for now, are published to: https://gateway.ipfs.io/ipns/12D3KooWKDDboo2aQzFxpHB7BXUUXudMr81ccC4d28eQPAfrgWQi

To pin the snapshots:
`ipfs pin add /ipns/12D3KooWKDDboo2aQzFxpHB7BXUUXudMr81ccC4d28eQPAfrgWQi`

To automatically resume the pinning when interrupted you can use the following command:
```
while [ 1 ]; do ipfs pin add --progress /ipns/12D3KooWKDDboo2aQzFxpHB7BXUUXudMr81ccC4d28eQPAfrgWQi; sleep 60; done
```

## Restoring
It should be possible to load the snapshots directly through a (local) IPFS gateway into Elasticsearch, although this has not yet been tested and it is most certainly advisable to pin the dataset as per the instructions above.

In order to load the snapshots, first make sure you're running a compatible (or equal) version of Elasticsearch and that there is enough disk space available (twice the current size of the index, so ~ 1TB as of the time of writing).

The steps are as follows:

1. [Run an IPFS Node](https://docs.ipfs.io/introduction/usage/)
2. Pin the index snapshot (as per instructions above)
3. [Run Elasticsearch 5.x](https://www.elastic.co/guide/en/elasticsearch/reference/5.6/install-elasticsearch.html)
4. Register the local IPFS gateway as a [readonly Elasticsearch snapshot repository](https://www.elastic.co/guide/en/elasticsearch/reference/5.6/modules-snapshots.html#_read_only_url_repository) through the URL `http://localhost:8080/ipns/12D3KooWKDDboo2aQzFxpHB7BXUUXudMr81ccC4d28eQPAfrgWQi/backup/` (assuming your local IPFS gateway is running on `localhost:8080`)
5. List available snapshots with `curl -X GET "localhost:9200/_snapshot/<repo_name>/_all?pretty"
`, testing your prior configuration
6. [Restore the snapshot](https://www.elastic.co/guide/en/elasticsearch/reference/5.6/modules-snapshots.html#_restore) of your choice: `curl -X POST "localhost:9200/_snapshot/<repo_name>/<snapshot_id>/_restore?pretty"
`
7. Wait... now...
8. Query away! You now run an exact copy of the ipfs-search.com index!

## Snapshot data license
[CC-BY-SA 4.0](https://github.com/idleberg/Creative-Commons-Markdown/blob/master/4.0/by-sa.markdown)
