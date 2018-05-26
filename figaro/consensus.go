// Package figaro is the main package for go-figaro
package figaro

// ConsensusEngine implements the rules by which things are validated. It handles
// enforcement of configurable rules, while types handle enforcement of static rules.
type ConsensusEngine interface {
	// BlockValidate should validate whether the block is valid according to consensus rules.
	ValidateBlock(db FullChainDataService, block *BlockHeader) bool
	// BlockValidate should validate whether the block is valid according to consensus rules. It
	// is responsible for cleanup and returning a valid chain at the point of the fork, the next block,
	// in the new canonical branch, and all pending future blocks in the new canonical branch.
	ChainReorg(db FullChainDataService, chain Chain, block *Block, futureblocks *BlockHeap) (Chain, *Block, *BlockHeap)
}
