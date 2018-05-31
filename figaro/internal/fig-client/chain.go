package internal

import (
	"bytes"
	"container/heap"
	"reflect"

	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figaro/internal/fig-node/figdb"
)

// HandleReceiveBlock handles validating and syncing a new block created from the network
func HandleReceiveBlock(db *figdb.DB, chain *figaro.Chain, block *figaro.BlockHeader, futureblocks *figaro.BlockHeap, engine figaro.ConsensusEngine) error {
	if !block.VerifySignature() {
		return figaro.ErrInvalidBlock
	}
	// If the block is the future, we'll come back to it. If it's in the
	// past we'll ignore it (we only reorg if we encounter a longer chain).
	if block.Number > chain.Depth+1 {
		heap.Push(futureblocks, block)
		return figaro.ErrInvalidBlock
	}
	if block.Number != chain.Depth+1 {
		return figaro.ErrInvalidBlock
	}
	if reflect.DeepEqual(block.ChainConfig, chain.ChainConfig) {
		return figaro.ErrInvalidBlock
	}
	next, err := engine.NextBlockProducer(db, chain.Head)
	if err != nil {
		return err
	}
	if !bytes.Equal(block.Producer, next) {
		return figaro.ErrInvalidBlock
	}
	if !bytes.Equal(block.ParentBlock, chain.Head) {
		chain, block, futureblocks, err = engine.ChainReorg(db, chain, block, futureblocks)
		if err != nil {
			return err
		}
	}
	err = chain.AppendBlock(db, block)
	if err != nil {
		return err
	}
	if futureblocks.PeekNextNumber() == chain.Depth+1 {
		block = heap.Pop(futureblocks).(*figaro.BlockHeader)
		return HandleReceiveBlock(db, chain, block, futureblocks, engine)
	}
	return nil
}
