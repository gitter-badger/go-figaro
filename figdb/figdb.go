// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"github.com/figaro-tech/figaro"
	"github.com/figaro-tech/figaro/pkg/figdb"
)

// DB is a domain Merkle database
type DB struct {
	DB figdb.FigDB
	SV figdb.StateValidator
	AV figdb.ArchiveValidator
	H  figaro.Hasher
	ED figaro.EncoderDecoder
}

// New returns a FigDB backed by a high-performance disk database
func New(dir string, hasher figaro.Hasher, encdec figaro.EncoderDecoder) *DB {
	db := figdb.New(dir, hasher, encdec)
	sv, av := newValidators(hasher)
	return &DB{db, sv, av, hasher, encdec}
}

// NewMemDB returns a FigDB backed by a high-performance memory database
func NewMemDB(hasher figaro.Hasher, encdec figaro.EncoderDecoder) *DB {
	db := figdb.NewMemDB(hasher, encdec)
	sv, av := newValidators(hasher)
	return &DB{db, sv, av, hasher, encdec}
}

func newValidators(hasher figaro.Hasher) (figdb.StateValidator, figdb.ArchiveValidator) {
	sv := figdb.NewStateValidator(hasher)
	av := figdb.NewArchiveValidator(hasher)
	return sv, av
}
