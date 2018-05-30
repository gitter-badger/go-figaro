# go-figaro

Official Figaro Blockchain - Go Version

## Performance Considerations

- Uses Blake2b hash functions, which perform well on multi-core 64-bit architectures
- Uses BadgerDB, which is optimized for fast SSD storage

## Data Storage

Even with space-optimized storage schemes, targeting 56000 tx/sec means we need to make liberal use of pruning. Since we must prune regardless, Figaro balances performance, storage, and master-node/light-client dichotomies with the following storage scheme:

- Canonical chain head saved in store under prefix+"head"
- Canonical chain blocks IDs saved in store under prefix+blockNum
- Block headers saved in store under block ID
- Block commits saved in archive trie w root in block header
- Block transactions saved in archive trie w root in block header
- Receipts saved in store under prefix+txid + root in block header
- Accounts are saved in state trie under world state root + root in block header
- Contract storage saved in state trie under account state root

### Block Variants

- Header: contains header information only
- CompBlock: contains Header + commits bloom + transactions bloom (variable size, fixed 3% false positive rate)
- RefBlock: contains Header + commits + transaction hashes (ids)
- Block: contains Header + commits + transactions

Each block variant is reproducible from headers alone by any full node variant.

## WIP Data/Validation Protocols

A light-client can request account data or contract storage data and a master-node can look up the address value under any valid state trie merkle root and provide an efficient proof if requested.

For transaction data, a light-client can request by transaction ID, and a master-node can lookup the receipt directly, optionally returning a proof via receipt=>block=>transaction root=>archive trie proof. This means there is one canonical receipt for a transaction, and chain reorgs will overwrite the receipt.

Light clients must keep local copies of block headers and canonical chain data, and reorg whenever a chain conflict is encountered. They may request block header (by id or by canonical number), transaction receipt/proofs, and state data/proofs from master-nodes at any time.

Master nodes, in addition to the requirements of a light client, must also validate incoming blocks by validating the block signature, provider, and header data; validating a heuristic subset of transactions; executing each transaction, updating state, and generating, validating, and saving each receipt; archiving the transactions and commits and validating the associated block roots. When providing a block, the same process must be followed, except that all transaction signatures must be verified. Master nodes may request missing block and/or transaction data from other master nodes.

PendingTxList is a list of received transactions awaiting validation, sorted by received timestamp. `QuickCheck` and `VerifySignature` should be called on each (using parallel threads if possible) and the transaction moved to the TxList.

TxList is a map of nonce-ordered future transaction lists by account and a received timestamp ordered list of current transactions to be executed. When a transaction is added to the TxList, it is sorted into the current or future queues based on account nonce. When a transaction is pulled from the current queue, it is added into a block, which will sequentialy validate and execute the transactions in block order.

## Demo

To enable graph visualizations, install Graphviz:

    brew install graphviz
