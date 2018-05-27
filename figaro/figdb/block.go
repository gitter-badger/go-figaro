// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"crypto/md5"

	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figcrypto/hasher"
)

// We prefix anything that is saved directly in the raw db, since
// the key we save under does not fully represent the data, as it
// would in archive and state tries.
var blockprefix = md5.Sum([]byte("figaro/block"))

// SaveBlockHeader saves a block hearder. It can be retreived by ID.
func (db *DB) SaveBlockHeader(header *figaro.BlockHeader) error {
	id, err := header.ID()
	if err != nil {
		return err
	}
	value, err := header.Encode()
	if err != nil {
		return err
	}
	key := hasher.Hash256(blockprefix[:], id)
	err = db.Store.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

// SaveBlock saves a block. It can be retreived by ID.
func (db *DB) SaveBlock(block *figaro.Block) error {
	id, err := block.ID()
	if err != nil {
		return err
	}
	// NOTE: this actually on saves the Header, as commits and txs
	// will have already been saved separately
	value, err := block.Encode()
	if err != nil {
		return err
	}
	key := hasher.Hash256(blockprefix[:], id)
	err = db.Store.Set(key, value)
	if err != nil {
		return err
	}
	// Save a BigBlock in the in-memory FIFO cache, so that
	// if we need to reference this block again soon (such as
	// checking for commits), we have it ready to go.
	big, err := block.Expand()
	if err != nil {
		return err
	}
	db.blockcache.Add(key, big)
	return nil
}

// FetchBlockHeader fetches just the block header by ID.
func (db *DB) FetchBlockHeader(id figaro.BlockHash) (header *figaro.BlockHeader, err error) {
	// First check the cache for a BigBlock
	// and just return its header if it exists
	key := hasher.Hash256(blockprefix[:], id)
	if big, ok := db.blockcache.Get(key); ok {
		header = big.(*figaro.BigBlock).BlockHeader
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
	if big, ok := db.blockcache.Get(key); ok {
		cblock = big.(*figaro.BigBlock).CompBlock
		return
	}
	// Get the full block and then compress it
	header, err := db.FetchBlockHeader(id)
	if err != nil {
		return
	}
	block, err := db.hydrateBlock(header)
	if err != nil {
		return
	}
	cblock, err = block.Compress()
	return
}

// FetchRefBlock returns a RefBlock, including Commits and TxIDs.
func (db *DB) FetchRefBlock(id figaro.BlockHash) (rblock *figaro.RefBlock, err error) {
	// First check the cache for a BigBlock
	key := hasher.Hash256(blockprefix[:], id)
	if big, ok := db.blockcache.Get(key); ok {
		rblock = &figaro.RefBlock{
			BlockHeader: big.(*figaro.BigBlock).BlockHeader,
			Commits:     big.(*figaro.BigBlock).Commits,
			TxIDs:       big.(*figaro.BigBlock).TxIDs,
		}
		return
	}
	// Otherwise fetch and hydrate the block
	var header *figaro.BlockHeader
	header, err = db.FetchBlockHeader(id)
	if err != nil {
		return
	}
	var block *figaro.Block
	block, err = db.hydrateBlock(header)
	if err != nil {
		return
	}
	rblock, err = block.Ref()
	return
}

// FetchBlock returns a Block, including Commits and Transactions.
func (db *DB) FetchBlock(id figaro.BlockHash) (block *figaro.Block, err error) {
	// First check the cache for a BigBlock
	key := hasher.Hash256(blockprefix[:], id)
	if big, ok := db.blockcache.Get(key); ok {
		block.BlockHeader = big.(*figaro.BigBlock).BlockHeader
		block.Commits = big.(*figaro.BigBlock).Commits
		block.Transactions = big.(*figaro.BigBlock).Transactions
		return
	}
	header, err := db.FetchBlockHeader(id)
	if err != nil {
		return nil, err
	}
	block, err = db.hydrateBlock(header)
	return
}

// FetchBigBlock returns a Block, including Commits and Transactions.
func (db *DB) FetchBigBlock(id figaro.BlockHash) (bblock *figaro.BigBlock, err error) {
	// First check the cache for a BigBlock
	key := hasher.Hash256(blockprefix[:], id)
	if big, ok := db.blockcache.Get(key); ok {
		bblock = big.(*figaro.BigBlock)
		return
	}
	// Build the BigBlock from the store
	var header *figaro.BlockHeader
	header, err = db.FetchBlockHeader(id)
	if err != nil {
		return
	}
	var block *figaro.Block
	block, err = db.hydrateBlock(header)
	if err != nil {
		return
	}
	bblock, err = block.Expand()
	return
}

func (db *DB) hydrateBlock(header *figaro.BlockHeader) (block *figaro.Block, err error) {
	block.BlockHeader = header
	block.Commits, err = db.RetrieveCommits(block.CommitsRoot)
	if err != nil {
		return
	}
	block.Transactions, err = db.RetrieveTransactions(block.TransactionsRoot)
	if err != nil {
		return
	}
	return
}
