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

// ValidateSet is a convenience wrapper for set.ValidateSet
func ValidateSet(key, data []byte) bool {
	return set.Validate(key, data)
}

// ValidateArchive is a convenience wrapper for trie.ValidateBMT
func ValidateArchive(root []byte, index int, data []byte, proof [][]byte) bool {
	return trie.ValidateBMT(root, index, data, proof)
}

// ValidateState is a convenience wrapper for trie.ValidateMPT
func ValidateState(root, key, value []byte, proof [][][]byte) bool {
	return trie.ValidateMPT(root, key, value, proof)
}

// FigDB is a high-performance merklized key/value database
type FigDB struct {
	Store   types.KeyStore
	Set     *set.Set
	Archive *trie.Archive
	State   *trie.State
}

// FigRawDB is a high-performance key/value database
type FigRawDB struct {
	Store types.KeyStore
}

// New returns a FigDB, ready to use
func New(datapath string) *FigDB {
	db := badger.NewKeyStore(datapath, nil)
	return &FigDB{
		Store: db,
		Set: &set.Set{
			KeyStore: db,
		},
		Archive: &trie.Archive{
			KeyStore: db,
		},
		State: &trie.State{
			KeyStore: db,
		},
	}
}

// NewRaw returns a FigRawDB, ready to use
func NewRaw(datapath string) *FigRawDB {
	db := badger.NewKeyStore(datapath, nil)
	return &FigRawDB{
		Store: db,
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
	RawType     CacheType
	RawSize     int
	SetType     CacheType
	SetSize     int
	ArchiveType CacheType
	ArchiveSize int
	StateType   CacheType
	StateSize   int
}

// NewWithCaches returns a FigDB, ready to use.
func NewWithCaches(datapath string, config CacheConfig) *FigDB {
	var rawCache types.Cache
	switch config.RawType {
	case LRU:
		rawCache = cache.NewBytesLRU(config.RawSize)
	case FIFO:
		rawCache = cache.NewBytesFIFO(config.RawSize)
	}
	db := badger.NewKeyStore(datapath, rawCache)
	var setCache types.Cache
	switch config.SetType {
	case LRU:
		setCache = cache.NewBytesLRU(config.SetSize)
	case FIFO:
		setCache = cache.NewBytesFIFO(config.SetSize)
	}
	var archiveCache types.Cache
	switch config.ArchiveType {
	case LRU:
		archiveCache = cache.NewBytesLRU(config.ArchiveSize)
	case FIFO:
		archiveCache = cache.NewBytesFIFO(config.ArchiveSize)
	}
	var stateCache types.Cache
	switch config.StateType {
	case LRU:
		stateCache = cache.NewBytesLRU(config.StateSize)
	case FIFO:
		stateCache = cache.NewBytesFIFO(config.StateSize)
	}
	return &FigDB{
		Store: db,
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
	}
}

// NewRawWithCache returns a FigRawDB, ready to use.
func NewRawWithCache(datapath string, config CacheConfig) *FigRawDB {
	var rawCache types.Cache
	switch config.RawType {
	case LRU:
		rawCache = cache.NewBytesLRU(config.RawSize)
	case FIFO:
		rawCache = cache.NewBytesFIFO(config.RawSize)
	}
	db := badger.NewKeyStore(datapath, rawCache)
	return &FigRawDB{
		Store: db,
	}
}

// NewMem returns a FigMemDB, ready to use
func NewMem() *FigDB {
	db := mock.NewKeyStore()
	return &FigDB{
		Store: db,
		Set: &set.Set{
			KeyStore: db,
		},
		Archive: &trie.Archive{
			KeyStore: db,
		},
		State: &trie.State{
			KeyStore: db,
		},
	}
}

// NewRawMem returns a FigMemDB, ready to use
func NewRawMem() *FigDB {
	db := mock.NewKeyStore()
	return &FigDB{
		Store: db,
	}
}
