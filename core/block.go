package core

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const zeroHash = "0x0000000000000000000000000000000000000000000000000000000000000000"

// Block is a collection of transactions that happen together
type Block struct {
	Index        uint64
	TimeStamp    time.Time
	Data         string
	PreviousHash string
	Hash         string
}

// NewBlock creates a new Block from data and the previous block hash
func NewBlock(index uint64, timestamp time.Time, data string, previousHash string) *Block {
	b := &Block{
		Index:        index,
		TimeStamp:    timestamp,
		Data:         data,
		PreviousHash: previousHash,
	}
	b.Seal()
	return b
}

// NewGenesisBlock creates a new genesis Block
func NewGenesisBlock() *Block {
	return NewBlock(0, time.Now(), "Genesis Block", zeroHash)
}

// Seal finalizes a block by hashing it and will return an error if previously sealed
func (b *Block) Seal() error {
	if b.Hash != "" {
		return errors.New("Block already sealed")
	}
	b.Hash = blockHash(b)
	return nil
}

// Validate generates a correct blockhash and checks it against the Block's hash
func (b *Block) Validate() bool {
	valid := blockHash(b)
	return b.Hash == valid
}

func blockHash(b *Block) string {
	h := sha256.New()
	h.Write([]byte(strconv.FormatUint(b.Index, 10)))
	h.Write([]byte(b.TimeStamp.String()))
	h.Write([]byte(b.Data))
	h.Write([]byte(b.PreviousHash))
	h.Write([]byte(b.Hash))
	return fmt.Sprintf("0x%x", h.Sum(nil))
}
