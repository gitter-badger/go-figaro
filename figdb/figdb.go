// Package figdb implements a high-performance Merkle database
package figdb

import (
	"github.com/figaro-tech/go-figaro/figdb/badger"
	"github.com/figaro-tech/go-figaro/figdb/mock"
	"github.com/figaro-tech/go-figaro/figdb/set"
	"github.com/figaro-tech/go-figaro/figdb/trie"
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
