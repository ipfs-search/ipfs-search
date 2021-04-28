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

### [Source](https://niverel.tymyrddin.space/en/play/stones/assistive/zkp
)



## Smart contracts


A smart contract is a computer protocol intended to digitally facilitate, verify, or enforce the negotiation or performance of a contract. Smart contracts allow the performance of credible transactions between disparate, anonymous parties without third parties (central authority, legal system, or external enforcement). These transactions are trackable and irreversible and serve to increase the integrity of the ledger.

### Resources


-   [Research Directions in Blockchain Data Management and Analytics](https://openproceedings.org/2018/conf/edbt/paper-227.pdf "https://openproceedings.org/2018/conf/edbt/paper-227.pdf"), Hoang Tam Vo, Ashish Kundu, Mukesh Mohania, 2018

### [Source](https://niverel.tymyrddin.space/en/play/stones/assistive/contracts)



## Voting


Voting models for consensus have certain security properties (they can be asynchronous Byzantine Fault Tolerant and achieve consensus even when some nodes are malicious and some messages are significantly delayed), but generally require that a large number of messages be sent amongst nodes to get to consensus. Hashgraph has come up with a voting mechanism where this would not be necessary …

### Resources
[I Want Your Vote! (Oh Wait I Already Know It)](https://medium.com/hashgraph/i-want-your-vote-oh-wait-i-already-know-it-e1faa50b31ad "https://medium.com/hashgraph/i-want-your-vote-oh-wait-i-already-know-it-e1faa50b31ad"), Paul Madsen, 2017

### [Source](https://niverel.tymyrddin.space/en/play/stones/assistive/voting)



## Gossip about gossip


Gossip about gossip, based on the gossip communication protocol, enables nodes to efficiently and rapidly exchange data with other nodes at internet scale. There are currently three known libraries that implement a gossip algorithm to discover nodes in a peer-to-peer network: Apache Gossip, gossip-python and Smudge.

### Resources


[Apache Gossip](https://github.com/apache/incubator-gossip "https://github.com/apache/incubator-gossip") (UDP, Java)    
[gossip-python](https://github.com/thomai/gossip-python "https://github.com/thomai/gossip-python") (TCP)    
[Smudge](https://github.com/clockworksoul/smudge "https://github.com/clockworksoul/smudge") (UDP, Go)

### [Source](https://niverel.tymyrddin.space/en/play/stones/assistive/gossip)