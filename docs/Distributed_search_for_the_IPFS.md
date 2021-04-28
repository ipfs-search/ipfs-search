# Making ipfs-search.com a distributed search engine

Our ambition is to make ipfs-search into a truly distributed search engine. We are grateful To [Nina](https://niverel.tymyrddin.space/doku.php?id=en/start) for the work she put in to researching distributed search on behalf of ipfs-search. The following is adapted from the results, first published on Nina's homepage.

## What does it mean for search to be distributed?

### Distributed search engines


Sadly most that existed have gone defunct, leaving [seekstorm (former FAROO](https://seekstorm.com/ "https://seekstorm.com/") (Propietary, English), [Seeks](https://beniz.github.io/seeks/ "https://beniz.github.io/seeks/") (Open Source, English) and [YaCy](https://yacy.net/ "https://yacy.net/") (Fully distributed).

### IPFS search engines

**Noetic**
Noetic searches IPFS content as indexed by conventional search engines. Noetic has had no commits over the past three years.    

**Trinity**
Trinity is a more promising but currently unfunctional effort to make a search engine and has had no commits for a year.    

**ipfssearch.xyz**
ipfsearch.xyz searches decentrally but requires a pre-compiled database.    

**dweb.page**
dweb.page uses IOTA, a permissionless distributed ledger allowing every single user to run their own search engine.

[Source](https://niverel.tymyrddin.space/en/play/bazaar/engines/distributed)




## The path to distributed search

Learning paths of technology evolutions which we travelled for getting transforming ideas for dweb search.

#### Distributed technology
* Blockchain
* Hashgraph
* DAG
* Holochain

#### Assistive technologies
* Zero knowledge proofs
* Smart contracts
* Voting
* Gossip about gossip

#### Current implementations
* IPFS
* BTFS
* FileCoin
* The usual peer crawling

#### Internet-facing demo
* Current IPFS Search architecture
* Good enough indexing

#### Technologies for distributed search
* IPFS
* Overlay networks
* Parsing
* Distributing the index on an IPFS  Cluster
* Querying the index
* Yggdrasil
* Testing

#### Moonshot for realizing distributed search


* Provider nodes that wish to participate, parse and index only the files they have added to a dweb (DHT hashes) and that have world file permissions.
* This local index is put on an (IPFS) cluster.
* A query can use the distributed index.
* Initial search functionalities are basic boolean search, to begin with.
* Settings functionality anticipates tuning.
* In the future, one can add to the search engine functionality with extensions.

[Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/start)

# Distributed technology

Technologies that can be considered 'distributed'.

## Blockchain 

Stuart Haber and W. Scott Stornetta already envisioned a cryptographically secured chain of blocks whereby no one could tamper with timestamps of documents in 1991. In 1992, they upgraded their system to use Merkle trees, increasing efficiency and enabling the collection of more documents on a single block. Satoshi Nakamoto, a person or a group of people, developed the first application of the digital ledger technology in 2008, BitCoin.

* Data is structured in blocks in order of transactions that are validated by miners.
* Each block produces a unique hash that identifies the transaction. If one attempts to alter the details of the transaction, a different hash will be generated. This can be evidence of a corrupted and invalid transaction.
* Transactions are published on a public ledger to which every node has access (transparency). The distributed nature of the public ledger makes it even more difficult for parties to tamper with information.
* Miners can postpone or even cancel a transaction.
* Traditional Blockchains rely on Proof of Work. These need many computations and as a result, the number of transactions per second is relatively low.
    * A transaction has to validate numerous transactions before being valid.
    * As blocks in blockchain multiply, it becomes increasingly difficult in terms of computations to achieve new blocks and mining becomes more power-intensive (expensive).

#### Use cases
Cryptocurrencies

### Resources 
[Bitcoin: A Peer-to-Peer Electronic Cash System](https://bitcoin.org/bitcoin.pdf)

[Source](https://niverel.tymyrddin.space/en/play/stones/dweb/blockchain)

## Hashgraph

The hashgraph algorithm was invented by Leemon Baird for achieving consensus quickly, fairly, efficiently, and securely.

-   Hashgraph achieves transaction success solely via consensus timestamping to make sure that transactions on the network agree with each node on the platform.    
-   On a Hashgraph network nodes do not have to validate transactions by _Proof of Work_ or _Proof of Stake_. Consensus is built with the _Gossip about Gossip_ and _Virtual Voting_ techniques instead, increasing the number of transactions per second.    
-   And consensus timestamping avoids the [Blockchain](https://niverel.tymyrddin.space/en/play/stones/dweb/blockchain "en:play:stones:dweb:blockchain") issues of cancelling transactions or by putting them on future blocks.    
-   These consensus techniques also facilitate fairness.    
-   Developers do not need a license but need the platform coin instead. API calls cost a micro-payment to the company.
    

#### Use cases


All use cases where trust is immutable and incorruptible, for example:

-   Cryptocurrency as a service for support for native micropayments    
-   Micro-storage in the form of a distributed file service that apps can use
-   Contracts    
-   Bank transfers    
-   Credential verification
    

### Resources
-   [Swirlds](https://www.swirlds.com/ "https://www.swirlds.com/")    
-   [Hedera Hashgraph](https://www.hedera.com/ "https://www.hedera.com/")

[Source](https://niverel.tymyrddin.space/en/play/stones/dweb/hashgraph)


## DAG

A DAG is a type of distributed ledger technology that relies on consensus algorithms. To prevail, transactions require majority support within the network. As a result, there is more cooperation and teamwork and nodes have equal rights. Such networks stick to the original goal of Distributed Ledger Technology, to democratise the internet economy.

-   No blocks. No chain. DAG is a structure that is connected like a mesh.
-   It connects current data transactions with previous ones.    
-   With nodes having equal rights, nodes do not have to refer to another node.    
-   A consensus-based system where nodes decide what happens to give a semblance of democracy as compared to platforms that go through a central command.    
-   For a transaction to succeed, it has to validate only two of the previous transactions.    
-   Transactions in DAGs adds throughput as many more validations happen.
    

#### Use cases
-   Cryptocurrencies    
-   Economic infrastructure for data sharing on the Internet of Things    
-   Remote Patient Monitoring    
-   Decentralised Peer-to-Peer energy trading
    

### Resources
-   [OByte](https://obyte.org/ "https://obyte.org/")    
-   [IoTA Use Cases](https://files.iota.org/comms/IOTA_Use_Cases.pdf "https://files.iota.org/comms/IOTA_Use_Cases.pdf")

[Source](https://niverel.tymyrddin.space/en/play/stones/dweb/dag)




## Holochain


Holochain combines Blockchain, BitTorrent and Github ideas for creating a technology that distributes among nodes to avoid any instance of centralised control of the flow of data. A truly revolutionary technology. Blockchain seeks to decentralise transactions such that people can interact directly without the need for a middle party. Holochain distributes the interactions.

-   MetaCurrency is the root, their next-generation Operating System is called Ceptr, a Holochain with Holo being their first real-world application system.    
-   Each node runs on a chain of its own, giving nodes the freedom to operate autonomously. Not all data needs to be shared with everyone. If two people wish to transfer value, and they agree, others do not have to know about it.    
-   There is no need for miners. Transaction fees are almost non-existent. There is no tokenization on the platform. Smart contracts rule.   
-   Users can store data using keys in the Holochain distributed hash table (DHT), but the data stays in actual locations “distributed” in various locations across the globe, relieving the network and improving scalability.    
-   Holochain creates a network composed of various distributed ledger technology networks.    
-   A developer will only need confirmation from the single-chain that makes up the whole DLT network.    
-   Holochain liberates us from corporate control over our choices and information.    
    -   Scalable distributed apps with data integrity        
    -   p2p networks with validating distributed hash tables        
    -   A technology inspired by nature
        

#### Use cases


Systems where not all parties need to participate:

-   Social networks    
-   Chat programs    
-   p2p platforms    
-   Shared document updates
    

### Resources


-   [r/holochain: Distributed Computing and Applications](https://www.reddit.com/r/holochain/ "https://www.reddit.com/r/holochain/")    
-   [Holochain projects](http://holochainprojects.com/ "http://holochainprojects.com/")    
-   [Decentralising the web: The key takeaways](https://www.computing.co.uk/ctg/news/3036546/decentralising-the-web-the-key-takeaways "https://www.computing.co.uk/ctg/news/3036546/decentralising-the-web-the-key-takeaways"), 2018    
-   [Holochains for Distributed Data Integrity](http://ceptr.org/projects/holochain "http://ceptr.org/projects/holochain")   
-   [Holochain](https://holochain.org/ "https://holochain.org/")

[Source](https://niverel.tymyrddin.space/en/play/stones/dweb/holochain)

# Assistive technologies


A truly decentralised web will require the network to provide privacy and trust **by design**. This requires algorithms that allow for trustless management. Zero knowledge proofs and/or cross-validation enables nodes to verify the existence and validity of exchanges. The challenge is to maintain a distributed consensus, without actually being able to see or make public any of the transaction details, guaranteeing privacy.

## Zero Knowledge proofs


A [blockchain](https://niverel.tymyrddin.space/en/play/stones/dweb/blockchain "en:play:stones:dweb: blockchain") is a data structure, a linear transaction log, replicated by devices whose users are rewarded for logging new transactions.

-   A change in any block invalidates every block after it, which means that an adversary can not tamper with historical transactions.    
-   A user only gets rewarded if they are working on the same chain as everyone else, so each participant has an incentive to go with the consensus. The result is a shared definitive historical record.
    

Devil's advocate:

-   It is not truly “trustless”, because most of its users are trusting the software, instead of trusting other people.    
-   Blockchain systems do not magically make the data in them accurate or the people entering the data trustworthy, they merely enable any user to audit whether it has been tampered with.    
-   Making it possible to download the blockchain from a broadcast node and decrypt the Merkle root from the Linux command line to independently verify transactions will only be used by the tech-savvy.
    

### Millionaire’s Problem


In Yao’s Millionaire’s Problem, two millionaires want to find out if they have the same amount of money without disclosing the exact amount. This problem is analogous to a more general problem where there are two numbers `a` and `b` and the goal is to determine whether the inequality `a ≥ b` is true or false without revealing the actual values of `a` and `b`. Problems like these can be used in Zero Knowledge Protocols or Zero Knowledge Password Proofs (ZKPs). The latter, Zero Knowledge Password Proof, is a way of doing authentication where no passwords are exchanged, which means they cannot be stolen. The first term, Zero Knowledge Protocol, has appeared more frequently within blockchain circles.

## Zero Knowledge Protocols


A protocol has to have the following properties to be considered Zero Knowledge Proof:

-   Completeness — If the statement is true, an honest verifier (that is, one following the protocol properly) will be convinced of this fact by an honest prover.    
-   Soundness — If the statement is false, no cheating prover can convince an honest verifier that it is true, except with some small probability.    
-   Zero-Knowledge — If the statement is true, no verifier learns anything other than the fact that the statement is true. In other words, just knowing the statement (not the secret) is sufficient to imagine a scenario showing that the prover knows the secret. This is formalized by showing that every verifier has some simulator that, given only the statement to be proved (and no access to the prover), can produce a transcript that “looks like” an interaction between the honest prover and the verifier in question.
    

Zero-knowledge proofs are probabilistic “proofs” (useful for “good enough” algorithms) rather than deterministic proofs, and techniques exist to decrease the soundness error to negligibly small values.

### Resources


-   September 2017, [the first ZKP was conducted on the Byzantium fork of Ethereum](https://cointelegraph.com/news/ethereum-upgrade-byzantium-is-live-verifies-first-zk-snark-proof "https://cointelegraph.com/news/ethereum-upgrade-byzantium-is-live-verifies-first-zk-snark-proof"), and perhaps we can adapt it for our purpose.    
-   [The knowledge complexity of interactive proof systems](https://people.csail.mit.edu/silvio/Selected%20Scientific%20Papers/Proof%20Systems/The_Knowledge_Complexity_Of_Interactive_Proof_Systems.pdf "https://people.csail.mit.edu/silvio/Selected%20Scientific%20Papers/Proof%20Systems/The_Knowledge_Complexity_Of_Interactive_Proof_Systems.pdf"), Shafi Goldwasser, Silvio Micali, Charles Rackoff, February 1989    
-   [How to Explain Zero-Knowledge Protocols to Your Children](http://pages.cs.wisc.edu/~mkowalcz/628.pdf "http://pages.cs.wisc.edu/~mkowalcz/628.pdf"), Jean-Jacques Quisquater, Louis Guillou, Thomas Berson, 1990    
-   [A Survey of Noninteractive Zero-KnowledgeProof System and Its Applications](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC4032740/pdf/TSWJ2014-560484.pdf "https://www.ncbi.nlm.nih.gov/pmc/articles/PMC4032740/pdf/TSWJ2014-560484.pdf"), Huixin Wu and Feng Wang, May 2014    
-   [A Mind-Bending Cryptographic Trick Promises to Take Blockchains Mainstream](https://www.technologyreview.com/s/609448/a-mind-bending-cryptographic-trick-promises-to-take-blockchains-mainstream/ "https://www.technologyreview.com/s/609448/a-mind-bending-cryptographic-trick-promises-to-take-blockchains-mainstream/"), Mike Orcutt, November 2017

[Source](https://niverel.tymyrddin.space/en/play/stones/assistive/zkp
)



## Smart contracts


A smart contract is a computer protocol intended to digitally facilitate, verify, or enforce the negotiation or performance of a contract. Smart contracts allow the performance of credible transactions between disparate, anonymous parties without third parties (central authority, legal system, or external enforcement). These transactions are trackable and irreversible and serve to increase the integrity of the ledger.

### Resources


-   [Research Directions in Blockchain Data Management and Analytics](https://openproceedings.org/2018/conf/edbt/paper-227.pdf "https://openproceedings.org/2018/conf/edbt/paper-227.pdf"), Hoang Tam Vo, Ashish Kundu, Mukesh Mohania, 2018

[Source](https://niverel.tymyrddin.space/en/play/stones/assistive/contracts)



## Voting


Voting models for consensus have certain security properties (they can be asynchronous Byzantine Fault Tolerant and achieve consensus even when some nodes are malicious and some messages are significantly delayed), but generally require that a large number of messages be sent amongst nodes to get to consensus. Hashgraph has come up with a voting mechanism where this would not be necessary …

### Resources
[I Want Your Vote! (Oh Wait I Already Know It)](https://medium.com/hashgraph/i-want-your-vote-oh-wait-i-already-know-it-e1faa50b31ad "https://medium.com/hashgraph/i-want-your-vote-oh-wait-i-already-know-it-e1faa50b31ad"), Paul Madsen, 2017

[Source](https://niverel.tymyrddin.space/en/play/stones/assistive/voting)



## Gossip about gossip


Gossip about gossip, based on the gossip communication protocol, enables nodes to efficiently and rapidly exchange data with other nodes at internet scale. There are currently three known libraries that implement a gossip algorithm to discover nodes in a peer-to-peer network: Apache Gossip, gossip-python and Smudge.

### Resources


[Apache Gossip](https://github.com/apache/incubator-gossip "https://github.com/apache/incubator-gossip") (UDP, Java)    
[gossip-python](https://github.com/thomai/gossip-python "https://github.com/thomai/gossip-python") (TCP)    
[Smudge](https://github.com/clockworksoul/smudge "https://github.com/clockworksoul/smudge") (UDP, Go)

[Source](https://niverel.tymyrddin.space/en/play/stones/assistive/gossip)




# Current implementations



## IPFS


[IPFS](https://github.com/ipfs/specs "https://github.com/ipfs/specs") combines Kademlia + BitTorrent protocol + Git ideas to provide a high-throughput (_handling data at a high rate_), content-addressed block storage model (_stores user data and its information as separate objects in a store that can hold any kind of file system_), with content-addressed hyperlinks (_making the information retrievable over internet connections based on its content, not its location_).

-   Each file and all of the blocks in it have a unique fingerprint (hash).
-   Version history of each file is tracked by IPFS. 
-   Duplications are removed across the network. 
-   A node stores only content it cares about and some indexing information (for who is storing what).
-   Looking for files means looking for nodes storing the content behind a hash (a unique fingerprint).
-   Every file can be found by its human-readable name using the IPNS naming system.
    

#### Use cases


IPFS was designed as a distributed storage technology, useful for archivists, service providers and content creators. The number of users have increased a lot recently, as well as its intended scope and its [applications](https://github.com/ipfs/ipfs/labels/applications%20of%20ipfs "https://github.com/ipfs/ipfs/labels/applications%20of%20ipfs"). For example, IPFS makes it possible to run a decentralized application (ÐAPP) and store complex and unstructured content with a distributed file system based on IPFS that is not connected to the public IPFS network, so that its content isn’t spread on other IPFS hosts outside the network.

Currently, IPFS intends to replace the HTTP protocol. There are some serious hurdles to take. Internet standards get set by publishing RFCs, having multiple independent implementations, and going through a standards process by the IETF. The [IRTF Decentralized Internet Infrastructure (DIN)](https://trac.ietf.org/trac/dinrg/wiki "https://trac.ietf.org/trac/dinrg/wiki") has noticed it as distributed storage technology, and in general, the need for distributed technologies.

#### Merkle DAGs


A future-proof system that allows for multiple different fingerprinting mechanisms (summaries of the content that can be used to address content) to coexist. How? Large files are chunked, hashed, and organised into an Interplanetary Linked Data (IPLD) structure, a Merkle DAG object.

#### Multihash format

Raw content is run through a hash function, to produce a digest. This digest is said to be guaranteed to be cryptographically unique to the contents of the file, and that file only. The [multihash format](https://github.com/multiformats/go-multihash/tree/master/multihash "https://github.com/multiformats/go-multihash/tree/master/multihash") provides a [wrapper around the hash](https://github.com/multiformats/go-multihash/blob/master/multihash.go#L146 "https://github.com/multiformats/go-multihash/blob/master/multihash.go#L146"): The hash itself specifies which hash function is used, and the length of the resultant hash in the first two bytes (`fn code` and `length`) of the multihash. The rest of it is the `hash digest`.

[![Multihash format](https://niverel.tymyrddin.space/_media/en/research/dawnbreaker/ipfs/multihash.png)](https://niverel.tymyrddin.space/_detail/en/research/dawnbreaker/ipfs/multihash.png?id=en%3Aplay%3Astones%3Acurrent%3Aipfs "en:research:dawnbreaker:ipfs:multihash.png")

IPFS comes with a default hash algorithm but can be recompiled to use another hash function as default or to change the importer code to add a way to specify the multihash choice. When the hashing algorithm used is changed from SHA256 to BLAKE2b, the prefixes in the wrapper will differ.

#### Base58

Base58 is a group of binary-to-text encoding schemes used to represent large integers as alphanumeric text. It is designed for humans to easily enter the data, copying from some visual source, but also allows easy copy and paste because a double-click will usually select the whole string. It is similar to Base64 but has been modified to avoid both non-alphanumeric characters and letters which might look ambiguous when printed - similar-looking letters are omitted: `0` _(zero)_, `O` _(capital o)_, `I` _(capital i)_ and `l` _(lower case L)_, and non-alphanumeric characters `+` _(plus)_ and `/` _(slash)_ are dropped.

The actual order of letters in the alphabet depends on the application, which is the reason why the term Base58 alone is not enough to fully describe the format. Base58Check is a Base58 encoding format that unambiguously encodes the type of data in the first few characters and includes an error detection code in the last few characters.

For example, the base58 letters `Qm` correspond with hexadecimal `12` (SHA-256 algorithm) and hexadecimal `20` (length 32 bytes).

#### CID version

A CID is a [self-describing content-addressed identifier](https://github.com/ipld/cid/blob/master/original-rfc.md "https://github.com/ipld/cid/blob/master/original-rfc.md"). It uses cryptographic hashes to achieve content addressing. It uses several multiformats to achieve flexible self-description, namely multihash for hashes, multicodec for data content types, and multibase to encode the CID itself into strings.

Concretely, it's typed content address: a tuple of (content-type, content-address)

#### Notes on zero duplication

-   The resulting hash is not only a result of the chosen hash algorithm (`hash` option), but also affected by the choice of chunking algorithm (`chunker` option), DAG format (`trickle` option) and CID version (`cid-version` option), so it is possible to have completely different hashes even if the format is marked the same.
    
-   The same file can be duplicated across the network: Someone could add a file, remove it, upgrade the IPFS client (or change to using a different one), add it again, and get a completely different hash. This requires intent and is complicated and time-consuming, so the probability of the existence of multiple hashes for the same file is low.
    
-   Duplication of files leads to redundancy, something users might want and even may consider necessary.
    

“Zero duplication” refers to not having wasteful duplicates.

### IPFS Components


[![Incentergy: IPFS the next generation internet procotol technical overview](https://niverel.tymyrddin.space/_media/en/research/dawnbreaker/ipfs/ipfs-and-ipns-block-overview.png "Incentergy: IPFS the next generation internet procotol technical overview")](https://www.incentergy.de/blog/2018/02/03/ipfs-the-next-generation-internet-procotol-technical-overview/ "https://www.incentergy.de/blog/2018/02/03/ipfs-the-next-generation-internet-procotol-technical-overview/")

#### Identity

The IPFS identities technology manages node identity generation and verification. Nodes are identified using a NodeID, which is a cryptographic hash of a public key. Each node stores its public and private keys, with the private key being encrypted with a passphrase.

The `NodeId`, the cryptographic hash3 of a public-key, is created with [S/Kademlia’s static crypto puzzle](http://www.scs.stanford.edu/~dm/home/papers/kpos.pdf "http://www.scs.stanford.edu/%7Edm/home/papers/kpos.pdf") (pdf).

#### Network

The IPFS network manages connections to other peers using various underlying network protocols. The network is configurable. IPFS nodes communicate with hundreds of other nodes in the network (or potentially, the entire internet). Some stack features of the IPFS network include transport, reliability, connectivity, integrity, and authenticity systems.

-   IPFS can use any transport protocol and is best suited for WebRTC Data Channels (browser connectivity) or uTP.
    
-   If underlying networks do not provide reliability, IPFS can provide it using uTP or SCTP.
    
-   For connectivity it also uses ICE NAT traversal techniques - STUN is a standardized set of methods and a network protocol for NAT hole punching, designed for UDP but which was extended to TCP. TURN is a NAT traversal relay protocol. ICE is a protocol for using STUN and/or TURN to do NAT traversal while picking the best network route available. It fills in some of the missing pieces and deficiencies that were not mentioned in the STUN specification.
    
-   Integrity of messages can be checked using a hash checksum.
    
-   Authenticity of messages can be checked using HMAC with the sender’s public key.
    

#### Routing

The IPFS routing system maintains information to locate specific peers and objects. It responds to both local and remote queries. Routing defaults to a DHT (Dynamic Hash Table), but it’s swappable.

IPFS uses a DHT based on S/Kademlia and Coral (DSHT). Coral stores the addresses of peers who can provide the data blocks taking advantage of data locality. Coral can distribute only subsets of the values to the nearest nodes avoiding hot spots. Coral organises a hierarchy of separate DSHTs called clusters depending on region and size. This enables nodes to query peers in their region first, “finding nearby data without querying distant nodes” and greatly reducing the latency of lookups.

Coral and mainline DHT use DHTs as a place to store – not the value, but – pointers to peers who have the actual value. IPFS uses the DHT in the same way. When a node advertises a block available for download, IPFS stores a record in the DHT with its own `Peer.ID`. This is called “providing” and the node becomes a “provider”. Requesters who wish to retrieve the content, query the DHT (or DSHT) and need only to retrieve a subset of providers, not all of them. Providing is done once per block because blocks (even sub-blocks) are independently addressable by their hash.

#### Exchange

IPFS uses a unique block exchange protocol called BitSwap to govern efficient block distribution. BitSwap is modelled like a market, and users have some minor incentives for data replication. Trade strategies are swappable.

Unlike [BitTorrent](http://www.bittorrent.org/beps/bep_0003.html "http://www.bittorrent.org/beps/bep_0003.html"), BitSwap is not limited to the blocks in one torrent. The blocks can come from completely unrelated files in the filesystem. BitSwap incentivises nodes to seed/serve blocks even when they do not need anything in particular. To avoid leeches (freeloading nodes that never share), peers track their balance (in bytes verified) with other nodes, and peers send blocks to debtor peers according to a function that falls as debt increases.

-   If a node is storing a node that is the parent (root/ancestor) of other nodes, then it is much more likely to also be storing the children. So when a requester attempts to pull down a large DAG, it first queries the DHT for providers of the root. Once the requester finds some and connects directly to retrieve the blocks, BitSwap will optimistically send them the “wantlist”, which will usually obviate any more DHT queries for that DAG.   
-   BitSwap only knows about Routing. And it only uses the `Provide` and `FindProviders` calls.
    

#### Objects

Merkle DAGs of content-addressed immutable objects with links are used to represent arbitrary data structures, including file hierarchies and communication systems.

Merkle DAGs provide:

-   _Content addressing_: All content is uniquely identified by its multihash checksum.    
-   _Tamper resistance_: All content is verified with its checksum.    
-   _Deduplication_: All objects that hold the same content are equal, and only stored once.
    

#### Files

IPFS uses a versioned file system hierarchy inspired by [Git](https://git-scm.com/ "https://git-scm.com/"). On Github, the complete file history and changes over time can be viewed.

IPFS defines a set of objects for modelling a versioned filesystem on top of the Merkle DAG. This object model is similar to the Git model:

-   A _block_ is a variable-size block of data.    
-   A _list_ is an ordered collection of blocks or other lists.    
-   A _tree_ is a collection of blocks, lists, or other trees.    
-   A _commit_ is a snapshot in the version history of a tree.
    

#### Naming

IPFS uses a self-certifying mutable name system called IPNS. IPNS was inspired by [SFS](https://pdos.csail.mit.edu/papers/sfs:euresti-meng.pdf "https://pdos.csail.mit.edu/papers/sfs:euresti-meng.pdf") (pdf) and is compatible with given services like [DNS](https://tools.ietf.org/html/rfc1035 "https://tools.ietf.org/html/rfc1035").

`NodeId` is obtained by `hash(node.PubKey)`. Then IPNS assigns every user a mutable namespace at: `/ipns/<NodeId>`. A user can publish an Object to this `/ipns/<NodeId>` path signed by his/her private key. When other users retrieve the object, they can check the signature matches the public key and `NodeId`. This verifies the authenticity of the _Object_ published by the user, achieving mutable state retrieval.

`<NodeId>` is a hash, it is not human friendly to pronounce and recall. That is where DNS TXT IPNS Records come in. If `/ipns/<domain>` is a valid domain name, IPFS looks up key `ipns` in its DNS TXT records. The `ipns` behaves as a symlink.

-   IPNS uses the `Put` and `GetValue` calls.
    

### Resources


[IPFS - Content Addressed, Versioned, P2P File System (DRAFT 3)](https://ipfs.io/ipfs/QmR7GSQM93Cx5eAg6a6yRzNde1FQv7uL6X1o4k7zrJa3LX/ipfs.draft3.pdf "https://ipfs.io/ipfs/QmR7GSQM93Cx5eAg6a6yRzNde1FQv7uL6X1o4k7zrJa3LX/ipfs.draft3.pdf"), Juan Benet

[Source](https://niverel.tymyrddin.space/en/play/stones/current/ipfs)


## BTFS


BTFS is a fork of [IPFS](https://niverel.tymyrddin.space/en/play/stones/current/ipfs "en:play:stones:current:ipfs").

## FileCoin


[IPFS](https://niverel.tymyrddin.space/en/play/stones/current/ipfs "en:play:stones:current:ipfs") with FileCoin offers storage on a global network of local providers who have the freedom to set prices based on supply and demand. It implements a generalised version of the BitTorrent exchange protocol and uses Proof-of-Storage (instead of a Proof-of-Work consensus algorithm like Bitcoin). Anyone can join the network, offer unused hard drive space, and get rewarded in FileCoin tokens for data storage and retrieval services. Filecoin is traded on several exchanges and supported by multiple cryptocurrency wallets, allowing the exchange FileCoin for other currencies like Euros, US Dollars, BTC and ETH.

### Resources


-   [Filecoin \[IOU\] (FIL)](https://www.coingecko.com/en/coins/filecoin "https://www.coingecko.com/en/coins/filecoin"), CoinGecko    
-   [Filecoin](https://filecoin.io/ "https://filecoin.io/")    
-   [Filecoin: A Decentralized Storage Network](https://filecoin.io/filecoin.pdf "https://filecoin.io/filecoin.pdf"), Protocol Labs, 2017

[Source](https://niverel.tymyrddin.space/en/play/stones/current/filecoin)


## The usual peer crawling


### Requirements


-   Full distribution of every task (no centralized coordination at all): Peers can perform their jobs independently and communicate required data. This will prevent link congestion to a central server. Functionally identically programmed nodes, distinguished by a unique identifier only.    
-   No single point of failure and gracefully dealing with permanent and transient failures (identify crawl traps and be tolerant to external failures)    
    -   Nodes do not crash on the failure of a single peer. ⇒ dynamic reallocation of addresses can be done across other peers.        
    -   Data sent to a node while it is down still needs to propagate properly throughout the system.        
    -   Data is still retrievable even with a node failing.        
    -   When a node comes back online after having failed, data that needs to be stored on that node is propagated to it and can be retrieved from it. 
    -   If a node failure is permanent, the node needs to be recoverable using data stored by other nodes.        
-   Scalability - Not relying on location implies latency can become an issue ⇒ minimising communication.    
-   Platform independence of nodes.    
-   Locally computable content address assignment based on consistent hashing - IPNS is consistent (check for the presence of DNSLink first for it is faster)    
-   Portability: The nodes can be configured to run on any kind of dweb network by just replacing it.    
-   Performance: Nodes run on spare CPU processing power and are not to put too large a load on a client machine.    
-   Freshness: After an object is initially acquired (processed), it may have to be periodically recrawled and checked for updates. In the simplest case, this could be done by starting another broad breadth-first crawl, or by requesting all items in the collection of a node again. Techniques for optimizing the “freshness” of such collections is usually based on observations about an item's update history (incremental crawling). Sadly, we cannot use this “as such” as an updated object in IPFS has a new hash. A variety of other heuristics can be used to recrawl as “more important” marked items. Good enough recrawling strategies are essential for maintaining an up-to-date search index with limited crawling bandwidth.    
-   Note that security and privacy have no requirements yet as this distributed crawler is a “hello world” function. The component with the distribution function can be further expanded once the distribution of the index is designed.
    

### Components


The below compartmentalised architecture helps to make the system portable. To make it pluggable we will need a high degree of functional independence between the components of the system so that we can run the crawler on top of different underlying dweb networks and different distribution functions (with minimal regression factor) to optimise the distribution function for the underlying dweb network architecture. Combining the three components can be done with a configuration file, created during installation, with the definition of communication interfaces of the dweb network and other required parameters. That means we would only need platform-dependent installers.

The three proposed crawler components:

-   An **_Overlay Network Layer_** responsible for formation and maintenance of a distributed search engine network, and communication between peers. These can be unstructured or structured networks. If a network is not scalable, a supernode architecture can be used to improve performance, hence a client must have support for flat as well as supernode architecture.    
-   A **_Peer and Content Distribution Function_** determines which clients to connect with. Each client has a copy of this function. Hash list & content range associated with a client can change due to joining/leaving of nodes in the network. The function will distribute hashes to crawl as well as content among peers and makes use of the underlying dweb network to provide load balancing and scalability and takes proximity of nodes into account. Initially, we will use a static distribution function, hash list and content range assignment functions can be hash functions.    
-   The** _Crawler_** downloads and extracts dweb objects and sends and receives data from other peers using the underlying network.
    

### Crawler


**Datastores:**
    
-   The **_Neighbourhood_** data store contains the identifiers of agents in a node's neighbourhood.        
-   The _**Seen Hash**_ data store contains a list of hashes that are already processed.
-   The **_Seen Content_** data store contains attributes of objects that are already processed.
         

Whenever a node receives a hash for processing from another node, the **_Preprocessor_** adds the hash to one of the **_Crawl Job Queues_**.
    
-   A classifier determines to which crawl job queue a request is added. Hashes from the same node are added to the same crawl job queue.
-   Rate of input is most likely much faster than the rate of output. We can implement an overflow controller that drops hashes after the crawl job queue overflows.
-   Rate throttling to prevent excessive requests to the same node.
-   Checks whether a hash has already been processed by accessing the **_Seen Hash_** data structure. If a hash is not already processed then the hash is added to the **_Seen Hash_** data structure.
-   We start with just one queue and can then experiment with adding more.
-   Each crawl job queue has one **_Content Fetcher_** thread associated with it that streams data into the extractor.    
    -   Uses the local IPFS gateway to fetch a (named) IPFS resource.      
    -   Checks file permissions.        
    -   If we keep the connection open to a node to minimise connection establishment overhead we will need exception handling.        
-   The **_Extractor_** extracts linked hashes/links from the object and passes them on to the hash validator. It can perhaps be multithreaded and process different pages simultaneously.    
-   The **_Hash Validator_** checks whether a hash is the responsibility of the node. If so, it sends the hash to the preprocessor. If not, it is sent out on the network.    
-   The **_Content Range Validator_** checks whether content lies in the range of the node.    
-   The **_Content Processor_** checks object for duplication from _Seen Content_ data. If the object is not already processed, then it is added to the _Seen Content_ data store.
    

### Resources


[UbiCrawler: A Scalable Fully Distributed Web Crawler](http://vigna.di.unimi.it/ftp/papers/UbiCrawler.pdf "http://vigna.di.unimi.it/ftp/papers/UbiCrawler.pdf"), Paolo Boldi, Bruno Codenotti, Massimo Santini, Sebastiano Vigna    

[Apoidea: A Decentralized Peer-to-Peer Architecture for Crawling the World Wide Web](https://www.cc.gatech.edu/~lingliu/papers/2003/apoidea-sigir03.pdf "https://www.cc.gatech.edu/~lingliu/papers/2003/apoidea-sigir03.pdf"), Aameek Singh, Mudhakar Srivatsa, LingLiu, and Todd Miller, 2003    

[Efficient Crawling Through URL Ordering](http://ilpubs.stanford.edu:8090/347/1/1998-51.pdf "http://ilpubs.stanford.edu:8090/347/1/1998-51.pdf"), Junghoo Cho, Hector Garcia-Molina, Lawrence Page, 1998

[Source](https://niverel.tymyrddin.space/en/play/stones/current/peer-crawling)


# Internet-facing demo


Mapping the theoretical limitations of decentralised/distributed search and how they might be ‘overcome’ with ‘good enough’ practical implementations/‘approximations’.


## Good enough indexing


If index construction needs to be done for different platforms with different hardware specs (full distribution of search engine), while still being scalable, using [single pass in-memory sorting (SPIMI)](https://niverel.tymyrddin.space/en/play/algos/index#single-pass-in-memory-sorting-spimi "en:play:algos:index") looks like a 'good enough' candidate. SPIMI uses terms instead of term-id's, writes each block's dictionary to disk and then starts a new dictionary for the next block. Otherwise, we can use [blocked sort based indexing (BSBI)](https://niverel.tymyrddin.space/en/play/stones/demo/indexing#blocked-sort-based-indexing-bsbi "en:play:stones:demo:indexing ↵"). Blocked sort-based indexing has awesome scaling properties, but needs a data structure for mapping terms to term-id. For very large collections, it may not fit into memory.
    
We can use known [distributed indexing algorithms](https://niverel.tymyrddin.space/en/play/algos/dindex "en:play:algos:dindex"). Problem is that MapReduce tasks must be written as acyclic dataflow programs, in essence a stateless mapper followed by a stateless reducer, that is executed by a batch job scheduler. This paradigm makes repeated querying of datasets difficult and imposes limitations on machine learning applications (like we may wish to use for PageRank and TrustRank), where iterative algorithms that revisit a single working set multiple times are the norm. The MapReduce model seems coupled to the Hadoop infrastructure. The [Bulk Synchronous Parallel (BSP)](http://www.bsp-worldwide.org/ "http://www.bsp-worldwide.org/") computing model for parallel programming may be a viable alternative. [Apache Hama](https://hama.apache.org/ "https://hama.apache.org/") implements BSP.
    
Pages and objects come in over time and need to be inserted, other pages and objects have been deleted and this means that all of the indexes have to be modified too: For documents, this means postings updates for terms already in the dictionary, new terms need to be added to the dictionary, and all associated indexes such as N-gram indexes will have to be updated for each added or deleted document. This can be made easy by making an auxiliary index and [logarithmic merge](https://niverel.tymyrddin.space/en/play/algos/dynindex#logarithmic-merge "en:play:algos:dynindex") of the main and auxiliary index.
    

[The Graph](https://medium.com/graphprotocol "https://medium.com/graphprotocol") is a protocol for building decentralized applications using GraphQL. In essence, [the graph](https://thegraph.com/ "https://thegraph.com/") is a decentralized index that works across blockchains (it can index data in multiple blockchains like Ethereum and BTC, but also on IPFS and Filecoin). It monitors the blockchains for new data and updates the index every time this happens. Once the index is updated, it tries to reach a consensus among the nodes that maintain it. Once consensus is reached, it ensures that the users of the index will have the latest data available. Not all challenges are solved yet. Who decides what when? Changes and how to adopt these changes? Topologies? Scalability? Update performance? How to ensure everyone has or has access to the same new “fresh” index? Time to contact that group …


# Distributed search

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

[Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/start)


## Overlay networks


An IPFS node can be fingerprinted through the content it stores. An overlay network needs to offer an “anonymous” mode that only enables features known to not leak information.

-   No local discovery.    
-   No transports other than, for example, via Tor (an overlay network consisting of more than seven thousand relays to conceal a user's location and usage from anyone conducting network surveillance or traffic analysis). 
-   Private routing to make the network non-enumerable.
    

And getting any of this wrong could put _some_ people in danger.

[Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/overlay)


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

[Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/extraction)

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


[Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/ipfs-cluster)

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



[Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/querying)

## Yggdrasil


Yggdrasil is an early-stage implementation of a fully end-to-end encrypted IPv6 network. It is lightweight, self-arranging, supported on multiple platforms and allows pretty much any IPv6-capable application to communicate securely with other Yggdrasil nodes. Yggdrasil does not require IPv6 Internet connectivity - it also works over IPv4.

Looking at it for its clustering and bootstrapping implementation.

### Resources


-   [Yggdrasil Version 0.3.6](https://yggdrasil-network.github.io/2019/08/03/release-v0-3-6.html "https://yggdrasil-network.github.io/2019/08/03/release-v0-3-6.html"), august 2019, first version with API    
-   [Yggdrasil](https://github.com/yggdrasil-network/yggdrasil-go "https://github.com/yggdrasil-network/yggdrasil-go"), Github

[Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/yggdrasil)

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


[Source](https://niverel.tymyrddin.space/en/play/stones/upsidedown/testing)




















# 