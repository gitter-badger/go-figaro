// Package internal is the main package for figdb
package internal

// KeyStore provides a key/value store
type KeyStore interface {
	Get(key []byte) []byte
	Set(key []byte, value []byte)
	Delete(key []byte)
	// Batch starts automatically batching Set() and Delete(),
	// holding onto the KeyStoreUpdateBatch until Save() is called
	Batch()
	// Write saves pending updates in the batch
	Write()
	// BatchUpdate updates a batch of updates that have been
	// manually queued
	BatchUpdate(updates KeyStoreUpdateBatch)
	Open(dir string)
	Close()
}

// KeyStoreUpdate update represents a single update in a batch of updates
type KeyStoreUpdate struct {
	Key   []byte
	Value []byte
}

// KeyStoreUpdateBatch is an ordered list of KeyStoreUpdate
type KeyStoreUpdateBatch []KeyStoreUpdate

// A Hasher provides cryptographic hashing
type Hasher interface {
	// Hash accepts 0 or more []byte and hashes them, returning a []byte
	Hash(b ...[]byte) []byte
}

// EncoderDecoder provides determinstic binary encoding/decoding
type EncoderDecoder interface {
	Encoder
	Decoder
}

// Encoder provides deterministic binary encoding
type Encoder interface {
	// Encode binary encodes a given source
	Encode(src interface{}) ([]byte, error)
}

// Decoder provides deterministic binary decoding
type Decoder interface {
	// Decode binary decodes data into a given destination
	Decode(dest interface{}, b []byte) error
}
