# Distributed and assistive technologies

## Distributed technologies

Technologies that can be considered 'distributed'.

### Blockchain 

Stuart Haber and W. Scott Stornetta already envisioned a cryptographically secured chain of blocks whereby no one could tamper with timestamps of documents in 1991. In 1992, they upgraded their system to use Merkle trees, increasing efficiency and enabling the collection of more documents on a single block. Satoshi Nakamoto, a person or a group of people, developed the first application of the digital ledger technology in 2008, BitCoin.

* Data is structured in blocks in order of transactions that are validated by miners.
* Each block produces a unique hash that identifies the transaction. If one attempts to alter the details of the transaction, a different hash will be generated. This can be evidence of a corrupted and invalid transaction.
* Transactions are published on a public ledger to which every node has access (transparency). The distributed nature of the public ledger makes it even more difficult for parties to tamper with information.
* Miners can postpone or even cancel a transaction.
* Traditional Blockchains rely on Proof of Work. These need many computations and as a result, the number of transactions per second is relatively low.
    * A transaction has to validate numerous transactions before being valid.
    * As blocks in blockchain multiply, it becomes increasingly difficult in terms of computations to achieve new blocks and mining becomes more power-intensive (expensive).

##### Use cases
Cryptocurrencies

#### Resources 
[Bitcoin: A Peer-to-Peer Electronic Cash System](https://bitcoin.org/bitcoin.pdf)

#### [Source](https://niverel.tymyrddin.space/en/play/stones/dweb/blockchain)

### Hashgraph

The hashgraph algorithm was invented by Leemon Baird for achieving consensus quickly, fairly, efficiently, and securely.

-   Hashgraph achieves transaction success solely via consensus timestamping to make sure that transactions on the network agree with each node on the platform.    
-   On a Hashgraph network nodes do not have to validate transactions by _Proof of Work_ or _Proof of Stake_. Consensus is built with the _Gossip about Gossip_ and _Virtual Voting_ techniques instead, increasing the number of transactions per second.    
-   And consensus timestamping avoids the [Blockchain](https://niverel.tymyrddin.space/en/play/stones/dweb/blockchain "en:play:stones:dweb:blockchain") issues of cancelling transactions or by putting them on future blocks.    
-   These consensus techniques also facilitate fairness.    
-   Developers do not need a license but need the platform coin instead. API calls cost a micro-payment to the company.
    

##### Use cases


All use cases where trust is immutable and incorruptible, for example:

-   Cryptocurrency as a service for support for native micropayments    
-   Micro-storage in the form of a distributed file service that apps can use
-   Contracts    
-   Bank transfers    
-   Credential verification
    

#### Resources
-   [Swirlds](https://www.swirlds.com/ "https://www.swirlds.com/")    
-   [Hedera Hashgraph](https://www.hedera.com/ "https://www.hedera.com/")

#### [Source](https://niverel.tymyrddin.space/en/play/stones/dweb/hashgraph)


### DAG

A DAG is a type of distributed ledger technology that relies on consensus algorithms. To prevail, transactions require majority support within the network. As a result, there is more cooperation and teamwork and nodes have equal rights. Such networks stick to the original goal of Distributed Ledger Technology, to democratise the internet economy.

-   No blocks. No chain. DAG is a structure that is connected like a mesh.
-   It connects current data transactions with previous ones.    
-   With nodes having equal rights, nodes do not have to refer to another node.    
-   A consensus-based system where nodes decide what happens to give a semblance of democracy as compared to platforms that go through a central command.    
-   For a transaction to succeed, it has to validate only two of the previous transactions.    
-   Transactions in DAGs adds throughput as many more validations happen.
    

##### Use cases
-   Cryptocurrencies    
-   Economic infrastructure for data sharing on the Internet of Things    
-   Remote Patient Monitoring    
-   Decentralised Peer-to-Peer energy trading
    

#### Resources
-   [OByte](https://obyte.org/ "https://obyte.org/")    
-   [IoTA Use Cases](https://files.iota.org/comms/IOTA_Use_Cases.pdf "https://files.iota.org/comms/IOTA_Use_Cases.pdf")

#### [Source](https://niverel.tymyrddin.space/en/play/stones/dweb/dag)




### Holochain


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
        

##### Use cases


Systems where not all parties need to participate:

-   Social networks    
-   Chat programs    
-   p2p platforms    
-   Shared document updates
    

#### Resources


-   [r/holochain: Distributed Computing and Applications](https://www.reddit.com/r/holochain/ "https://www.reddit.com/r/holochain/")    
-   [Holochain projects](http://holochainprojects.com/ "http://holochainprojects.com/")    
-   [Decentralising the web: The key takeaways](https://www.computing.co.uk/ctg/news/3036546/decentralising-the-web-the-key-takeaways "https://www.computing.co.uk/ctg/news/3036546/decentralising-the-web-the-key-takeaways"), 2018    
-   [Holochains for Distributed Data Integrity](http://ceptr.org/projects/holochain "http://ceptr.org/projects/holochain")   
-   [Holochain](https://holochain.org/ "https://holochain.org/")

#### [Source](https://niverel.tymyrddin.space/en/play/stones/dweb/holochain)

## Assistive technologies


A truly decentralised web will require the network to provide privacy and trust **by design**. This requires algorithms that allow for trustless management. Zero knowledge proofs and/or cross-validation enables nodes to verify the existence and validity of exchanges. The challenge is to maintain a distributed consensus, without actually being able to see or make public any of the transaction details, guaranteeing privacy.

### Zero Knowledge proofs


A [blockchain](https://niverel.tymyrddin.space/en/play/stones/dweb/blockchain "en:play:stones:dweb: blockchain") is a data structure, a linear transaction log, replicated by devices whose users are rewarded for logging new transactions.

-   A change in any block invalidates every block after it, which means that an adversary can not tamper with historical transactions.    
-   A user only gets rewarded if they are working on the same chain as everyone else, so each participant has an incentive to go with the consensus. The result is a shared definitive historical record.
    

Devil's advocate:

-   It is not truly “trustless”, because most of its users are trusting the software, instead of trusting other people.    
-   Blockchain systems do not magically make the data in them accurate or the people entering the data trustworthy, they merely enable any user to audit whether it has been tampered with.    
-   Making it possible to download the blockchain from a broadcast node and decrypt the Merkle root from the Linux command line to independently verify transactions will only be used by the tech-savvy.
    

#### Millionaire’s Problem


In Yao’s Millionaire’s Problem, two millionaires want to find out if they have the same amount of money without disclosing the exact amount. This problem is analogous to a more general problem where there are two numbers `a` and `b` and the goal is to determine whether the inequality `a ≥ b` is true or false without revealing the actual values of `a` and `b`. Problems like these can be used in Zero Knowledge Protocols or Zero Knowledge Password Proofs (ZKPs). The latter, Zero Knowledge Password Proof, is a way of doing authentication where no passwords are exchanged, which means they cannot be stolen. The first term, Zero Knowledge Protocol, has appeared more frequently within blockchain circles.

### Zero Knowledge Protocols


A protocol has to have the following properties to be considered Zero Knowledge Proof:

-   Completeness — If the statement is true, an honest verifier (that is, one following the protocol properly) will be convinced of this fact by an honest prover.    
-   Soundness — If the statement is false, no cheating prover can convince an honest verifier that it is true, except with some small probability.    
-   Zero-Knowledge — If the statement is true, no verifier learns anything other than the fact that the statement is true. In other words, just knowing the statement (not the secret) is sufficient to imagine a scenario showing that the prover knows the secret. This is formalized by showing that every verifier has some simulator that, given only the statement to be proved (and no access to the prover), can produce a transcript that “looks like” an interaction between the honest prover and the verifier in question.
    

Zero-knowledge proofs are probabilistic “proofs” (useful for “good enough” algorithms) rather than deterministic proofs, and techniques exist to decrease the soundness error to negligibly small values.

#### Resources


-   September 2017, [the first ZKP was conducted on the Byzantium fork of Ethereum](https://cointelegraph.com/news/ethereum-upgrade-byzantium-is-live-verifies-first-zk-snark-proof "https://cointelegraph.com/news/ethereum-upgrade-byzantium-is-live-verifies-first-zk-snark-proof"), and perhaps we can adapt it for our purpose.    
-   [The knowledge complexity of interactive proof systems](https://people.csail.mit.edu/silvio/Selected%20Scientific%20Papers/Proof%20Systems/The_Knowledge_Complexity_Of_Interactive_Proof_Systems.pdf "https://people.csail.mit.edu/silvio/Selected%20Scientific%20Papers/Proof%20Systems/The_Knowledge_Complexity_Of_Interactive_Proof_Systems.pdf"), Shafi Goldwasser, Silvio Micali, Charles Rackoff, February 1989    
-   [How to Explain Zero-Knowledge Protocols to Your Children](http://pages.cs.wisc.edu/~mkowalcz/628.pdf "http://pages.cs.wisc.edu/~mkowalcz/628.pdf"), Jean-Jacques Quisquater, Louis Guillou, Thomas Berson, 1990    
-   [A Survey of Noninteractive Zero-KnowledgeProof System and Its Applications](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC4032740/pdf/TSWJ2014-560484.pdf "https://www.ncbi.nlm.nih.gov/pmc/articles/PMC4032740/pdf/TSWJ2014-560484.pdf"), Huixin Wu and Feng Wang, May 2014    
-   [A Mind-Bending Cryptographic Trick Promises to Take Blockchains Mainstream](https://www.technologyreview.com/s/609448/a-mind-bending-cryptographic-trick-promises-to-take-blockchains-mainstream/ "https://www.technologyreview.com/s/609448/a-mind-bending-cryptographic-trick-promises-to-take-blockchains-mainstream/"), Mike Orcutt, November 2017

#### [Source](https://niverel.tymyrddin.space/en/play/stones/assistive/zkp)



### Smart contracts


A smart contract is a computer protocol intended to digitally facilitate, verify, or enforce the negotiation or performance of a contract. Smart contracts allow the performance of credible transactions between disparate, anonymous parties without third parties (central authority, legal system, or external enforcement). These transactions are trackable and irreversible and serve to increase the integrity of the ledger.

#### Resources


-   [Research Directions in Blockchain Data Management and Analytics](https://openproceedings.org/2018/conf/edbt/paper-227.pdf "https://openproceedings.org/2018/conf/edbt/paper-227.pdf"), Hoang Tam Vo, Ashish Kundu, Mukesh Mohania, 2018

#### [Source](https://niverel.tymyrddin.space/en/play/stones/assistive/contracts)



### Voting


Voting models for consensus have certain security properties (they can be asynchronous Byzantine Fault Tolerant and achieve consensus even when some nodes are malicious and some messages are significantly delayed), but generally require that a large number of messages be sent amongst nodes to get to consensus. Hashgraph has come up with a voting mechanism where this would not be necessary …

#### Resources
[I Want Your Vote! (Oh Wait I Already Know It)](https://medium.com/hashgraph/i-want-your-vote-oh-wait-i-already-know-it-e1faa50b31ad "https://medium.com/hashgraph/i-want-your-vote-oh-wait-i-already-know-it-e1faa50b31ad"), Paul Madsen, 2017

#### [Source](https://niverel.tymyrddin.space/en/play/stones/assistive/voting)



### Gossip about gossip


Gossip about gossip, based on the gossip communication protocol, enables nodes to efficiently and rapidly exchange data with other nodes at internet scale. There are currently three known libraries that implement a gossip algorithm to discover nodes in a peer-to-peer network: Apache Gossip, gossip-python and Smudge.

#### Resources


[Apache Gossip](https://github.com/apache/incubator-gossip "https://github.com/apache/incubator-gossip") (UDP, Java)    
[gossip-python](https://github.com/thomai/gossip-python "https://github.com/thomai/gossip-python") (TCP)    
[Smudge](https://github.com/clockworksoul/smudge "https://github.com/clockworksoul/smudge") (UDP, Go)

#### [Source](https://niverel.tymyrddin.space/en/play/stones/assistive/gossip)