# API

## REST API
We're using the [querystring query API](https://www.elastic.co/guide/en/opensearch/reference/current/query-dsl-query-string-query.html), allowing filters by field like so: `references.name:epub` or like so `last-seen:>now-1M`.

An up-to-date list of available fields can be found in the index mapping definition for [files](https://github.com/ipfs-search/ipfs-search/blob/master/docs/indices/files.json) and [directories](https://github.com/ipfs-search/ipfs-search/blob/master/docs/indices/directories.json).

In addition, [interactive API documentation](https://api.ipfs-search.com/) is automatically generated from our [OpenAPI spec](https://github.com/ipfs-search/ipfs-search-api/blob/master/openapi-v1.yaml).

## Go documentstaiton
The API of the crawler is fully annotated, documentation is available at [go.dev](https://pkg.go.dev/github.com/ipfs-search/ipfs-search).
