// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figdb"
)

// ArchiveTransactions archives Transactions, returning the merkle root of the archive.
func (db *DB) ArchiveTransactions(txs []*figaro.Transaction) (root figaro.Root, err error) {
	encoded := make([][]byte, len(txs))
	var e []byte
	for i, tx := range txs {
		e, err = tx.Encode()
		if err != nil {
			return
		}
		encoded[i] = e
	}
	root, err = db.Archive.Save(encoded)
	return
}

// RetrieveTransactions retrieves an archive of Transactions from a merkle root.
func (db *DB) RetrieveTransactions(root figaro.Root) (txs []*figaro.Transaction, err error) {
	var encoded [][]byte
	encoded, err = db.Archive.Retrieve(root)
	if err != nil {
		return
	}
	txs = make([]*figaro.Transaction, len(encoded))
	for i, e := range encoded {
		tx := &figaro.Transaction{}
		err = tx.Decode(e)
		if err != nil {
			return
		}
		txs[i] = tx
	}
	return
}

// GetTransaction gets the Transaction at index in from the archive in the merkle root.
func (db *DB) GetTransaction(root figaro.Root, index int) (tx *figaro.Transaction, err error) {
	var e []byte
	e, err = db.Archive.Get(root, index)
	if err != nil || len(e) == 0 {
		return
	}
	err = tx.Decode(e)
	return
}

// GetAndProveTransaction gets the Transaction at index in from the archive in the merkle root, providing a merkle proof.
func (db *DB) GetAndProveTransaction(root figaro.Root, index int) (tx *figaro.Transaction, proof [][]byte, err error) {
	var e []byte
	e, proof, err = db.Archive.GetAndProve(root, index)
	if err != nil || len(e) == 0 {
		return
	}
	err = tx.Decode(e)
	return
}

// ValidateTransaction validates whether a proof is valid for a given Transaction in root at index.
func (db *DB) ValidateTransaction(root figaro.Root, index int, tx figaro.Transaction, proof [][]byte) bool {
	e, err := tx.Encode()
	if err != nil {
		return false
	}
	return figdb.ValidateArchive(root, index, e, proof)
}
