// Package figdb implements a high-performance Merkle database
package figdb

import (
	"github.com/figaro-tech/go-figaro/figdb/badger"
	"github.com/figaro-tech/go-figaro/figdb/cache"
	"github.com/figaro-tech/go-figaro/figdb/mock"
	"github.com/figaro-tech/go-figaro/figdb/set"
	"github.com/figaro-tech/go-figaro/figdb/trie"
	"github.com/figaro-tech/go-figaro/figdb/types"
)

// StateValidate validates a proof for a given root, key, and value
func StateValidate(root, key, value []byte, proof [][][]byte) bool {
	return trie.Validate(root, key, value, proof)
}

type trieDB struct {
	Set     *set.Set
	Archive *trie.Archive
	State   *trie.State
}

// FigDB is a high-performance merklized key/value database
type FigDB struct {
	trieDB
	DB *badger.KeyStore
}

// FigMemDB is a in-memory merklized key/value database,
// useful for testing or demos... state is not saved to
// disk and does not survive reboots
type FigMemDB struct {
	trieDB
	DB *mock.KeyStore
}

// New returns a FigDB, ready to use
func New(datapath string) *FigDB {
	db := badger.NewKeyStore(datapath)
	return &FigDB{
		DB: db,
		trieDB: trieDB{
			Set: &set.Set{
				KeyStore: db,
			},
			Archive: &trie.Archive{
				KeyStore: db,
			},
			State: &trie.State{
				KeyStore: db,
			},
		},
	}
}

// CacheType indicates whether to use a LRU or FIFO cache
type CacheType int

const (
	// LRU is a LRU cache
	LRU CacheType = iota
	// FIFO is a FIFO cache
	FIFO
)

// CacheConfig sets the configuration for caches
type CacheConfig struct {
	SetType     CacheType
	SetSize     int
	ArchiveType CacheType
	ArchiveSize int
	StateType   CacheType
	StateSize   int
}

// NewWithCaches returns a FigDB, ready to use.
func NewWithCaches(datapath string, config CacheConfig) *FigDB {
	db := badger.NewKeyStore(datapath)
	var setCache types.Cache
	switch config.SetType {
	case LRU:
		setCache = cache.NewLRU(config.SetSize)
	case FIFO:
		setCache = cache.NewFIFO(config.SetSize)
	}
	var archiveCache types.Cache
	switch config.ArchiveType {
	case LRU:
		archiveCache = cache.NewLRU(config.ArchiveSize)
	case FIFO:
		archiveCache = cache.NewFIFO(config.ArchiveSize)
	}
	var stateCache types.Cache
	switch config.StateType {
	case LRU:
		stateCache = cache.NewLRU(config.StateSize)
	case FIFO:
		stateCache = cache.NewFIFO(config.StateSize)
	}
	return &FigDB{
		DB: db,
		trieDB: trieDB{
			Set: &set.Set{
				KeyStore: db,
				Cache:    setCache,
			},
			Archive: &trie.Archive{
				KeyStore: db,
				Cache:    archiveCache,
			},
			State: &trie.State{
				KeyStore: db,
				Cache:    stateCache,
			},
		},
	}
}

// NewMem returns a FigMemDB, ready to use
func NewMem() *FigMemDB {
	db := mock.NewKeyStore()
	return &FigMemDB{
		DB: db,
		trieDB: trieDB{
			Set: &set.Set{
				KeyStore: db,
			},
			Archive: &trie.Archive{
				KeyStore: db,
			},
			State: &trie.State{
				KeyStore: db,
			},
		},
	}
}
