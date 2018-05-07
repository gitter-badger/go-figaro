// Package badger implements BadgerDb key/value store
package badger

import (
	"log"

	"github.com/dgraph-io/badger"
	"github.com/figaro-tech/go-figaro/figdb/internal"
)

// Store sets up an BadgerDB key/value store
type Store struct {
	DB      *badger.DB
	batch   bool
	pending internal.KeyStoreUpdateBatch
}

// NewStore returns a badger.KeyStore ready for use
func NewStore(dir string) *Store {
	s := &Store{}
	s.Open(dir)
	return s
}

// Get gets a value from the db at key
func (s *Store) Get(key []byte) []byte {
	var r []byte
	err := s.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		r = make([]byte, len(val))
		copy(r, val)
		return nil
	})
	if err == badger.ErrKeyNotFound {
		return r
	}
	if err != nil {
		log.Panic(err)
	}
	return r
}

// Set sets a value in the db at key
func (s *Store) Set(key []byte, value []byte) {
	defer s.GC()
	if s.batch {
		s.pending = append(s.pending, internal.KeyStoreUpdate{Key: key, Value: value})
		return
	}
	err := s.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
	if err != nil {
		log.Panic(err)
	}
}

// Delete deletes a value in the db at key
func (s *Store) Delete(key []byte) {
	defer s.GC()
	if s.batch {
		s.pending = append(s.pending, internal.KeyStoreUpdate{Key: key, Value: nil})
		return
	}
	err := s.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})
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
	defer s.GC()
	err := s.DB.Update(func(txn *badger.Txn) error {
		for _, update := range updates {
			if update.Value == nil {
				err := txn.Delete(update.Key)
				if err != nil {
					return err
				}
			} else {
				err := txn.Set(update.Key, update.Value)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// Open will create a db from dir
func (s *Store) Open(dir string) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	s.DB = db
}

// Close closes the lock on the underlying db
func (s *Store) Close() {
	err := s.DB.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// GC removes older versions and GCs the values,
// run it early and often
func (s *Store) GC() {
	var err error
	err = s.DB.PurgeOlderVersions()
	if err != nil {
		log.Panic(err)
	}
	// Recommended space reclemation ratio
	err = s.DB.RunValueLogGC(0.5)
	if err == badger.ErrRejected {
		log.Print(err)
		return
	}
	if err == badger.ErrNoRewrite {
		log.Print(err)
		return
	}
	if err != nil {
		log.Panic(err)
	}
}
