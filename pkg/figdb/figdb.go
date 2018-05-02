// Package figdb implements a high-performance Merkle database
package figdb

import (
	"io"
	"log"

	"github.com/figaro-tech/figaro/pkg/figcrypto"
	"github.com/figaro-tech/figaro/pkg/figdb/internal"
	"github.com/figaro-tech/figaro/pkg/figdb/internal/badger"
	"github.com/figaro-tech/figaro/pkg/figdb/internal/level"
	"github.com/figaro-tech/figaro/pkg/figdb/internal/mem"
	"github.com/figaro-tech/figaro/pkg/figdb/internal/trie"
	"github.com/figaro-tech/figaro/pkg/figdb/mock"
)

type figdbType int

const (
	levelDB figdbType = iota
	badgerDB
)
const dbType = levelDB

// FigDB is the interface into FigDB types
type FigDB interface {
	// State gets the figdb.StateTrie
	State() StateTrie
	// Archive gets the figdb.ArchiveTrie
	Archive() ArchiveTrie
	// Calling Batch will cause all updates to be saved
	// to a pending pool until Write is called
	Batch()
	// Write saves any pending writes
	Write()
	// Close closes the underlying db
	Close()
	// Backup backs up the underlying store
	Backup(writer io.Writer, since uint64) (uint64, error)
	// Restore from a backup
	Restore(reader io.Reader) error
	// Force garbage collection
	GC()
}

// StateTrie stores each  state update as an entry in the
// backing keystore, supporting retrieval of the value of any key at
// any given world state root, as well as efficient Merkle proof of
// the value of any key at any given world state root. States use
// modified Merkle Patricia tries.
type StateTrie interface {
	// Set updates the value at key in the given Merkle root, returning a new Merkle root
	Set(root, key, value []byte) []byte
	// Get returns the value at key in given Merkle root
	Get(root, key []byte) []byte
	// Prove returns the value at key in a given Merkle root,
	// along with a Merkle proof of the result
	Prove(root, key []byte) ([]byte, [][][]byte)
}

// StateValidator validates a proof of a value at a key in a given Merkle root,
// with zero knowledge of the world state
type StateValidator interface {
	Validate(root, key, data []byte, proof [][][]byte) bool
}

// ArchiveTrie stores a collection of data as a single archive
// in a backing keystore, supporting retrieval of the entire archive
// or single item, as well as efficient Merkle proofs of inclusion of
// an item at an index. Archives use binary Merkle tries.
type ArchiveTrie interface {
	// Save saves a collection of data, returning a Merkle root of the archive
	Save(collection [][]byte) []byte
	// Retrieve returns a collection of data stored at a Merkle root
	Retrieve(root []byte) [][]byte
	// Get returns the data at index in a collection stored at a Merkle root
	Get(root []byte, index int) []byte
	// Prove returns the data at index in a collection stored at a Merkle root,
	// along with a Merkle proof of the result
	Prove(root []byte, index int) ([]byte, [][]byte)
}

// ArchiveValidator validates a proof of a single value at index of an archive
// with zero knowledge of the batch itself
type ArchiveValidator interface {
	Validate(root []byte, index int, data []byte, proof [][]byte) bool
}

// New returns a FigDB backed by a high-performance disk database
func New(dir string) FigDB {
	var store internal.KeyStore
	if dbType == levelDB {
		store = level.NewStore(dir)
	} else if dbType == badgerDB {
		store = badger.NewStore(dir)
	} else {
		log.Panicf("Unsupported dbType %d", dbType)
	}
	state, archive := newTries(store)
	return &db{store, state, archive}
}

// NewMemDB returns a FigDB backed by a high-performance memory database
func NewMemDB() FigDB {
	store := mem.NewStore()
	state, archive := newTries(store)
	return &db{store, state, archive}
}

// NewStateValidator returns a Merkle validator for FigDB state
func NewStateValidator() StateValidator {
	hasher := &figcrypto.Hasher{}
	return &trie.StateValidator{Hasher: hasher}
}

// NewArchiveValidator returns a Merkle validator for FigDB archives
func NewArchiveValidator() ArchiveValidator {
	hasher := &figcrypto.Hasher{}
	return &trie.ArchiveValidator{Hasher: hasher}
}

type db struct {
	store   internal.KeyStore
	state   StateTrie
	archive ArchiveTrie
}

func (db *db) State() StateTrie {
	return db.state
}

func (db *db) Archive() ArchiveTrie {
	return db.archive
}

func (db *db) Batch() {
	db.store.Batch()
}

func (db *db) Write() {
	db.store.Write()
}

func (db *db) Close() {
	db.store.Close()
}

func (db *db) Backup(writer io.Writer, since uint64) (uint64, error) {
	if dbType == badgerDB {
		return db.store.(*badger.Store).DB.Backup(writer, since)
	}
	panic("not implemented")
}

func (db *db) Restore(reader io.Reader) error {
	if dbType == badgerDB {
		return db.store.(*badger.Store).DB.Load(reader)
	}
	panic("not implemented")
}

func (db *db) GC() {
	if dbType == badgerDB {
		db.store.(*badger.Store).GC()
	}
	panic("not implemented")
}

func newTries(store internal.KeyStore) (*trie.State, *trie.Archive) {
	hasher := &figcrypto.Hasher{}
	encdec := &mock.EncoderDecoder{} // obviously replace this with RLP
	state := &trie.State{KeyStore: store, Hasher: hasher, Encdec: encdec}
	archive := &trie.Archive{KeyStore: store, Hasher: hasher, Encdec: encdec}
	return state, archive
}
