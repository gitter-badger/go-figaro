package trie

import (
	"errors"
	"hash"
)

var (
	errProve = errors.New("trie prove: data does not exist at path")
)

// Store provides a key/value store to a StateTrie
type Store interface {
	Get(key []byte) []byte
	Set(key []byte, value []byte)
}

// Cypher provides hashing capabilities for state tries
type Cypher interface {
	NewHash() hash.Hash
	Hash([]byte) []byte
}

// EncoderDecoder provides binary encoding/decoding for state tries
type EncoderDecoder interface {
	Encode(interface{}) ([]byte, error)
	Decode(interface{}, []byte) error
}
