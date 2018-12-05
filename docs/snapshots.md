# ipfs-search snapshots
ipfs-search makes daily [elasticsearch snapshots](https://www.elastic.co/guide/en/elasticsearch/reference/current/modules-snapshots.html) of the indexed data.

As soon as [ipfs/go-ipfs#5815](https://github.com/ipfs/go-ipfs/issues/5815) is solved, they will be automatically published on [IPFS](https://ipfs.io/). Inthe meantime, you may contact the team on IRC/Matrix (#ipfssearch on Freenode or #ipfs-search:chat.weho.st on Matrix) for less recent 'manual' shares of the snaphshot.

## Pinning
To pin the snapshots:
`ipfs pin add $hash`

To automatically resume the pinning when interrupted you can use the following command:
```
while [ 1 ]; do ipfs pin add $hash; sleep 60; done
```

## Restoring 
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

For further information use [elasticsearch snapshots](https://www.elastic.co/guide/en/elasticsearch/reference/current/modules-snapshots.html).

## License
[CC-BY-SA 4.0](https://github.com/idleberg/Creative-Commons-Markdown/blob/master/4.0/by-sa.markdown)
