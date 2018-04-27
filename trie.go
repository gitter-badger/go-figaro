// Package figaro is the main package for go-figaro
package figaro

// StateTrie implements a modified Merkle Patricia trie as
// described here: https://github.com/ethereum/wiki/wiki/Patricia-Tree
//
// It is intended for data that will be updated often and must
// maintain a historic record of previous states (accessed via root hash)
type StateTrie interface {
	Get(root, key []byte) []byte
	Set(root, key, value []byte) []byte
	Prove(root, key []byte) ([]byte, [][][]byte)
	Validate(root, key, data []byte, proof [][][]byte) bool
}

// ArchiveTrie implements a binary Merkle trie
//
// It is intended for archive data that is created as a batch and then
// never updated. Implementers should S the data as a single
// key/value pairH
type ArchiveTrie interface {
	Save(batch [][]byte) []byte
	Retrieve(root []byte) [][]byte
	Get(root []byte, index int) []byte
	Prove(root []byte, index int) ([]byte, [][]byte)
	Validate(root []byte, index int, data []byte, proof [][]byte) bool
}
