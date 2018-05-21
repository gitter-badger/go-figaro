// Package trie implements merkle tries over the database
package trie

import (
	"bytes"
	"errors"
	"math"

	"github.com/figaro-tech/go-figaro/figbuf"
	"github.com/figaro-tech/go-figaro/figcrypto/hash"
	"github.com/figaro-tech/go-figaro/figdb/types"
)

// ErrIndexOutOfRange is a self-explanatory error
var ErrIndexOutOfRange = errors.New("figdb archive: index is out of range for archive")

// Archive impelements a binary merkle trie over the database
// for archival of rarely changing data
type Archive struct {
	KeyStore types.KeyStore
	Cache    types.Cache
}

// Save creates a new entry for the batch of data
// and returns a Merkle root
func (tr *Archive) Save(data [][]byte) ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	root := trie(data)
	value := enc.EncodeBytesSlice(data)
	if tr.Cache != nil {
		tr.Cache.Add(root, value)
	}
	err := tr.KeyStore.Set(root, value)
	if err != nil {
		return nil, err
	}
	return root, nil
}

// Retrieve returns a batch of data given a Merkle root
func (tr *Archive) Retrieve(root []byte) ([][]byte, error) {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	var value []byte
	var err error
	if tr.Cache != nil {
		if c, ok := tr.Cache.Get(root); ok {
			value = c
		}
	}
	if value == nil {
		value, err = tr.KeyStore.Get(root)
	}
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	data, err := dec.DecodeBytesSlice(value)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Get returns the value at the Merkle root and index
func (tr *Archive) Get(root []byte, index int) ([]byte, error) {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	var value []byte
	var err error
	if tr.Cache != nil {
		if c, ok := tr.Cache.Get(root); ok {
			value = c
		}
	}
	if value == nil {
		value, err = tr.KeyStore.Get(root)
	}
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	data, err := dec.DecodeBytesSlice(value)
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

	var value []byte
	var err error
	if tr.Cache != nil {
		if c, ok := tr.Cache.Get(root); ok {
			value = c
		}
	}
	if value == nil {
		value, err = tr.KeyStore.Get(root)
	}
	if err != nil {
		return nil, nil, err
	}
	if value == nil {
		return nil, nil, nil
	}
	data, err := dec.DecodeBytesSlice(value)
	if err != nil {
		return nil, nil, err
	}
	if index > len(data)-1 {
		return nil, nil, ErrIndexOutOfRange
	}
	proof, err := proof(data, index)
	if err != nil {
		return nil, nil, err
	}
	return data[index], proof, nil
}

// ValidateBMT validates a binary merkle trie proof that data exists in root at index
func ValidateBMT(root []byte, index int, data []byte, proof [][]byte) bool {
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

func trie(data [][]byte) []byte {
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

func proof(data [][]byte, index int) ([][]byte, error) {
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
