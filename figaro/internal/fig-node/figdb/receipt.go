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
var receiptprefix = md5.Sum([]byte("figaro/receipt"))

// SaveReceipt saves a receipt underneath the associated txid.
func (db *DB) SaveReceipt(r figaro.Receipt) error {
	b, err := r.Encode()
	if err != nil {
		return err
	}
	key := hasher.Hash256(receiptprefix[:], r.TxID)
	return db.Store.Set(key, b)
}

// FetchReceipt returns a receipt associated with the txid, if it exists.
func (db *DB) FetchReceipt(txid figaro.TxHash) (r *figaro.Receipt, err error) {
	key := hasher.Hash256(receiptprefix[:], txid)
	var b []byte
	b, err = db.Store.Get(key)
	if err != nil || len(b) == 0 {
		return
	}
	err = r.Decode(b)
	if err != nil {
		return
	}
	r.TxID = txid
	return
}
