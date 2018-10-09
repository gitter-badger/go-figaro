// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	fdb "github.com/figaro-tech/go-fig-db"
	"github.com/figaro-tech/go-figaro/figaro"
)

// ArchiveCommits archives Commits, returning the merkle root of the archive.
func (db *DB) ArchiveCommits(commits []figaro.Commit) (root figaro.Root, err error) {
	b := make([][]byte, len(commits))
	for i, c := range commits {
		b[i] = c
	}
	root, err = db.Archive.Save(b)
	if err != nil {
		return
	}
	return
}

// RetrieveCommits retrieves an archive of Commits from a merkle root.
func (db *DB) RetrieveCommits(root figaro.Root) (commits []figaro.Commit, err error) {
	var bb [][]byte
	bb, err = db.Archive.Retrieve(root)
	if err != nil {
		return
	}
	commits = make([]figaro.Commit, len(bb))
	for i, b := range bb {
		commits[i] = b
	}
	return
}

// GetCommit gets the Commit at index in from the archive in the merkle root.
func (db *DB) GetCommit(root figaro.Root, index int) (figaro.Commit, error) {
	return db.Archive.Get(root, index)
}

// GetAndProveCommit gets the Commit at index in from the archive in the merkle root, providing a merkle proof.
func (db *DB) GetAndProveCommit(root figaro.Root, index int) (Commits figaro.Commit, proof [][]byte, err error) {
	return db.Archive.GetAndProve(root, index)
}

// ValidateCommit validates whether a proof is valid for a given Commit in root at index.
func (db *DB) ValidateCommit(root figaro.Root, index int, commit figaro.Commit, proof [][]byte) bool {
	return fdb.ValidateArchive(root, index, commit, proof)
}
