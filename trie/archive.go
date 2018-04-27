package trie

import (
	"bytes"
)

// ArchiveTrie implements a binary Merkle trie
//
// It is intended for archive data that is created
// as a batch and then never updated. The entire trie
// is saved as a single entry in a key/value store,
// but efficient proofs can be provided that a single
// value resides at a given index
type ArchiveTrie struct {
	store  Store
	cypher Cypher
	encdec EncoderDecoder
}

// NewArchiveTrie creates a new ArchiveTrie ready for use
func NewArchiveTrie(store Store, cypher Cypher, encdec EncoderDecoder) *ArchiveTrie {
	return &ArchiveTrie{store: store, cypher: cypher, encdec: encdec}
}

// Archive stores data in a key/value store, returning a root hash
func (t *ArchiveTrie) Archive(data [][]byte) ([]byte, error) {
	// construct a root hash, end then store the encoded data at the root
	// if data has an odd length, append nil to make it even
	if len(data)&1 == 1 {
		data = append(data, nil)
	}
	// hash the first layer
	level := make([][]byte, 0, len(data))
	for _, d := range data {
		level = append(level, t.cypher.Hash(d))
	}
	// build up the trie until we have a root
	for len(level) > 1 {
		levelHash := make([][]byte, 0, len(level)/2)
		for i := 0; i < len(level); i += 2 {
			levelHash = append(levelHash, t.pairHash(level[i], level[i+1]))
		}
		level = levelHash
	}
	// set and return the root hash
	h := level[0]
	e, err := t.encdec.Encode(data)
	if err != nil {
		return nil, err
	}
	t.store.Set(h, e)
	return h, nil
}

// Retrieve gets data at a given index, given a root
func (t *ArchiveTrie) Retrieve(root []byte, index int) []byte {
	b := t.store.Get(root)
	if b == nil {
		return nil
	}
	var data [][]byte
	err := t.encdec.Decode(&data, b)
	if err != nil {
		panic(err)
	}
	if index > len(data)-1 {
		return nil
	}
	return data[index]
}

// Prove returns a merkle proof of data at a given index into a root
func (t *ArchiveTrie) Prove(root []byte, index int, datum []byte) ([][]byte, error) {
	b := t.store.Get(root)
	if b == nil {
		return nil, errProve
	}
	var data [][]byte
	err := t.encdec.Decode(&data, b)
	if err != nil {
		return nil, err
	}
	if index > len(data)-1 {
		return nil, errProve
	}
	if !bytes.Equal(data[index], datum) {
		return nil, errProve
	}
	var proof [][]byte
	proof = append(proof, t.cypher.Hash(datum))
	// hash the first layer
	level := make([][]byte, 0, len(data))
	for _, d := range data {
		dh := t.cypher.Hash(d)
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
			levelHash = append(levelHash, t.pairHash(level[i], level[i+1]))
		}
		level = levelHash
		index = index / 2
	}
	// add the root to proof
	proof = append(proof, level[0])
	return proof, nil
}

func (t *ArchiveTrie) pairHash(one, two []byte) []byte {
	h := t.cypher.NewHash()
	h.Write(one)
	h.Write(two)
	return h.Sum(nil)
}
