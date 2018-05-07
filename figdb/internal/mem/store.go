// Package mem implements in-memory key/value store
package mem

import (
	"fmt"

	"github.com/figaro-tech/go-figaro/figdb/internal"
)

// Store sets up an in-memory key/value store
type Store struct {
	DB      map[string][]byte
	batch   bool
	pending internal.KeyStoreUpdateBatch
}

// NewStore returns an initialized *Store
func NewStore() *Store {
	store := &Store{}
	store.Open("")
	return nil
}

// Get returns a trie value given a trie key
func (s *Store) Get(key []byte) []byte {
	v := s.DB[string(key)]
	if v == nil {
		return nil
	}
	c := make([]byte, len(v))
	copy(c, v)
	return c
}

// Set updates a trie key with a trie value
func (s *Store) Set(key []byte, value []byte) {
	if value != nil {
		c := make([]byte, len(value))
		copy(c, value)
		s.DB[string(key)] = c
	}
}

// Delete removes a trie key/value
func (s *Store) Delete(key []byte) {
	delete(s.DB, string(key))
}

// Batch causes all writes to be batched until Save() is called
func (s *Store) Batch() {
	s.batch = true
}

// Write avees all pending updates in the batch
func (s *Store) Write() {
	s.batch = false
	s.BatchUpdate(s.pending)
	s.pending = s.pending[:0]
}

// BatchUpdate executs multiple set or delets at once
func (s *Store) BatchUpdate(updates internal.KeyStoreUpdateBatch) {
	for _, update := range updates {
		if update.Value == nil {
			s.Delete(update.Key)
		} else {
			s.Set(update.Key, update.Value)
		}
	}
}

// Open sets up the db
func (s *Store) Open(dir string) {
	s.DB = make(map[string][]byte)
}

// Close resets the store for a memdb
func (s *Store) Close() {
	s.DB = nil
}

func (s *Store) String() string {
	var str string
	for k, v := range s.DB {
		str = str + fmt.Sprintf("% x:% x\n", k, v)
	}
	return str
}

// Len returns the number of entries in the store
func (s *Store) Len() int {
	return len(s.DB)
}
