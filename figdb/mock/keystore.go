package mock

import (
	"github.com/figaro-tech/go-figaro/figdb/types"
)

// KeyStore sets up an in-memory key/value store
type KeyStore struct {
	DB      map[string][]byte
	batch   bool
	pending types.KeyStoreUpdateBatch
}

// NewKeyStore makes a new KeyStore
func NewKeyStore() *KeyStore {
	ks := &KeyStore{
		DB: make(map[string][]byte),
	}
	return ks
}

// Get returns a trie value given a trie key
func (ks *KeyStore) Get(key []byte) ([]byte, error) {
	v := ks.DB[string(key)]
	if v == nil {
		return nil, nil
	}
	c := make([]byte, len(v))
	copy(c, v)
	return c, nil
}

// Set updates a trie key with a trie value
func (ks *KeyStore) Set(key []byte, value []byte) error {
	if value != nil {
		c := make([]byte, len(value))
		copy(c, value)
		ks.DB[string(key)] = c
	}
	return nil
}

// Delete removes a trie key/value
func (ks *KeyStore) Delete(key []byte) error {
	delete(ks.DB, string(key))
	return nil
}

// Batch causes all writes to be batched until Save() is called
func (ks *KeyStore) Batch() {
	ks.batch = true
}

// Write avees all pending updates in the batch
func (ks *KeyStore) Write() error {
	ks.batch = false
	ks.BatchUpdate(ks.pending)
	ks.pending = ks.pending[:0]
	return nil
}

// BatchUpdate executs multiple set or delets at once
func (ks *KeyStore) BatchUpdate(updates types.KeyStoreUpdateBatch) error {
	for _, update := range updates {
		if update.Value == nil {
			ks.Delete(update.Key)
		} else {
			ks.Set(update.Key, update.Value)
		}
	}
	return nil
}
