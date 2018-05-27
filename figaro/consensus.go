// Package figaro is the main package for go-figaro
package figaro

// ConsensusEngine implements the rules by which things are validated. It handles
// enforcement of configurable rules, while types handle enforcement of static rules.
type ConsensusEngine interface {
	// ValidateBlockProducer should validate whether the block is valid according to consensus rules. It
	// has sole responsibility for determining whether the producer is valid, but can assume
	// the signature is already validated.
	ValidateBlockProducer(db FullChainDataService, block *Block) bool

	// ValidateBlockTxs should validate whether the block is valid according to consensus rules. It
	// has sole responsibility for determining whether transactions signatures are valid or fraudulent,
	// and handle the case of fraudulent signatures in the block.
	ValidateBlockTxs(db FullChainDataService, block *Block) bool

	// HandleFraudHeaders has sole responsiblity for handling the case where a block is found to contain
	// invalid headers (roots don't validate), including determining and penalizing fraud.
	HandleInvalidHeaders(db FullChainDataService, block *Block) error

	// HandleDivergence is responsibile for determining a canonical chain in the event of divergent,
	// but otherwise valid, chains/blocks. It should return the new chain whose head is at the fork,
	// along with the next and future blocks in the canonical chain.
	HandleDivergence(db FullChainDataService, chain Chain, block *Block, futureblocks *BlockHeap) (Chain, Block, BlockHeap, error)
}
