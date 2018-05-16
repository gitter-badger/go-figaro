package mock

import (
	"fmt"

	"github.com/figaro-tech/go-figaro/figdb/types"
)

// KeyStore sets up an in-memory key/value store
type KeyStore struct {
	DB      map[string]string
	batch   bool
	pending types.KeyStoreUpdateBatch
}

// NewKeyStore makes a new KeyStore
func NewKeyStore() *KeyStore {
	ks := &KeyStore{
		DB: make(map[string]string),
	}
	return ks
}

// Get returns a trie value given a trie key
func (ks *KeyStore) Get(key types.Key) ([]byte, error) {
	v := ks.DB[key.String()]
	if v == "" {
		return nil, nil
	}
	return []byte(v), nil
}

// Set updates a trie key with a trie value
func (ks *KeyStore) Set(key types.Key, value []byte) error {
	if value != nil {
		ks.DB[key.String()] = string(value)
	}
	return nil
}

// Delete removes a trie key/value
func (ks *KeyStore) Delete(key types.Key) error {
	delete(ks.DB, key.String())
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

func (ks *KeyStore) String() string {
	s := ""
	for k, v := range ks.DB {
		s += fmt.Sprintf("%#x : %#x\n", k, v)
	}
	return s
}
