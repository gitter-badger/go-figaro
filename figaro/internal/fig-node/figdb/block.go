// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"crypto/md5"

	"github.com/figaro-tech/go-fig-crypto/hasher"
	"github.com/figaro-tech/go-figaro/figaro"
)

// We prefix anything that is saved directly in the raw db, since
// the key we save under does not fully represent the data, as it
// would in archive and state tries.
var blockprefix = md5.Sum([]byte("figaro/block"))

type blockCacheItem struct {
	header figaro.BlockHeader
	block  figaro.Block
	ref    figaro.RefBlock
	comp   figaro.CompBlock
}

// SaveBlock saves a block. It can be retreived by ID.
func (db *DB) SaveBlock(block *figaro.Block) error {
	// We only save the header, everything else is saved separately
	value, err := block.BlockHeader.Encode()
	if err != nil {
		return err
	}
	key := hasher.Hash256(blockprefix[:], block.ID)
	err = db.Store.Set(key, value)
	if err != nil {
		return err
	}
	ref := block.Ref()
	comp := block.Compress()
	item := &blockCacheItem{
		header: *(block.BlockHeader),
		block:  *block,
		ref:    *ref,
		comp:   *comp,
	}
	db.blockcache.Add(key, item)
	return nil
}

// FetchBlockHeader fetches just the block header by ID.
func (db *DB) FetchBlockHeader(id figaro.BlockHash) (header *figaro.BlockHeader, err error) {
	// First check the cache for a BigBlock
	// and just return its header if it exists
	key := hasher.Hash256(blockprefix[:], id)
	if item, ok := db.blockcache.Get(key); ok {
		*header = item.(*blockCacheItem).header
		return
	}
	// Fetch the header from the store
	var b []byte
	b, err = db.Store.Get(key)
	if err != nil {
		return
	}
	err = header.Decode(b)
	if err != nil {
		return
	}
	return
}

// FetchCompBlock returns a CompBlock, including CommitsBloom and TxBloom.
func (db *DB) FetchCompBlock(id figaro.BlockHash) (cblock *figaro.CompBlock, err error) {
	// First check the cache for a BigBlock
	key := hasher.Hash256(blockprefix[:], id)
	if item, ok := db.blockcache.Get(key); ok {
		*cblock = item.(*blockCacheItem).comp
		return
	}
	// Get the full block and then compress it
	header, err := db.FetchBlockHeader(id)
	if err != nil {
		return
	}
	block, err := db.HydrateBlock(header)
	if err != nil {
		return
	}
	cblock = block.Compress()
	return
}

// FetchRefBlock returns a RefBlock, including Commits and TxIDs.
func (db *DB) FetchRefBlock(id figaro.BlockHash) (rblock *figaro.RefBlock, err error) {
	// First check the cache for a BigBlock
	key := hasher.Hash256(blockprefix[:], id)
	if item, ok := db.blockcache.Get(key); ok {
		*rblock = item.(*blockCacheItem).ref
		return
	}
	// Otherwise fetch and hydrate the block
	var header *figaro.BlockHeader
	header, err = db.FetchBlockHeader(id)
	if err != nil {
		return
	}
	var block *figaro.Block
	block, err = db.HydrateBlock(header)
	if err != nil {
		return
	}
	rblock = block.Ref()
	return
}

// FetchBlock returns a Block, including Commits and Transactions.
func (db *DB) FetchBlock(id figaro.BlockHash) (block *figaro.Block, err error) {
	// First check the cache for a BigBlock
	key := hasher.Hash256(blockprefix[:], id)
	if item, ok := db.blockcache.Get(key); ok {
		*block = item.(*blockCacheItem).block
		return
	}
	header, err := db.FetchBlockHeader(id)
	if err != nil {
		return nil, err
	}
	block, err = db.HydrateBlock(header)
	return
}

// HydrateBlock creates a block from a BlockHeader by retreiving missing
// data from the database.
func (db *DB) HydrateBlock(header *figaro.BlockHeader) (block *figaro.Block, err error) {
	block.BlockHeader = header
	block.Commits, err = db.RetrieveCommits(block.CommitsRoot)
	if err != nil {
		return
	}
	block.Transactions, err = db.RetrieveTransactions(block.TransactionsRoot)
	if err != nil {
		return
	}
	err = block.SetBlooms()
	return
}
