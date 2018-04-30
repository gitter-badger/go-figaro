// Package level implements a LevelDB store
package level

import (
	"log"

	"github.com/figaro-tech/figaro/pkg/figdb/internal"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Store implements a KeyStore
type Store struct {
	DB      *leveldb.DB
	batch   bool
	pending internal.KeyStoreUpdateBatch
}

// NewStore returns a level.KeyStore ready for use
func NewStore(dir string) *Store {
	s := &Store{}
	s.Open(dir)
	return s
}

// Get gets a value from the db at key
func (s *Store) Get(key []byte) []byte {
	data, err := s.DB.Get(key, nil)
	if err != nil {
		log.Panic(err)
	}
	r := make([]byte, len(data))
	copy(r, data)
	return r
}

// Set sets a value in the db at key
func (s *Store) Set(key []byte, value []byte) {
	err := s.DB.Put(key, value, nil)
	if err != nil {
		log.Panic(err)
	}
}

// Delete deletes a value in the db at key
func (s *Store) Delete(key []byte) {
	err := s.DB.Delete(key, nil)
	if err != nil {
		log.Panic(err)
	}
}

// Batch causes all writes to be batched until Save() is called
func (s *Store) Batch() {
	s.batch = true
}

// Write savees all pending updates in the batch
func (s *Store) Write() {
	s.batch = false
	s.BatchUpdate(s.pending)
	s.pending = s.pending[:0]
}

// BatchUpdate executes multiple set or delets at once
func (s *Store) BatchUpdate(updates internal.KeyStoreUpdateBatch) {
	batch := new(leveldb.Batch)
	for _, update := range updates {
		if update.Value == nil {
			batch.Delete(update.Key)
		} else {
			batch.Put(update.Key, update.Value)
		}
	}
	err := s.DB.Write(batch, nil)
	if err != nil {
		log.Panic(err)
	}
}

// Open will create a db from dir
func (s *Store) Open(dir string) {
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}
	db, err := leveldb.OpenFile(dir, o)
	if err != nil {
		log.Fatal(err)
	}
	s.DB = db
}

// Close closes the lock on the underlying db
func (s *Store) Close() {
	s.DB.Close()
}
