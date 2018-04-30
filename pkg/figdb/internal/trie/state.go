// Package trie implements the figaro.*Trie interfaces
package trie

import (
	"bytes"
	"log"

	"github.com/figaro-tech/figaro/pkg/figdb/internal"
)

// State implements a pkg.StateTrie
type State struct {
	KeyStore internal.KeyStore
	Hasher   internal.Hasher
	Encdec   internal.EncoderDecoder
}

// StateValidator implements a pkg.StateValidator
type StateValidator struct {
	Hasher internal.Hasher
}

// Get takes a Merkle root and key and returns a value
func (tr *State) Get(root, key []byte) []byte {
	p := nibbleSlice(key)
	return tr.get(root, p)
}

func (tr *State) get(nodehash []byte, path []int8) []byte {
	// node is a [][]byte
	// with outer slice index 0..16 being hashes of other nodes
	// and index 17 being a (possibly nil) value in binary format
	b := tr.KeyStore.Get(nodehash)
	if b == nil {
		// there is no data at this path, so return nil
		return nil
	}
	var node [][]byte
	err := tr.Encdec.Decode(&node, b)
	if err != nil {
		log.Fatal(err)
	}
	// No more path to traverse
	if path == nil || len(path) == 0 {
		// return the value, which is in the last index
		return node[len(node)-1]
	}
	// recurse using the hash in node at the next path step, and the remainder of the path
	return tr.get(node[path[0]], path[1:])
}

// Set updates the value at a Merkle root and key and returns a Merkle root
func (tr *State) Set(root, key, value []byte) []byte {
	p := nibbleSlice(key)
	if newroot := tr.set(root, p, value); newroot != nil {
		return newroot
	}
	return tr.Hasher.Hash(make([][]byte, 17, 17)...)
}

func (tr *State) set(nodehash []byte, path []int8, value []byte) []byte {
	// node is a [][]byte
	// with outer slice index 0..16 being hashes of other nodes
	// and index 17 being a (possibly nil) value in binary format
	var node [][]byte
	b := tr.KeyStore.Get(nodehash)
	// if node doesn't exist, initialize an empty one
	if b == nil {
		node = make([][]byte, 17, 17)
	} else {
		err := tr.Encdec.Decode(&node, b)
		if err != nil {
			log.Fatal(err)
		}
	}
	// No more path to traverse
	if path == nil || len(path) == 0 {
		// set the value, which is in the last index
		node[len(node)-1] = value
	} else {
		// recurse to the next node in the path, setting the new hash for this node
		node[path[0]] = tr.set(node[path[0]], path[1:], value)
	}
	// If node is all nils, then we can just return nil
	if isNil(node) {
		return nil
	}
	// hash this node, set the value in the store, and then return the hash
	e, err := tr.Encdec.Encode(node)
	if err != nil {
		log.Fatal(err)
	}
	node[len(node)-1] = tr.Hasher.Hash(node[len(node)-1])
	h := tr.Hasher.Hash(node...)
	tr.KeyStore.Set(h, e)
	return h
}

// Prove returns the value and proof of a Merkle root and key
func (tr *State) Prove(root, key []byte) ([]byte, [][][]byte) {
	p := nibbleSlice(key)
	return tr.prove(root, p)
}

func (tr *State) prove(nodehash []byte, path []int8) ([]byte, [][][]byte) {
	var proof [][][]byte
	// node is a [][]byte
	// with outer slice index 0..16 being hashes of other nodes
	// and index 17 being a (possibly nil) value in binary format
	b := tr.KeyStore.Get(nodehash)
	if b == nil {
		// there is no data at this path, so return nil
		return nil, nil
	}
	var node [][]byte
	err := tr.Encdec.Decode(&node, b)
	if err != nil {
		log.Fatal(err)
	}
	data := node[len(node)-1]
	node[len(node)-1] = tr.Hasher.Hash(node[len(node)-1])
	// No more path to traverse
	if path == nil || len(path) == 0 {
		proof = make([][][]byte, 0, 64)
		proof = append(proof, node)
		return data, proof
	}
	// recurse using the hash in node at the next path step, and the remainder of the path
	data, proof = tr.prove(node[path[0]], path[1:])
	proof = append(proof, node)
	return data, proof
}

// Validate confirms whether the proof is valid for the given Merkle root, path, and data
func (tr *StateValidator) Validate(root, key, data []byte, proof [][][]byte) bool {
	p := nibbleSlice(key)
	return tr.validate(root, p, data, proof)
}

func (tr *StateValidator) validate(root []byte, path []int8, data []byte, proof [][][]byte) bool {
	// If this doesn't match, then it's defacto invalid
	if len(path)+1 != len(proof) {
		return false
	}
	node := proof[0]
	h := tr.Hasher.Hash(data)
	if !bytes.Equal(node[len(node)-1], h) {
		return false
	}
	for i := 0; i < len(proof)-1; i++ {
		node = proof[i]
		if !bytes.Equal(node[path[i]], h) {
			return false
		}
		h = tr.Hasher.Hash(node...)
	}
	if !bytes.Equal(root, h) {
		return false
	}
	return true
}

// nibbleSlice converts a slice of bytes into a slice of nibbles
// where a nibble is an int with value 0..f
func nibbleSlice(bytes []byte) []int8 {
	ns := make([]int8, 0, len(bytes)*2)
	for _, b := range bytes {
		np := nibblePair(b)
		ns = append(ns, np[:]...)
	}
	return ns
}

// nibblePair returns a []int8 pair of nibbles from a byte
func nibblePair(b byte) []int8 {
	var np [2]int8
	np[0] = int8(b >> 4)
	np[1] = int8(b & 0xf)
	return np[:]
}

func isNil(node [][]byte) bool {
	for _, h := range node {
		if h != nil && len(h) > 0 {
			return false
		}
	}
	return true
}
