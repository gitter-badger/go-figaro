// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"github.com/figaro-tech/go-figaro/figdb"
	"github.com/figaro-tech/go-figaro/figdb/trie"
)

// DB is a domain Merkle database
type DB struct {
	DB      *figdb.FigDB
	Archive *trie.Archive
	State   *trie.State
}

// MemDB is a domain Merkle database, in-memory only
type MemDB struct {
	DB      *figdb.FigMemDB
	Archive *trie.Archive
	State   *trie.State
}

// New returns a FigDB backed by a high-performance disk database
func New(dir string) *DB {
	db := figdb.New(dir)
	return &DB{db, db.Archive, db.State}
}

// NewMem returns a FigDB backed by a high-performance memory database
func NewMem() *MemDB {
	db := figdb.NewMem()
	return &MemDB{db, db.Archive, db.State}
}
