package core

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

// This should get moved to a crypto package in the future
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

// Seal finalizes a block by hashing it and will return an error if previously sealed
func (b *Block) Seal() *Block {
	b.Hash = blockHash(b)
	// In the future, this will probably need to do some additional work
	// like cryptographically signing the block
	return b
}

// Valid generates a correct blockhash and checks it against the Block's hash
func (b *Block) Valid() bool {
	return b.Hash == blockHash(b)
}

func blockHash(b *Block) string {
	h := sha256.New()
	h.Write([]byte(strconv.FormatUint(b.Index, 10)))
	h.Write([]byte(b.TimeStamp.String()))
	h.Write([]byte(b.Data))
	h.Write([]byte(b.PreviousHash))
	return fmt.Sprintf("0x%x", h.Sum(nil))
}
