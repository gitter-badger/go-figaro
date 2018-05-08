// Package trie implements the figaro.*Trie interfaces
package trie

import (
	"bytes"
	"errors"
	"math"

	"github.com/figaro-tech/go-figaro/figcrypto"

	"github.com/figaro-tech/go-figaro/figbuf"
	"github.com/figaro-tech/go-figaro/figdb/types"
)

// ErrIndexOutOfRange is a self-explanatory error
var ErrIndexOutOfRange = errors.New("figdb archive: index is out of range for archive")

// Archive impelements a pkg.ArchiveTrie
type Archive struct {
	KeyStore types.KeyStore
}

func (tr *Archive) hash(d ...[]byte) []byte {
	h := figcrypto.HasherPool.Get().(*figcrypto.Hasher)
	defer figcrypto.HasherPool.Put(h)
	return h.Hash(d...)
}

// ArchiveValidator implements a pkg.ArchiveValidator
type ArchiveValidator struct {
}

func (tr *ArchiveValidator) hash(d ...[]byte) []byte {
	h := figcrypto.HasherPool.Get().(*figcrypto.Hasher)
	defer figcrypto.HasherPool.Put(h)
	return h.Hash(d...)
}

// Save creates a new entry for the batch of data
// and returns a Merkle root
func (tr *Archive) Save(data [][]byte) ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	e := enc.EncodeBytesSlice(data)
	if len(data)&1 == 1 {
		data = append(data, nil)
	}
	trie := make([][]byte, len(data))
	for i, d := range data {
		trie[i] = tr.hash(d)
	}
	for {
		for i, j := 0, 0; i < len(trie); i, j = i+2, j+1 {
			trie[j] = tr.hash(trie[i], trie[i+1])
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
	h := trie[0]
	err := tr.KeyStore.Set(h, e)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// Retrieve returns a batch of data given a Merkle root
func (tr *Archive) Retrieve(root []byte) ([][]byte, error) {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	b, err := tr.KeyStore.Get(root)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}
	data, _, err := dec.DecodeBytesSlice(b)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Get returns the value at the Merkle root and index
func (tr *Archive) Get(root []byte, index int) ([]byte, error) {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	b, err := tr.KeyStore.Get(root)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}
	data, _, err := dec.DecodeBytesSlice(b)
	if err != nil {
		return nil, err
	}
	if index > len(data)-1 {
		return nil, ErrIndexOutOfRange
	}
	return data[index], nil
}

// GetAndProve returns the value and proof at the Merkle root and index
func (tr *Archive) GetAndProve(root []byte, index int) ([]byte, [][]byte, error) {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	b, err := tr.KeyStore.Get(root)
	if err != nil {
		return nil, nil, err
	}
	if b == nil {
		return nil, nil, nil
	}
	data, _, err := dec.DecodeBytesSlice(b)
	if err != nil {
		return nil, nil, err
	}
	if index > len(data)-1 {
		return nil, nil, ErrIndexOutOfRange
	}
	if len(data)&1 == 1 {
		data = append(data, nil)
	}
	datum := data[index]
	trie := make([][]byte, len(data))
	for i, d := range data {
		trie[i] = tr.hash(d)
	}
	proof := make([][]byte, merkleLen(len(data)))
	proof[0] = trie[index]
	for k := 1; ; k++ {
		proof[k] = trie[index+1-(index&1*2)]
		for i, j := 0, 0; i < len(trie); i, j = i+2, j+1 {
			trie[j] = tr.hash(trie[i], trie[i+1])
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
	return datum, proof, nil
}

// Validate confirms whether the proof is valid for
// the given Merkle root, index, and data
func (tr *ArchiveValidator) Validate(root []byte, index int, data []byte, proof [][]byte) bool {
	if proof == nil || len(proof) == 0 {
		return false
	}
	if root == nil || !bytes.Equal(proof[len(proof)-1], root) {
		return false
	}
	h := tr.hash(data)
	if !bytes.Equal(proof[0], h) {
		return false
	}
	for _, p := range proof[1 : len(proof)-1] {
		if index&1 == 0 {
			h = tr.hash(h, p)
		} else {
			h = tr.hash(p, h)
		}
		index = index / 2
	}
	if !bytes.Equal(h, proof[len(proof)-1]) {
		return false
	}
	return true
}

func merkleLen(l int) int {
	return int(math.Ceil(math.Log2(float64(l)))) + 2
}
