// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figdb"
)

const txCommitFp = 0.0000001 // 1 in 10 million

// CreateTxCommits saves the commits, returning the key and set in binary format.
func (db *DB) CreateTxCommits(commits ...figaro.TxCommit) (key []byte, set []byte, err error) {
	bb := make([][]byte, len(commits))
	for i, c := range commits {
		bb[i] = c
	}
	key, set, err = db.Set.Create(bb, txCommitFp)
	return
}

// HasTxCommits determines whether the set at key contains each commit.
func (db *DB) HasTxCommits(root []byte, commits ...figaro.TxCommit) []bool {
	bb := make([][]byte, len(commits))
	for i, c := range commits {
		bb[i] = c
	}
	return db.Set.HasBatch(root, bb)
}

// HasTxCommit determines whether the set at key contains a commit.
func (db *DB) HasTxCommit(root []byte, commit figaro.TxCommit) bool {
	return db.Set.Has(root, commit)
}

// ArchiveTxs archives Txs, returning the merkle root of the archive.
func (db *DB) ArchiveTxs(txs ...*figaro.Tx) (root []byte, err error) {
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

// RetrieveTxs retrieves an archive of Txs from a merkle root.
func (db *DB) RetrieveTxs(root []byte) (txs []*figaro.Tx, err error) {
	var encoded [][]byte
	encoded, err = db.Archive.Retrieve(root)
	if err != nil {
		return
	}
	txs = make([]*figaro.Tx, len(encoded))
	for i, e := range encoded {
		tx := &figaro.Tx{}
		err = tx.Decode(e)
		if err != nil {
			return
		}
		txs[i] = tx
	}
	return
}

// GetTx gets the Tx at index in from the archive in the merkle root.
func (db *DB) GetTx(root []byte, index int) (tx *figaro.Tx, err error) {
	var e []byte
	e, err = db.Archive.Get(root, index)
	if err != nil || len(e) == 0 {
		return
	}
	tx = &figaro.Tx{}
	tx.Decode(e)
	return
}

// GetAndProveTx gets the Tx at index in from the archive in the merkle root, providing a merkle proof.
func (db *DB) GetAndProveTx(root []byte, index int) (tx *figaro.Tx, proof [][]byte, err error) {
	var e []byte
	e, proof, err = db.Archive.GetAndProve(root, index)
	if err != nil || len(e) == 0 {
		return
	}
	tx = &figaro.Tx{}
	tx.Decode(e)
	return
}

// ValidateTx validates whether a proof is valid for a given Tx in root at index.
func (db *DB) ValidateTx(root []byte, index int, tx *figaro.Tx, proof [][]byte) bool {
	e, err := tx.Encode()
	if err != nil {
		return false
	}
	return figdb.ValidateArchive(root, index, e, proof)
}
