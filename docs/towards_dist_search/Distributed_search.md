# Making ipfs-search distributed

How distributed search could be realized for the IPFS:

-   Provider nodes that wish to participate, parse and index only the files they have added to a dweb (DHT hashes) and that have world file permissions. 
-   This local index is put on an (IPFS) cluster. 
-   A query can use the distributed index.
-   Initial search functionalities are basic [boolean search] https://niverel.tymyrddin.space/en/play/algos/boolean "en:play:algos:boolean") to begin with.
-   Settings functionality anticipates [tuning](https://niverel.tymyrddin.space/en/problems/psearch/design "en:problems:psearch:design").
-   In the future, one can add to the search engine functionality with extensions.
    

Control for users.

### Sketch


The ballon d'essai can consist of:

-   A [distributed index using an IPFS cluster](https://niverel.tymyrddin.space/en/play/stones/upsidedown/ipfs-cluster "en:play:stones:upsidedown:ipfs-cluster")
-   An indexer package with which content providers can index what they provide and add such an index to the distributed index, starting with indexing documents
-   A thin, separate client with which people can query the distributed index and receive results ranked relevant to the query

### [Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/start)


## Overlay networks


An IPFS node can be fingerprinted through the content it stores. An overlay network needs to offer an “anonymous” mode that only enables features known to not leak information.

-   No local discovery.    
-   No transports other than, for example, via Tor (an overlay network consisting of more than seven thousand relays to conceal a user's location and usage from anyone conducting network surveillance or traffic analysis). 
-   Private routing to make the network non-enumerable.
    

And getting any of this wrong could put _some_ people in danger.

### [Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/overlay)


## Parsing


We could code different parsers for each type of file but that is not our main focus at the moment, and because a Python port of the Apache Tika library exists that according to the documentation supports text extraction from over 1500 file formats, we go with that, at least for now. But it is slow, and in the future we may reconsider.

This parser is pointed to the root of a site or a collection, parses its content (thereby creating a corpus) and adds the objects to ipfs, rather than fetching the ipfs hashtable and taking it from there. Again, we wish to focus on the indexing and clustering of a distributed index, not on finding out how to use the ipfs hashtable (for now).

_Code on this page is just a first shot, and to be read as pseudocode snippets._

```
import os, os.path
import ipfsApi
from tika import parser
from multiprocessing import Pool

def tika_parser(file_path):
    # Extract text from document
    content = parser.from_file(file_path)
    if 'content' in content:
        text = content['content']
    else:
        return
    # Convert to string
    text = str(text)
    # Normalisation to utf-8 format
    safe_text = text.encode('utf-8', errors='ignore')
    # Escape any \ issues
    safe_text = str(safe_text).replace('\\', '\\\\').replace('"', '\\"')
    # Add hash (as filename) and content of file to corpus dataframe
    ...

def walkthrough ()
    corpus_root = os.getxxx (path_to_root)
    walk through the directory structure to fetch each file_path and 
        add each encountered object to ipfs (if duplicate, will not be pinned)
        add hash and file_path to paths
    return paths

    pool = Pool()
    pool.map(tika_parser, paths)
    return corpus
```

### Resources


-   [Tika Supported Document Formats](https://tika.apache.org/1.4/formats.html "https://tika.apache.org/1.4/formats.html")    
-   [IPFS API Bindings for Python](https://pypi.org/project/ipfs-api/ "https://pypi.org/project/ipfs-api/")

### [Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/extraction)

## Distributing the index on an IPFS Cluster


-   IPFS does not guarantee redundancy. We can use IPFS clustering.    
-   Only popular indexes will be able to get a decent speed.    
    -   We can run a few web agent type ipfs nodes in a cluster that pin all the indexes. Give these enough bandwidth and we have some basis nodes that can act as mirrors and can also be served via HTTPS (the internet-facing demo version).        
    -   IPFS can replace mirror indexes with IPNS addresses. We will still need reliable hosting for these initial seeders.
        

### Risks
IPFS is still in alpha development. That means there are a lot of (undiscovered) bugs and vulnerabilities and the code is not stable. This could create (security) problems.
    

### Resources


-   [IPFS Cluster Github](https://github.com/ipfs/ipfs-cluster/ "https://github.com/ipfs/ipfs-cluster/")
-   [IPFS Cluster Documentation](https://cluster.ipfs.io/documentation/ "https://cluster.ipfs.io/documentation/") 
-   [IPFS Cluster Architecture overview](https://cluster.ipfs.io/documentation/deployment/architecture/ "https://cluster.ipfs.io/documentation/deployment/architecture/")


### [Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/ipfs-cluster)

## Querying the index


Our intention is to support [boolean queries](https://niverel.tymyrddin.space/en/play/algos/boolean "en:play:algos:boolean") and [phrase queries](https://niverel.tymyrddin.space/en/play/algos/phrase "en:play:algos:phrase").
-   Sanitize the query (stemming all the words, making all letters lowercase, removing punctuation)    
-   Tokenise the query (split into words)    
-   Get term lists from the distributed index, which documents they appear in, and union the lists
    



### Boolean query


For each inverted index from self and received from neighbours:

```
def one_word_query(word, invertedIndex):
	pattern = re.compile('[\W_]+')
	word = pattern.sub(' ',word)
	if word in invertedIndex.keys():
		return [filename for filename in invertedIndex[word].keys()]
	else:
		return []
```

**OR**

**Aggregate lists and union**
```
def free_text_query(string):
	pattern = re.compile('[\W_]+')
	string = pattern.sub(' ',string)
	result = []
	for word in string.split():
		result += one_word_query(word)
	return list(set(result))
```

**AND**

For an AND use an intersection instead of a union to aggregate the results of the single word queries.

### Phrase query


```
def phrase_query(string, invertedIndex):
	pattern = re.compile('[\W_]+')
	string = pattern.sub(' ',string)
	listOfLists, result = [],[]
	for word in string.split():
		listOfLists.append(one_word_query(word))
	setted = set(listOfLists[0]).intersection(*listOfLists)
	for filename in setted:
		temp = []
		for word in string.split():
			temp.append(invertedIndex[word][filename][:])
		for i in range(len(temp)):
			for ind in range(len(temp[i])):
				temp[i][ind] -= i
		if set(temp[0]).intersection(*temp):
			result.append(filename)
	return rankResults(result, string)
```



### [Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/querying)

## Yggdrasil


Yggdrasil is an early-stage implementation of a fully end-to-end encrypted IPv6 network. It is lightweight, self-arranging, supported on multiple platforms and allows pretty much any IPv6-capable application to communicate securely with other Yggdrasil nodes. Yggdrasil does not require IPv6 Internet connectivity - it also works over IPv4.

Looking at it for its clustering and bootstrapping implementation.

### Resources


-   [Yggdrasil Version 0.3.6](https://yggdrasil-network.github.io/2019/08/03/release-v0-3-6.html "https://yggdrasil-network.github.io/2019/08/03/release-v0-3-6.html"), august 2019, first version with API    
-   [Yggdrasil](https://github.com/yggdrasil-network/yggdrasil-go "https://github.com/yggdrasil-network/yggdrasil-go"), Github

### [Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/yggdrasil)

## Testing

-   Scalability indicators    
    -   Number of hashes crawled per second per-peer versus the number of peers        
    -   Number of downloaded bytes per second versus the number of peers        
-   Performance indicators    
    -   Number of hashes crawled per second versus different CPU loads/platforms        
    -   Throughput of a peer versus the number of crawled job queues (to determine the optimal number of crawl job queues) per platform (differentiate using agent attributes).        
-   Node failure    
    -   If automated, this may require adding data entry points in the API that are only used for testing.        
    -   Add test data, check that it has been added and has propagated throughout the neighbourhood.
    -   Take an agent offline (check that it has gone down and is inaccessible) and verify that all the data appears to be working.        
    -   Pull data manually from each data store (check there are no errors as a result) on the agent, and verify that the data is still retrievable from the system.
    -   Bring the downed node back online. The data that belongs on this node (hopefully) begins to flow back into the node.        
    -   After a while, pull the data from the agent to check that data that was sent to its neighbours when it was down is stored correctly.        
-   Predictive analysis 
    -   Test for false negatives and false positives of the various classifiers with unlabelled traffic data
        

### Resources


-   [Grid'5000](https://www.fed4fire.eu/testbeds/grid5000/ "https://www.fed4fire.eu/testbeds/grid5000/")


### [Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/testing)