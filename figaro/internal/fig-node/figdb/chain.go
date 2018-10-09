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
var (
	chainprefix = md5.Sum([]byte("figaro/chain"))
	chainhead   = hasher.Hash256([]byte("figaro/chainhead"))
)

// SaveChain saves the canonical chain.
func (db *DB) SaveChain(chain *figaro.Chain) error {
	b, err := chain.Encode()
	if err != nil {
		return err
	}
	err = db.Store.Set(chainhead, b)
	if err != nil {
		return err
	}
	key := hasher.Hash256(chainprefix[:], []byte(string(chain.Depth)))
	return db.Store.Set(key, chain.Head)
}

// FetchChain fetches the canonical chain.
func (db *DB) FetchChain() (chain *figaro.Chain, err error) {
	var b []byte
	b, err = db.Store.Get(chainhead)
	if err != nil || len(b) == 0 {
		return
	}
	chain.Decode(b)
	return
}

// FetchChainBlock fetches the Block at index in the canonical chain.
func (db *DB) FetchChainBlock(index uint64) (bhash figaro.BlockHash, err error) {
	key := hasher.Hash256(chainprefix[:], []byte(string(index)))
	bhash, err = db.Store.Get(key)
	return
}
