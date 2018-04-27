package trie

import (
	"bytes"
)

// A StateProver is able to validate state trie proofs
type StateProver struct {
	cypher Cypher
}

// NewStateProver returns a state prover ready for use
func NewStateProver(cypher Cypher) *StateProver {
	return &StateProver{cypher: cypher}
}

// hashing for proof, where data is already hashed
func (t *StateProver) proofHash(node [][]byte) []byte {
	h := t.cypher.NewHash()
	for _, b := range node {
		h.Write(b)
	}
	return h.Sum(nil)
}

// VerifyProof verifies whether a given proof is valid for a given root, path, and data
func (t *StateProver) VerifyProof(root, path, data []byte, proof [][][]byte) bool {
	p := nibbleSlice(path)
	return t.verifyProof(root, p, data, proof)
}

func (t *StateProver) verifyProof(root []byte, path []int8, data []byte, proof [][][]byte) bool {
	if len(path)+1 != len(proof) {
		return false
	}
	h := t.cypher.Hash(data)
	for i := len(proof) - 1; i > -1; i-- {
		node := proof[i]
		if i == len(proof)-1 {
			if !bytes.Equal(node[len(node)-1], h) {
				return false
			}
			h = t.proofHash(node)
			continue
		}
		if !bytes.Equal(node[path[i]], h) {
			return false
		}
		h = t.proofHash(node)
	}
	if !bytes.Equal(root, h) {
		return false
	}
	return true
}
