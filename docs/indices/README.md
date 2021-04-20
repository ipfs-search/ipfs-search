# Indices

## Elasticsearch index mapping and settings
To be used in the [Create index API](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html):
```
$ curl -d <index.json> -X PUT http://localhost:9200/<index-name>
```

* [Files](https://github.com/ipfs-search/ipfs-search/blob/master/docs/indices/files.json)
* [Directories](https://github.com/ipfs-search/ipfs-search/blob/master/docs/indices/directories.json)
* [Invalids](https://github.com/ipfs-search/ipfs-search/blob/master/docs/indices/invalids.json)
* [Partials](https://github.com/ipfs-search/ipfs-search/blob/master/docs/indices/partials.json)

## Example entries

Examples of real-life crawled content are available for a [file](https://github.com/ipfs-search/ipfs-search/blob/master/docs/example_file.json) and a [directory](https://github.com/ipfs-search/ipfs-search/blob/master/docs/example_directory.json).

## Reindexing
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
