package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

// Block is a collection of transactions that happen together
type Block struct {
	Index        uint64
	TimeStamp    time.Time
	Data         string
	PreviousHash string
	Hash         string
}

// Log prints out information about a Block
func (b *Block) Log() {
	fmt.Println("\n=================================================")
	fmt.Println(fmt.Sprintf("\nIndex: %d", b.Index))
	fmt.Println(fmt.Sprintf("TimeStamp: %s", b.TimeStamp))
	fmt.Println(fmt.Sprintf("Data: %s", b.Data))
	fmt.Println(fmt.Sprintf("PreviousHash: %s", b.PreviousHash))
	fmt.Println(fmt.Sprintf("Hash: %s", b.Hash))
}

// ToHash returns a hashed representation of the Block
func (b *Block) ToHash() string {
	h := sha256.New()
	h.Write([]byte(strconv.FormatUint(b.Index, 10)))
	h.Write([]byte(b.TimeStamp.String()))
	h.Write([]byte(b.Data))
	h.Write([]byte(b.PreviousHash))
	h.Write([]byte(b.Hash))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// NewBlock creates a new Block and immediately hashes it
func NewBlock(index uint64, timestamp time.Time, data string, previousHash string) *Block {
	b := Block{
		Index:        index,
		TimeStamp:    timestamp,
		Data:         data,
		PreviousHash: previousHash,
	}
	b.Hash = b.ToHash()
	return &b
}

// CreateGenesisBlock creates a Block to start a chain
func CreateGenesisBlock() *Block {
	return NewBlock(0, time.Now(), "Genesis Block", "0")
}

// NextBlock creates a new Block from an existing Block
func NextBlock(lastBlock *Block) *Block {
	return NewBlock(
		lastBlock.Index+1,
		time.Now(),
		fmt.Sprintf("This is block number %d", lastBlock.Index+1),
		lastBlock.Hash)
}

func main() {
	chain := []*Block{CreateGenesisBlock()}
	previousBlock := chain[0]

	for i := 1; i <= 20; i++ {
		newBlock := NextBlock(previousBlock)
		chain = append(chain, newBlock)
		previousBlock = chain[i]
		newBlock.Log()
	}
}
