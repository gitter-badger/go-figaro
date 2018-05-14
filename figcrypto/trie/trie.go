// Package trie provides cryptographic functions
package trie

import (
	"bytes"
	"errors"
	"math"

	"github.com/figaro-tech/go-figaro/figcrypto/hash"
)

// ErrIndexOutOfRange is a self-explanatory error
var ErrIndexOutOfRange = errors.New("figcrypto trie: index is out of range for data")

// Trie constructs a root hash from an ordered list of data
// based on the binary merkle algorithm
func Trie(data [][]byte) []byte {
	if len(data)&1 == 1 {
		data = append(data, nil)
	}
	trie := make([][]byte, len(data))
	for i, d := range data {
		trie[i] = hash.Hash256(d)
	}
	for {
		for i, j := 0, 0; i < len(trie); i, j = i+2, j+1 {
			trie[j] = hash.Hash256(trie[i], trie[i+1])
		}
		l := len(trie) / 2
		if l == 1 {
			break
		}
		if l&1 == 1 {
			trie = trie[:len(trie)/2+1]
		} else {
			trie = trie[:len(trie)/2]
		}
	}
	return trie[0]
}

// Proof construct a merkle proof of the datum in data at index
func Proof(data [][]byte, index int) ([][]byte, error) {
	if index > len(data)-1 {
		return nil, ErrIndexOutOfRange
	}
	if len(data)&1 == 1 {
		data = append(data, nil)
	}
	trie := make([][]byte, len(data))
	for i, d := range data {
		trie[i] = hash.Hash256(d)
	}

	proof := make([][]byte, int(math.Ceil(math.Log2(float64(len(data)))))+1)
	for k := 0; ; k++ {
		proof[k] = trie[index+1-(index&1*2)]
		for i, j := 0, 0; i < len(trie); i, j = i+2, j+1 {
			trie[j] = hash.Hash256(trie[i], trie[i+1])
		}
		l := len(trie) / 2
		if l == 1 {
			break
		}
		index = index / 2
		if l&1 == 1 {
			trie = trie[:len(trie)/2+1]
		} else {
			trie = trie[:len(trie)/2]
		}
	}
	proof[len(proof)-1] = trie[0]
	return proof, nil
}

// Validate validates the proof that data exists in root at index
func Validate(root []byte, index int, data []byte, proof [][]byte) bool {
	dh := hash.Hash256(data)
	for _, p := range proof[:len(proof)-1] {
		if index&1 == 0 {
			dh = hash.Hash256(dh, p)
		} else {
			dh = hash.Hash256(p, dh)
		}
		index = index / 2
	}
	if !bytes.Equal(dh, proof[len(proof)-1]) || !bytes.Equal(proof[len(proof)-1], root) {
		return false
	}
	return true
}
