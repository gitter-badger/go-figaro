// Package trie implements merkle tries over the database
package trie

import (
	"errors"

	"github.com/figaro-tech/go-fig-buf"
	"github.com/figaro-tech/go-fig-crypto/trie"
	"github.com/figaro-tech/go-fig-db/types"
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

	root := trie.Trie(data)
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
	proof, err := trie.Proof(data, index)
	if err != nil {
		return nil, nil, err
	}
	return data[index], proof, nil
}

// ValidateBMT validates the proof of data at index in root.
func ValidateBMT(root []byte, index int, data []byte, proof [][]byte) bool {
	return trie.Validate(root, index, data, proof)
}
