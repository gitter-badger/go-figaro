// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"github.com/figaro-tech/figaro"
	"github.com/figaro-tech/figaro/pkg/figdb"
)

// DB is a domain Merkle database
type DB struct {
	DB      figdb.FigDB
	State   figdb.StateTrie
	Archive figdb.ArchiveTrie
	EncDec  figaro.EncoderDecoder
}

// Validator is a domain Merkle validator
type Validator struct {
	State   figdb.StateValidator
	Archive figdb.ArchiveValidator
	EncDec  figaro.EncoderDecoder
}

// New returns a FigDB backed by a high-performance disk database
func New(dir string, encdec figaro.EncoderDecoder) *DB {
	db := figdb.New(dir)
	return &DB{db, db.State(), db.Archive(), encdec}
}

// NewMemDB returns a FigDB backed by a high-performance memory database
func NewMemDB(encdec figaro.EncoderDecoder) *DB {
	db := figdb.NewMemDB()
	return &DB{db, db.State(), db.Archive(), encdec}
}

// NewValidator returns a Merkle validator ready to use
func NewValidator(encdec figaro.EncoderDecoder) *Validator {
	sv := figdb.NewStateValidator()
	av := figdb.NewArchiveValidator()
	return &Validator{sv, av, encdec}
}
