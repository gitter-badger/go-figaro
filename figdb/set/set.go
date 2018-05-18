// Package set implements a probabilisitc archive set with a configurable false-positive rate.
// Items saved in the archive can be queried for existence, but cannot be retrieved. Each probabilistic archive
// is saved in the underlying key/value store under a hash of the set.
package set

// TODO: explore use of cuckoo filter instead... this will likely require creating
// our own implementation, as no solid Go implementation has been found yet

import (
	"github.com/figaro-tech/go-figaro/figcrypto/hash"
	"github.com/figaro-tech/go-figaro/figdb/filter/bbloom"
	"github.com/figaro-tech/go-figaro/figdb/types"
)

// Set impelements a probabilistic archive of set membership.
type Set struct {
	KeyStore types.KeyStore
	Cache    types.Cache
}

// Save creates a bloom filter of the members of data with the target
// false positivate rate, fp, returning a unique key for querying
// set membership in the future.
func (s *Set) Save(data [][]byte, fp float64) ([]byte, error) {
	bloom := bbloom.NewWithEstimates(uint64(len(data)), fp)
	for _, datum := range data {
		bloom.Add(datum)
	}
	v, err := bloom.Marshal()
	if err != nil {
		return nil, err
	}
	k := hash.Hash256(v)
	if s.Cache != nil {
		s.Cache.Add(k, v)
	}
	s.KeyStore.Set(k, v)
	return k, nil
}

// Get returns a bloom filter, intended for cases where multiple tests will
// occur on the same filter
func (s *Set) Get(key types.Key) (*bbloom.Bloom, error) {
	var v []byte
	var err error
	if s.Cache != nil {
		if c, ok := s.Cache.Get(key); ok {
			v = c
		}
	}
	if v == nil {
		v, err = s.KeyStore.Get(key)
	}
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, nil
	}
	bloom, err := bbloom.Unmarshal(v)
	if err != nil {
		return nil, err
	}
	return bloom, nil
}

// Has tests whether datum is in the Set.
func (s *Set) Has(key types.Key, datum []byte) bool {
	var v []byte
	var err error
	if s.Cache != nil {
		if c, ok := s.Cache.Get(key); ok {
			v = c
		}
	}
	if v == nil {
		v, err = s.KeyStore.Get(key)
	}
	if err != nil || len(v) == 0 {
		return false
	}
	bloom, err := bbloom.Unmarshal(v)
	if err != nil {
		return false
	}
	return bloom.Has(datum)
}

// HasBatch tests whether each datum in data is in the Set,
// returning an ordered []bool array of results.
func (s *Set) HasBatch(key types.Key, data [][]byte) (ins []bool) {
	var v []byte
	var err error
	if s.Cache != nil {
		if c, ok := s.Cache.Get(key); ok {
			v = c
		}
	}
	if v == nil {
		v, err = s.KeyStore.Get(key)
	}
	if err != nil || len(v) == 0 {
		return
	}
	bloom, err := bbloom.Unmarshal(v)
	if err != nil {
		return
	}
	ins = make([]bool, len(data))
	for i, datum := range data {
		ins[i] = bloom.Has(datum)
	}
	return
}
