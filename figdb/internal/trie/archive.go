// Package trie implements the figaro.*Trie interfaces
package trie

import (
	"bytes"
	"log"

	"github.com/figaro-tech/go-figaro/figdb/internal"
)

// Archive impelements a pkg.ArchiveTrie
type Archive struct {
	KeyStore internal.KeyStore
	Hasher   internal.Hasher
	Encdec   internal.EncoderDecoder
}

// ArchiveValidator implements a pkg.ArchiveValidator
type ArchiveValidator struct {
	Hasher internal.Hasher
}

// Save creates a new entry for the batch of data
// and returns a Merkle root
func (tr *Archive) Save(batch [][]byte) []byte {
	// construct a root hash, end then store the encoded data at the root
	// if data has an odd length, append nil to make it even
	if len(batch)&1 == 1 {
		batch = append(batch, nil)
	}
	// hash the first layer
	level := make([][]byte, 0, len(batch))
	for _, d := range batch {
		level = append(level, tr.Hasher.Hash(d))
	}
	// build up the trie until we have a root
	for len(level) > 1 {
		levelHash := make([][]byte, 0, len(level)/2)
		for i := 0; i < len(level); i += 2 {
			levelHash = append(levelHash, tr.Hasher.Hash(level[i], level[i+1]))
		}
		level = levelHash
	}
	// set and return the root hash
	h := level[0]
	e, err := tr.Encdec.Encode(batch)
	if err != nil {
		log.Fatal(err)
	}
	tr.KeyStore.Set(h, e)
	return h
}

// Retrieve returns a batch of data given a Merkle root
func (tr *Archive) Retrieve(root []byte) [][]byte {
	b := tr.KeyStore.Get(root)
	if b == nil {
		return nil
	}
	var data [][]byte
	err := tr.Encdec.Decode(&data, b)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// Get returns the value at the Merkle root and index
func (tr *Archive) Get(root []byte, index int) []byte {
	b := tr.KeyStore.Get(root)
	if b == nil {
		return nil
	}
	var data [][]byte
	err := tr.Encdec.Decode(&data, b)
	if err != nil {
		log.Fatal(err)
	}
	if index > len(data)-1 {
		return nil
	}
	return data[index]
}

// Prove returns the value and proof at the Merkle root and index
func (tr *Archive) Prove(root []byte, index int) ([]byte, [][]byte) {
	b := tr.KeyStore.Get(root)
	if b == nil {
		return nil, nil
	}
	var data [][]byte
	err := tr.Encdec.Decode(&data, b)
	if err != nil {
		log.Fatal(err)
	}
	if index > len(data)-1 {
		return nil, nil
	}
	datum := data[index]

	var proof [][]byte
	proof = append(proof, tr.Hasher.Hash(datum))
	// hash the first layer
	level := make([][]byte, 0, len(data))
	for _, d := range data {
		dh := tr.Hasher.Hash(d)
		level = append(level, dh)
	}
	// build up the trie until we have a root
	for len(level) > 1 {
		if index&1 == 0 {
			// if index is even, grab right twin
			proof = append(proof, level[index+1])
		} else if index&1 == 1 {
			// if index is odd, grab left twin
			proof = append(proof, level[index-1])
		}
		levelHash := make([][]byte, 0, len(level)/2)
		for i := 0; i < len(level); i += 2 {
			levelHash = append(levelHash, tr.Hasher.Hash(level[i], level[i+1]))
		}
		level = levelHash
		index = index / 2
	}
	// add the root to proof
	proof = append(proof, level[0])
	return datum, proof
}

// Validate confirms whether the proof is valid for
// the given Merkle root, index, and data
func (tr *ArchiveValidator) Validate(root []byte, index int, data []byte, proof [][]byte) bool {
	// No such thing as zero length proofs
	if proof == nil || len(proof) == 0 {
		return false
	}
	// The last proof is the rooth hash of the data, so check it
	if root == nil || !bytes.Equal(proof[len(proof)-1], root) {
		return false
	}
	h := tr.Hasher.Hash(data)
	// The first proof is the hash of the data, so check it
	if !bytes.Equal(proof[0], h) {
		return false
	}
	// Starting with the second member of the proof, up to
	// but not including the root hash, hash h with its twin
	for _, p := range proof[1 : len(proof)-1] {
		if index&1 == 0 {
			// for even indexes, twin is right twin
			h = tr.Hasher.Hash(h, p)
		} else {
			// for odd indexes, twin is left twin
			h = tr.Hasher.Hash(p, h)
		}
		index = index / 2
	}
	// check h against the root hash
	if !bytes.Equal(h, proof[len(proof)-1]) {
		return false
	}
	return true
}
