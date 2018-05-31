package internal

import (
	"bytes"
	"container/heap"
	"reflect"

	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figaro/internal/fig-node/figdb"
)

// HandleReceiveBlock handles validating and syncing a new block received from the network
func HandleReceiveBlock(db *figdb.DB, chain *figaro.Chain, block *figaro.Block, futureblocks *figaro.BlockHeap, engine figaro.ConsensusEngine) error {
	// If the block is the future, we'll come back to it.
	if block.Number > chain.Depth+1 {
		err := db.SaveBlock(block)
		if err != nil {
			return err
		}
		// TODO: wrap all these heap methods for type safety in `figaro`
		heap.Push(futureblocks, block.BlockHeader)
		return nil
	}
	// If the block is in the past, skip it, as we've already got a longer chain.
	// NOTE: this skipped block could be canonical, but we'll wait until we encounter
	// a longer chain and then request missing blocks from the network.
	if block.Number < chain.Depth+1 {
		// TODO: maybe should probably delete the block, since we don't need it anymore
		header := heap.Pop(futureblocks).(*figaro.BlockHeader)
		block, err := db.HydrateBlock(header)
		if err != nil {
			return err
		}
		return HandleReceiveBlock(db, chain, block, futureblocks, engine)
	}
	err := HandleNextBlock(db, chain, block, engine)
	if err != nil && err != figaro.ErrReorgRequired {
		return err
	}
	if err == figaro.ErrReorgRequired {
		var header *figaro.BlockHeader
		// This will also handle syncing the database after the reorg, so we'll have the block
		// data available to us by the time this returns
		chain, header, futureblocks, err = engine.ChainReorg(db, chain, block.BlockHeader, futureblocks)
		if err != nil {
			return err
		}
		block, err := db.HydrateBlock(header)
		if err != nil {
			return err
		}
		return HandleReceiveBlock(db, chain, block, futureblocks, engine)
	}
	// TODO: try to process future blocks recursively
	if futureblocks.PeekNextNumber() == chain.Depth+1 {
		header := heap.Pop(futureblocks).(*figaro.BlockHeader)
		block, err := db.HydrateBlock(header)
		if err != nil {
			return err
		}
		return HandleReceiveBlock(db, chain, block, futureblocks, engine)
	}
	return nil
}

// HandleNextBlock handles validating and syncing the next block recevied from the network
func HandleNextBlock(db *figdb.DB, chain *figaro.Chain, block *figaro.Block, engine figaro.ConsensusEngine) error {
	if !block.VerifySignature() {
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
	if !VerifyTxSignatures(block) {
		err := engine.HandleFraud(db, block.BlockHeader)
		if err != nil {
			return err
		}
		return figaro.ErrInvalidBlock
	}
	// If there's a conflict, we'll get back a new chain, block, and futureblocks and can continue as normal
	// This will also handle cleaning up invalid data from the non-canonical chain
	if !bytes.Equal(block.ParentBlock, chain.Head) {
		return figaro.ErrReorgRequired
	}
	prevbl, err := db.FetchBlock(chain.Head)
	if err != nil {
		return err
	}
	err = SyncBlock(db, prevbl, block)
	if err != nil {
		return err
	}
	err = chain.AppendBlock(db, block.BlockHeader)
	if err != nil {
		return err
	}
	return nil
}

// HandleProduceBlock handles the case where it is this nodes turn to produce a block.
func HandleProduceBlock() {

}
