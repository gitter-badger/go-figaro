// Package figaro is the main package for go-figaro
package figaro

import (
	"errors"

	"github.com/figaro-tech/go-fig-buf"
)

var (
	// ErrInvalidBlock is a self-explantory error
	ErrInvalidBlock = errors.New("figaro chain: invalid block for chain")
	// ErrReorgRequired is a self-explantory error
	ErrReorgRequired = errors.New("figaro chain: chain reorg required")
)

// ChainConfig represents the current config for the chain. It will be saved in each
// block header for future reference.
type ChainConfig struct {
	Stake      uint64
	CommitFee  uint32
	TxFee      uint32
	WaitBlocks uint8
}

// Encode deterministically encodes a Chain to binary format.
func (cc ChainConfig) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextUint64(buf, cc.Stake)
		buf = enc.EncodeNextUint32(buf, cc.CommitFee)
		buf = enc.EncodeNextUint32(buf, cc.TxFee)
		buf = enc.EncodeNextUint8(buf, cc.WaitBlocks)
		return buf
	})
}

// Decode decodes a deterministically encoded Chain from binary format.
func (cc *ChainConfig) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	return dec.DecodeList(buf, func(r []byte) []byte {
		cc.Stake, r = dec.DecodeNextUint64(r)
		cc.CommitFee, r = dec.DecodeNextUint32(r)
		cc.TxFee, r = dec.DecodeNextUint32(r)
		cc.WaitBlocks, r = dec.DecodeNextUint8(r)
		return r
	})
}

// Chain is a singly-linked list where each block in the chain links
// to the previous block in the chain via cryptographically secure IDs.
// There can be only one canononical chain.
type Chain struct {
	Depth uint64
	Head  BlockHash
	ChainConfig
}

// NextBlock generates a new block that is the child of
// the current block, ready for processing.
func (chain *Chain) NextBlock() *Block {
	return &Block{
		BlockHeader: &BlockHeader{
			Number:      chain.Depth + 1,
			ParentBlock: chain.Head,
			ChainConfig: chain.ChainConfig,
		},
	}
}

// AppendBlock adds a block to the chain, saving
// the chain in the database. It will not verify that the
// chain is not gapped, only that the block number is greater
// than the current chain depth.
func (chain *Chain) AppendBlock(db ChainDataService, bl *BlockHeader) error {
	if bl.Number <= chain.Depth {
		return ErrInvalidBlock
	}
	chain.Head = bl.ID
	chain.Depth = bl.Number
	return db.SaveChain(chain)
}

// Encode deterministically encodes a Chain to binary format.
func (chain Chain) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextBytes(buf, chain.Head)
		buf = enc.EncodeNextUint64(buf, chain.Depth)
		cfg, err := chain.ChainConfig.Encode()
		if err != nil {
			panic("could not encode chain config")
		}
		buf = enc.EncodeNextBytes(buf, cfg)
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

		var cfg []byte
		cfg, r = dec.DecodeNextBytes(r)
		err := chain.ChainConfig.Decode(cfg)
		if err != nil {
			panic("invalid data for chain config")
		}
		return r
	})
}

// ChainDataService should save chain directly into a key/value store.
type ChainDataService interface {
	// SaveChain should the chain itself, as well as
	// a reference to the current block head by number.
	SaveChain(*Chain) error
	FetchChain() (*Chain, error)
	FetchChainBlock(index uint64) (BlockHash, error)
}
