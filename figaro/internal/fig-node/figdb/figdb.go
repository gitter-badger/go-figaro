// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	fdb "github.com/figaro-tech/go-fig-db"
	"github.com/figaro-tech/go-fig-db/cache"
)

// DB is a domain Merkle database
type DB struct {
	*fdb.FigDB
	blockcache *cache.FIFO
}

// New returns a FigDB backed by a high-performance disk database.
func New(dir string, blockcachesize int) *DB {
	db := fdb.New(dir)
	return &DB{db, cache.NewFIFO(blockcachesize)}
}

// NewMem returns a FigDB backed by a high-performance memory database.
func NewMem(dir, blockcachesize int) *DB {
	db := fdb.NewMem()
	return &DB{db, cache.NewFIFO(blockcachesize)}
}
