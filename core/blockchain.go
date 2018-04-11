package core

import (
	"errors"
	"time"
)

// BlockChain is a collection of ordered blocks
type BlockChain struct {
	Blocks []*Block
}

// NewBlockChain will create a new Blockchain given a genesis block
func NewBlockChain(genesis *Block) *BlockChain {
	return &BlockChain{Blocks: []*Block{genesis}}
}

// Head returns the head of the blockchain (most recent block)
func (c *BlockChain) Head() *Block {
	return c.Blocks[len(c.Blocks)-1]
}

// CreateBlock will create a new Block from data and add it to the chain
func (c *BlockChain) CreateBlock(data string) *Block {
	h := c.Head()
	n := NewBlock(h.Index+1, time.Now(), data, h.Hash)
	c.Blocks = append(c.Blocks, n)
	return n
}

// ReceiveBlock will validate a Block and add it to chain
func (c *BlockChain) ReceiveBlock(block *Block) error {
	h := c.Head()
	// In the real world, blocks could be received out of order
	// so we'll need to revisit this at some point
	if block.Index != h.Index+1 {
		return errors.New("Block received out of order")
	}
	// Validate the previous hash is the current head
	if block.PreviousHash != h.Hash {
		return errors.New("Invalid previous hash")
	}
	// Validate that the current hash
	if !block.Validate() {
		return errors.New("Invalid hash")
	}
	c.Blocks = append(c.Blocks, block)
	return nil
}
