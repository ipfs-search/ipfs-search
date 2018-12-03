## ipfs-search snapshots
ipfs-search provides the daily snapshot for all of the indexed data using 
[elasticsearch snapshots](https://www.elastic.co/guide/en/elasticsearch/reference/current/modules-snapshots.html).
The snapshots are available on [IPFS](https://ipfs.io/) with the following hash: `Qmc3RxfyZTPf7omWN1XxDkaZhp93ukfLSY14CTC8n1v5Hv`

To pin the snapshots:
`ipfs pin add Qmc3RxfyZTPf7omWN1XxDkaZhp93ukfLSY14CTC8n1v5Hv`

For now its size is about 325GB, So to automatically resume the pinning when interrupted you can use the following command:
```
while [ 1 ]; do ipfs pin add Qmc3RxfyZTPf7omWN1XxDkaZhp93ukfLSY14CTC8n1v5Hv; sleep 60; done
```

## Restoring the snapshot
Download the `ipfs-search` snapshot with  `ipfs add pin` then mount it using 
[FUSE](https://github.com/ipfs/go-ipfs/blob/master/docs/fuse.md) or use `ipfs get`\

Add the path of the downloaded snapshot to the configuration file:
1. Open `elasticsearch.yml`
2. Add: `path.repo: ["path/to/snapshot"]`
3. Run the elasticsearch
 Now use the following command to register the repository with any name (for example: ipfs_search)
```
curl -X PUT "localhost:9200/_snapshot/ipfs_search" -H 'Content-Type: application/json' -d'
{
    "type": "fs",
    "settings": {
        "location": "path/to/snapshot",
        "compress": true
    }
}'
```
To list all of the available snapshots:searchsearch\
`curl -X GET "localhost:9200/_snapshot/ipfs_search/_all?pretty"`
 To show a specific snapshot (for example the snapshot snapshot_181025-0316)\
`curl -X GET "localhost:9200/_snapshot/ipfs_search/snapshot_181025-0316?pretty"`
 Restore the specified snapshot with this command\
`curl -X POST "localhost:9200/_snapshot/elastic_search/snapshot_181025-0316/_restore?wait_for_completion=true"`
