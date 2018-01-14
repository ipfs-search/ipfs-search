# How to reindex

1. Stop crawler.
```
$ systemctl stop ipfs-crawler
```

1. Create snapshot to allow for rollback:
```
PUT
/_snapshot/ipfs/snapshot_v<old>
```

2. Create new index:
```
PUT /ipfs_v<new>
<<< index-json >>>
```

3. Reindex old to new:
```
POST /_reindex
{
  "source": {
    "index": "ipfs_v<old>"
  },
  "dest": {
    "index": "ipfs_v<new>"
  }
}
```
(Go fetch some coffee for this one.)

4. Remove old alias, create new alias:
```
POST /_aliases
{
    "actions" : [
        { "remove" : { "index" : "ipfs_v<old>", "alias" : "ipfs" } },
        { "add" : { "index" : "ipfs_v<new>", "alias" : "ipfs" } }
    ]
}
```

5. Restart crawler:
```
$ systemctl start ipfs-crawler
```

5. Remove old index (after verifying everything is ok):
```
DELETE /ipv_v<old>
```
