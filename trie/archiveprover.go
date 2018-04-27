package trie

import (
	"bytes"
)

// An ArchiveProver is able to validate archive trie proofs
type ArchiveProver struct {
	cypher Cypher
}

// NewArchiveProver returns a archive prover ready for use
func NewArchiveProver(cypher Cypher) *ArchiveProver {
	return &ArchiveProver{cypher: cypher}
}

// VerifyProof verifies whether a given proof is valid for a given root, index, and data
func (t *ArchiveProver) VerifyProof(root []byte, index int, data []byte, proof [][]byte) bool {
	// No such thing as zero length proofs
	if proof == nil || len(proof) == 0 {
		return false
	}
	// The last proof is the rooth hash of the data, so check it
	if root == nil || !bytes.Equal(proof[len(proof)-1], root) {
		return false
	}
	h := t.cypher.Hash(data)
	// The first proof is the hash of the data, so check it
	if !bytes.Equal(proof[0], h) {
		return false
	}
	// Starting with the second member of the proof, up to
	// but not including the root hash, hash h with its twin
	for _, p := range proof[1 : len(proof)-1] {
		if index&1 == 0 {
			// for even indexes, twin is right twin
			h = t.pairHash(h, p)
		} else {
			// for odd indexes, twin is left twin
			h = t.pairHash(p, h)
		}
		index = index / 2
	}
	// check h against the root hash
	if !bytes.Equal(h, proof[len(proof)-1]) {
		return false
	}
	return true
}

func (t *ArchiveProver) pairHash(one, two []byte) []byte {
	h := t.cypher.NewHash()
	h.Write(one)
	h.Write(two)
	return h.Sum(nil)
}
