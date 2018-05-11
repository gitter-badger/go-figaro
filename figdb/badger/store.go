// Package badger implements BadgerDb key/value store
package badger

import (
	"errors"
	"io"
	"log"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/figaro-tech/go-figaro/figdb/types"
)

// ErrCalledAfterClose is a self-explanatory error
var ErrCalledAfterClose = errors.New("figdb: operation called after DB was closed")

// KeyStore is a key/value store backed by BadgerDB
type KeyStore struct {
	block   sync.RWMutex
	gclock  sync.RWMutex
	DB      *badger.DB
	batchDB map[string]string
	batch   bool
}

// NewKeyStore returns a badger.KeyStore ready for use
func NewKeyStore(dir string) *KeyStore {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return &KeyStore{DB: db}
}

// Close closes the DB and releases the file lock
func (ks *KeyStore) Close() {
	ks.gclock.RLock()
	defer ks.gclock.RUnlock()

	err := ks.DB.Close()
	if err != nil {
		log.Fatal(err)
	}
	ks.DB = nil
}

// NOTE: read locks are taken throughout to prevent concurrent updates while loading
// in a database backup (which takes a full lock). Locks serve no other purpose, as
// Badger is safe for concurrent access via transactions

// Get gets a value from the db at key
func (ks *KeyStore) Get(key []byte) ([]byte, error) {
	if ks.DB == nil {
		log.Panic(ErrCalledAfterClose)
	}

	ks.gclock.RLock()
	defer ks.gclock.RUnlock()

	if ks.batch {
		v, err := ks.getFromBatch(key)
		if err != badger.ErrKeyNotFound {
			return v, err
		}
	}

	var r []byte
	err := ks.DB.View(func(txn *badger.Txn) error {
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
		return r, nil
	}
	if err != nil {
		return r, err
	}
	return r, nil
}

// Set sets a value in the db at key
func (ks *KeyStore) Set(key, value []byte) error {
	if ks.DB == nil {
		log.Panic(ErrCalledAfterClose)
	}

	ks.gclock.RLock()
	defer ks.gclock.RUnlock()

	if ks.batch {
		return ks.addToBatch(key, value)
	}

	err := ks.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
	return err
}

// Delete deletes a value in the db at key
func (ks *KeyStore) Delete(key []byte) error {
	if ks.DB == nil {
		log.Panic(ErrCalledAfterClose)
	}

	ks.gclock.RLock()
	defer ks.gclock.RUnlock()

	if ks.batch {
		return ks.addToBatch(key, nil)
	}

	err := ks.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
	return err
}

// Batch causes all writes to be batched until Save() is called
func (ks *KeyStore) Batch() {
	if ks.DB == nil {
		log.Panic(ErrCalledAfterClose)
	}

	ks.block.Lock()
	defer ks.block.Unlock()

	ks.batchDB = make(map[string]string)
	ks.batch = true
}

// Write savees all pending updates in the batch
func (ks *KeyStore) Write() error {
	if ks.DB == nil {
		log.Panic(ErrCalledAfterClose)
	}

	ks.block.Lock()
	defer ks.block.Unlock()

	pending := make(types.KeyStoreUpdateBatch, len(ks.batchDB))
	var i int
	for k, v := range ks.batchDB {
		pending[i].Key = []byte(k)
		pending[i].Value = []byte(v)
		i++
	}

	err := ks.BatchUpdate(pending)

	ks.batch = false
	return err
}

// BatchUpdate executes multiple set or delets at once
func (ks *KeyStore) BatchUpdate(updates types.KeyStoreUpdateBatch) error {
	if ks.DB == nil {
		log.Panic(ErrCalledAfterClose)
	}

	ks.gclock.RLock()
	defer ks.gclock.RUnlock()

	err := ks.DB.Update(func(txn *badger.Txn) error {
		var err error
		for _, update := range updates {
			if update.Value == nil {
				err = txn.Delete(update.Key)
			} else {
				err = txn.Set(update.Key, update.Value)
			}
			if err != nil {
				return err
			}
		}
		return err
	})
	return err
}

// GC removes older versions and GCs the values,
// run it early and often
func (ks *KeyStore) GC() {
	ks.gclock.Lock()
	defer ks.gclock.Unlock()

	var err error
	err = ks.DB.PurgeOlderVersions()
	if err != nil {
		log.Panic(err)
	}
	for {
		err = ks.DB.RunValueLogGC(0.5)
		if err == badger.ErrNoRewrite {
			break
		}
		if err == badger.ErrRejected {
			log.Print(err)
			break
		}
		if err != nil {
			log.Panic(err)
		}
	}
}

// Backup will backup a snapshot of all entries newer than the provided timestamp,
// returning a new timestamp that can be passed to future invocations of Backup
func (ks *KeyStore) Backup(w io.Writer, since uint64) (uint64, error) {
	ks.gclock.RLock()
	defer ks.gclock.RUnlock()

	return ks.DB.Backup(w, since)
}

// Load will restore a backup after acquiring a lock
func (ks *KeyStore) Load(r io.Reader) error {
	ks.gclock.RLock()
	defer ks.gclock.RUnlock()

	return ks.DB.Load(r)
}

func (ks *KeyStore) addToBatch(key, value []byte) error {
	ks.block.Lock()
	defer ks.block.Unlock()

	ks.batchDB[string(key)] = string(value)
	return nil
}

func (ks *KeyStore) getFromBatch(key []byte) ([]byte, error) {
	ks.block.RLock()
	defer ks.block.RUnlock()

	v := ks.batchDB[string(key)]
	if v == "" {
		return nil, badger.ErrKeyNotFound
	}
	return []byte(v), nil
}
