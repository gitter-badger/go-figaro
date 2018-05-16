// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figcrypto/trie"
)

// ArchiveTxCommits archives commits, returning the merkle root of the archive.
func (db *DB) ArchiveTxCommits(ed figaro.TransactionEncodingService, commits ...*figaro.TxCommit) ([]byte, error) {
	encoded := make([][]byte, len(commits))
	for i, c := range commits {
		e, err := ed.EncodeTxCommit(c)
		if err != nil {
			return nil, err
		}
		encoded[i] = e
	}
	return db.Archive.Save(encoded)
}

// RetrieveTxCommits retrieves an archive of commits from a merkle root.
func (db *DB) RetrieveTxCommits(ed figaro.TransactionEncodingService, root []byte) ([]*figaro.TxCommit, error) {
	encoded, err := db.Archive.Retrieve(root)
	if err != nil {
		return nil, err
	}
	commits := make([]*figaro.TxCommit, len(encoded))
	for i, e := range encoded {
		c, err := ed.DecodeTxCommit(e)
		if err != nil {
			return nil, err
		}
		commits[i] = c
	}
	return commits, nil
}

// GetTxCommit gets the commit at index in from the archive in the merkle root.
func (db *DB) GetTxCommit(ed figaro.TransactionEncodingService, root []byte, index int) (*figaro.TxCommit, error) {
	e, err := db.Archive.Get(root, int(index))
	if err != nil {
		return nil, err
	}
	if len(e) == 0 {
		return nil, nil
	}
	return ed.DecodeTxCommit(e)
}

// GetAndProveTxCommit gets the commit at index in from the archive in the merkle root, providing a merkle proof.
func (db *DB) GetAndProveTxCommit(ed figaro.TransactionEncodingService, root []byte, index int) (*figaro.TxCommit, [][]byte, error) {
	e, p, err := db.Archive.GetAndProve(root, int(index))
	if err != nil {
		return nil, nil, err
	}
	if len(e) == 0 {
		return nil, nil, nil
	}
	c, err := ed.DecodeTxCommit(e)
	if err != nil {
		return nil, nil, err
	}
	return c, p, nil
}

// ValidateTxCommit validates whether a proof is valid for a given commit in root at index.
func (db *DB) ValidateTxCommit(ed figaro.TransactionEncodingService, root []byte, index int, commit *figaro.TxCommit, proof [][]byte) bool {
	e, err := ed.EncodeTxCommit(commit)
	if err != nil {
		return false
	}
	return trie.Validate(root, index, e, proof)
}

// ArchiveTransactions archives transactions, returning the merkle root of the archive.
func (db *DB) ArchiveTransactions(ed figaro.TransactionEncodingService, txs ...*figaro.Transaction) ([]byte, error) {
	encoded := make([][]byte, len(txs))
	for i, c := range txs {
		e, err := ed.EncodeTransaction(c)
		if err != nil {
			return nil, err
		}
		encoded[i] = e
	}
	return db.Archive.Save(encoded)
}

// RetrieveTransactions retrieves an archive of transactions from a merkle root.
func (db *DB) RetrieveTransactions(ed figaro.TransactionEncodingService, root []byte) ([]*figaro.Transaction, error) {
	encoded, err := db.Archive.Retrieve(root)
	if err != nil {
		return nil, err
	}
	commits := make([]*figaro.Transaction, len(encoded))
	for i, e := range encoded {
		c, err := ed.DecodeTransaction(e)
		if err != nil {
			return nil, err
		}
		commits[i] = c
	}
	return commits, nil
}

// GetTransaction gets the transaction at index in from the archive in the merkle root.
func (db *DB) GetTransaction(ed figaro.TransactionEncodingService, root []byte, index int) (*figaro.Transaction, error) {
	e, err := db.Archive.Get(root, int(index))
	if err != nil {
		return nil, err
	}
	if len(e) == 0 {
		return nil, nil
	}
	return ed.DecodeTransaction(e)
}

// GetAndProveTransaction gets the transaction at index in from the archive in the merkle root, providing a merkle proof.
func (db *DB) GetAndProveTransaction(ed figaro.TransactionEncodingService, root []byte, index int) (*figaro.Transaction, [][]byte, error) {
	e, p, err := db.Archive.GetAndProve(root, int(index))
	if err != nil {
		return nil, nil, err
	}
	if len(e) == 0 {
		return nil, nil, nil
	}
	c, err := ed.DecodeTransaction(e)
	if err != nil {
		return nil, nil, err
	}
	return c, p, nil
}

// ValidateTransaction validates whether a proof is valid for a given transaction in root at index.
func (db *DB) ValidateTransaction(ed figaro.TransactionEncodingService, root []byte, index int, tx *figaro.Transaction, proof [][]byte) bool {
	e, err := ed.EncodeTransaction(tx)
	if err != nil {
		return false
	}
	return trie.Validate(root, index, e, proof)
}
