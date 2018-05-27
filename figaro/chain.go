// Package figaro is the main package for go-figaro
package figaro

import (
	"bytes"
	"container/heap"
	"errors"
	"reflect"

	"github.com/figaro-tech/go-figaro/figbuf"
)

// ErrInvalidBlock is returned if appending a produced block would break the chain.
var ErrInvalidBlock = errors.New("figaro chain: invalid next block for chain")

// ChainConfig represents the current config for the chain. It will be saved in each
// block header for future reference.
type ChainConfig struct {
	Authority  Address
	Stake      uint64
	CommitFee  uint32
	TxFee      uint32
	WaitBlocks uint8
}

// Chain is a singly-linked list where each block in the chain links
// to the previous block in the chain via cryptographically secure IDs.
// There can be only one canononical chain.
type Chain struct {
	Depth uint64
	Head  BlockHash
	ChainConfig
}

// NewBlock creates a new block based on the current head,
// ready for adding commits and transactions.
func (chain *Chain) NewBlock(provider Address, beneficiary Address) *Block {
	return &Block{
		BlockHeader: &BlockHeader{
			Provider:    provider,
			Beneficiary: beneficiary,
			Number:      chain.Depth + 1,
			ParentBlock: chain.Head,
			ChainConfig: chain.ChainConfig,
		},
	}
}

// AppendBlock will append a block to the head, saving both the block
// and the chain. It performs only basic checks to ensure the chain is
// not gapped, and assumes that the block is produced locally
func (chain *Chain) AppendBlock(db FullChainDataService, block *Block) error {
	if !block.WellFormed() {
		return ErrInvalidBlock
	}
	if !block.VerifySignature() {
		return ErrInvalidBlock
	}
	if reflect.DeepEqual(block.ChainConfig, chain.ChainConfig) {
		return ErrInvalidBlock
	}
	if block.Number > chain.Depth+1 {
		return ErrInvalidBlock
	}
	if !bytes.Equal(block.ParentBlock, chain.Head) {
		return ErrInvalidBlock
	}
	if reflect.DeepEqual(block.ChainConfig, chain.ChainConfig) {
		return ErrInvalidBlock
	}
	err := db.SaveBlock(block)
	if err != nil {
		return err
	}
	chain.Depth++
	chain.Head, err = block.ID()
	if err != nil {
		return err
	}
	db.SaveChain(chain)
	return nil
}

// ReceiveBlock processes a block received from the network. It performs only
// basic checks, and assumes that the block is trusted if the signature and
// producer are valid. If the block is the next block in the chain, it will be
// appended. If the chain is gapped, the block will be added to `futureblocks`.
// If there is a conflict, `engine` will be used to initiate a chain reorg.
func (chain *Chain) ReceiveBlock(db FullChainDataService, block *Block, futureblocks *BlockHeap, engine ConsensusEngine) error {
	if !block.WellFormed() {
		return nil
	}
	if !block.VerifySignature() {
		return nil
	}
	if reflect.DeepEqual(block.ChainConfig, chain.ChainConfig) {
		return nil
	}
	if block.Number > chain.Depth+1 {
		heap.Push(futureblocks, block)
		return nil
	}
	if !engine.ValidateBlockProducer(db, block) {
		return nil
	}
	// If there's a conflict, we'll get back a new chain, block, and futureblocks and can continue as normal
	if !bytes.Equal(block.ParentBlock, chain.Head) {
		var err error
		*chain, *block, *futureblocks, err = engine.HandleDivergence(db, *chain, block, futureblocks)
		if err != nil {
			return err
		}
	}
	return chain.AppendBlock(db, block)
}

// ReceiveAndSyncBlock processes a block received from the network.If the block
// is the next block in the chain, it will be fully validated and synced. If
// the chain is gapped, the block will be added to `futureblocks`. If there is a
// conflict, `engine` will be used to initiate a chain reorg.
func (chain *Chain) ReceiveAndSyncBlock(db FullChainDataService, block *Block, futureblocks *BlockHeap, engine ConsensusEngine) error {
	if !block.WellFormed() {
		return nil
	}
	if !block.VerifySignature() {
		return nil
	}
	if reflect.DeepEqual(block.ChainConfig, chain.ChainConfig) {
		return nil
	}
	if block.Number > chain.Depth+1 {
		heap.Push(futureblocks, block)
		return nil
	}
	if !engine.ValidateBlockProducer(db, block) {
		return nil
	}
	if !engine.ValidateBlockTxs(db, block) {
		return nil
	}
	if !block.Sync(db) {
		err := engine.HandleInvalidHeaders(db, block)
		if err != nil {
			return err
		}
		return nil
	}
	// If there's a conflict, we'll get back a new chain, block, and futureblocks and can continue as normal
	if !bytes.Equal(block.ParentBlock, chain.Head) {
		var err error
		*chain, *block, *futureblocks, err = engine.HandleDivergence(db, *chain, block, futureblocks)
		if err != nil {
			return err
		}
	}
	return chain.AppendBlock(db, block)
}

// Encode deterministically encodes a Chain to binary format.
func (chain Chain) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextBytes(buf, chain.Head)
		buf = enc.EncodeNextUint64(buf, chain.Depth)
		buf = enc.EncodeNextBytes(buf, chain.Authority)
		buf = enc.EncodeNextUint64(buf, chain.Stake)
		buf = enc.EncodeNextUint32(buf, chain.CommitFee)
		buf = enc.EncodeNextUint32(buf, chain.TxFee)
		buf = enc.EncodeNextUint8(buf, chain.WaitBlocks)
		return buf
	})
}

// Decode decodes a deterministically encoded Chain from binary format.
func (chain *Chain) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	return dec.DecodeList(buf, func(r []byte) []byte {
		chain.Head, r = dec.DecodeNextBytes(r)
		chain.Depth, r = dec.DecodeNextUint64(r)
		chain.Authority, r = dec.DecodeNextBytes(r)
		chain.Stake, r = dec.DecodeNextUint64(r)
		chain.CommitFee, r = dec.DecodeNextUint32(r)
		chain.TxFee, r = dec.DecodeNextUint32(r)
		chain.WaitBlocks, r = dec.DecodeNextUint8(r)
		return r
	})
}

// FullChainDataService provides full data for all chain types.
type FullChainDataService interface {
	AccountDataService
	CommitDataService
	TransactionDataService
	ReceiptDataService
	BlockDataService
	ChainDataService
}

// ChainDataService should save chain directly into a key/value store.
type ChainDataService interface {
	SaveChain(*Chain) error
	FetchChain() (*Chain, error)
	FetchChainBlock(index uint64) (BlockHash, error)
}
