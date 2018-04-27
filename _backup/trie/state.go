package trie

import (
	"bytes"
)

// StateTrie implements a modified Merkle Patricia trie as
// described here: https://github.com/ethereum/wiki/wiki/Patricia-Tree
//
// It is intended for data that will be updated often, but which must
// maintain a historic record of previous states (accessed via root hash),
// though it also supports explicit pruning of old states
type StateTrie struct {
	store  Store
	cypher Cypher
	encdec EncoderDecoder
}

// NewStateTrie creates a new StateTrie ready for use
func NewStateTrie(store Store, cypher Cypher, encdec EncoderDecoder) *StateTrie {
	return &StateTrie{store: store, cypher: cypher, encdec: encdec}
}

// hashing for normal ops, where data needs to be hashed yet
func (t *StateTrie) nodeHash(node [][]byte) []byte {
	h := t.cypher.NewHash()
	for i, b := range node {
		if i == len(node)-1 {
			h.Write(t.cypher.Hash(b))
		} else {
			h.Write(b)
		}
	}
	return h.Sum(nil)
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

func isNilNode(node [][]byte) bool {
	for _, h := range node {
		if h != nil && len(h) > 0 {
			return false
		}
	}
	return true
}

// Get fetches a value at a path, along with a proof
func (t *StateTrie) Get(root, path []byte) (data []byte) {
	p := nibbleSlice(path)
	return t.get(root, p)
}

// hash is key for node, path is remaining path to traverse before grabbing value
func (t *StateTrie) get(nodehash []byte, path []int8) []byte {
	// node is a [][]byte
	// with outer slice index 0..16 being hashes of other nodes
	// and index 17 being a (possibly nil) value in binary format
	b := t.store.Get(nodehash)
	if b == nil {
		// there is no data at this path, so return nil
		return nil
	}
	var node [][]byte
	err := t.encdec.Decode(&node, b)
	if err != nil {
		panic(err)
	}
	// No more path to traverse
	if path == nil || len(path) == 0 {
		// return the value, which is in the last index
		return node[len(node)-1]
	}
	// recurse using the hash in node at the next path step, and the remainder of the path
	return t.get(node[path[0]], path[1:])
}

// Update does stuff and returns the new root hash
func (t *StateTrie) Update(root, path, value []byte) (newroot []byte) {
	p := nibbleSlice(path)
	if newroot = t.update(root, p, value); newroot != nil {
		return newroot
	}
	return t.nodeHash(make([][]byte, 17, 17))
}

func (t *StateTrie) update(nodehash []byte, path []int8, value []byte) []byte {
	// node is a [][]byte
	// with outer slice index 0..16 being hashes of other nodes
	// and index 17 being a (possibly nil) value in binary format
	var node [][]byte
	b := t.store.Get(nodehash)
	// if node doesn't exist, initialize an empty one
	if b == nil {
		node = make([][]byte, 17, 17)
	} else {
		err := t.encdec.Decode(&node, b)
		if err != nil {
			panic(err)
		}
	}
	// No more path to traverse
	if path == nil || len(path) == 0 {
		// set the value, which is in the last index
		node[len(node)-1] = value
	} else {
		// recurse to the next node in the path, setting the new hash for this node
		node[path[0]] = t.update(node[path[0]], path[1:], value)
	}
	// hash this node, set the value in the store, and then return the hash
	if isNilNode(node) {
		return nil
	}
	h := t.nodeHash(node)
	e, err := t.encdec.Encode(node)
	if err != nil {
		panic(err)
	}
	t.store.Set(h, e)
	return h
}

// Delete removes the value at path
func (t *StateTrie) Delete(root, path []byte) (newroot []byte) {
	return t.Update(root, path, nil)
}

// Prove provides a proof for a given root, path, or data
func (t *StateTrie) Prove(root, path, data []byte) ([][][]byte, error) {
	p := nibbleSlice(path)
	return t.prove(root, p, data)
}

func (t *StateTrie) prove(nodehash []byte, path []int8, data []byte) ([][][]byte, error) {
	proof := make([][][]byte, 0, 64)
	// node is a [][]byte
	// with outer slice index 0..16 being hashes of other nodes
	// and index 17 being a (possibly nil) value in binary format
	b := t.store.Get(nodehash)
	if b == nil {
		// there is no data at this path, so return nil
		return nil, errProve
	}
	var node [][]byte
	err := t.encdec.Decode(&node, b)
	if err != nil {
		panic(err)
	}
	// No more path to traverse
	if path == nil || len(path) == 0 {
		if !bytes.Equal(node[len(node)-1], data) {
			return nil, errProve
		}
		return append(proof, t.proofNode(node)), nil
	}
	proof = append(proof, t.proofNode(node))
	// recurse using the hash in node at the next path step, and the remainder of the path
	p, err := t.prove(node[path[0]], path[1:], data)
	if err != nil {
		return nil, err
	}
	return append(proof, p...), nil
}

func (t *StateTrie) proofNode(node [][]byte) [][]byte {
	node[len(node)-1] = t.cypher.Hash(node[len(node)-1])
	return node
}
