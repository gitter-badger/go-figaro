// Package figaro is the main package for go-figaro
package figaro

// ConsensusEngine implements the rules by which things are validated. It handles
// enforcement of configurable rules, while types handle enforcement of static rules.
type ConsensusEngine interface {
	// NextBlockProducer must deterministically decide on the next block producer
	// masternode address based on the previous block.
	NextBlockProducer(db FullDataService, prevblock BlockHash) (Address, error)

	// HandleFraud is responsible for enforcing conensus rules when a block is found to
	// contain fraudulent transactions or headers.
	HandleFraud(db FullDataService, fraudblock *BlockHeader) error

	// ChainReorg is responsibile for determining a canonical chain in the event of divergent,
	// but otherwise valid, chains/blocks. It receives the current chain, the block that would conflict
	// with the current chain head, and the list of pending fugure blocks. It should return the new chain,
	// along with the next and future blocks in the canonical chain.
	ChainReorg(db FullDataService, chain *Chain, forkblock *BlockHeader, futureblocks *BlockHeap) (*Chain, *BlockHeader, *BlockHeap, error)
}

// FullDataService provides full data for all chain types.
type FullDataService interface {
	AccountDataService
	CommitDataService
	TransactionDataService
	ReceiptDataService
	BlockDataService
	ChainDataService
}
