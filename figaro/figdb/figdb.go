// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"github.com/figaro-tech/go-figaro/figdb"
	"github.com/figaro-tech/go-figaro/figdb/cache"
)

// DB is a domain Merkle database
type DB struct {
	*figdb.FigDB
	blockcache *cache.FIFO
}

// New returns a FigDB backed by a high-performance disk database.
func New(dir string, blockcachesize int) *DB {
	db := figdb.New(dir)
	return &DB{db, cache.NewFIFO(blockcachesize)}
}

// NewMem returns a FigDB backed by a high-performance memory database.
func NewMem(dir, blockcachesize int) *DB {
	db := figdb.NewMem()
	return &DB{db, cache.NewFIFO(blockcachesize)}
}
