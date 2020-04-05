# ipfs-search documentation

## Architecture

ipfs-search consists of the following components:
* Sniffer
* Queue
* Crawler
* Metadata extractor
* Search backend
* API
* Frontend

### Sniffer
The sniffer listens to gossip between our IPFS node and others and adds hashes for which a provider is offered to the `hashes` queue, filtering for (currently) unparseable data and items recently updated.

### Queue: RabbitMQ
RabbitMQ holds a `files` and a `hashes` queue with items to be crawled, in a soon-to-be well-defined JSON-format.

### Crawler: ipfs-search
#### Hashes (directories or files)
The crawler takes items of the `hashes` queue and attempts to list the items using the IPFS RPC API. This will tell it whether the item is a file, a directory or some other type.

In case it's a directory, the directory listing will be added and the referred items will be added to the `hashes` queue in case they are directories and to the `files` queue in case they are files.

In the case the crawled item is a file, it will be added to the `files` queue and no further action is taken.

#### Files (only files)
Jobs taken from the `files` queue are guaranteed to be files, metadata extraction and content type detection will be attempted by IPFS TIKA.

#### Updating items
All indexed items will be initially given a `first-seen` field and, when seen again, will have their `last-seen` field set or updated.

#### References
When an item is referred to from a directory, i.e. when it's found to be a directory item in the hashes queue, it's referenced name and parent directory will be added to the list of references for that given item. This will happen both for new as well as existing items.

### Metadata extractor: ipfs-tika
IPFS-TIKA uses the local IPFS gateway to fetch a (named) IPFS resource and streams the resulting data into an Apache TIKA metadata extractor.

It currently extracts body text up to a certain limit, links and any available metadata. In the future we hope to detect the language as well.

### Search backend: Elasticsearch
Any crawled items will be stored in Elasticsearch, which has a custom mapping defined to prevent the many returned metadata fields from all being indexed (for obvious efficiency reasons).

It has been found that it is necessary to regularly update the index to circumvent occasional problems with indexing, performance, queries or other factors.

### API
The API provides a layer on top of the search backend, providing filtered output and a limited query functionality, as well as reformatting the resulting items.

In the near future we hope to provide an endpoint for adding new items to the crawl queue as well.

### Frontend
The frontend is nothing more than a static front to the search API.
