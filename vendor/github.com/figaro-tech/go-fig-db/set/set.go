// Package set implements a probabilisitc set with a configurable false-positive rate.
// Items saved in the archive can be queried for existence, but cannot be retrieved. Each probabilistic archive
// is saved in the underlying key/value store under a hash of the set. Package `bloom` can be used separately
// if one would rather save the binary representation of the bloom filter inside another structure directly.
package set

// TODO: explore use of cuckoo filter instead... this will likely require creating
// our own implementation, as no solid Go implementation has been found yet

import (
	"bytes"
	"errors"

	"github.com/figaro-tech/go-fig-crypto/hasher"
	"github.com/figaro-tech/go-fig-db/bloom"
	"github.com/figaro-tech/go-fig-db/types"
)

// Set impelements a probabilistic archive of set membership.
type Set struct {
	KeyStore types.KeyStore
	Cache    types.Cache
}

// Create creates a bloom filter of the members of data with the target
// false positivate rate, fp, returning a unique key for querying
// set membership in the future.
func (s *Set) Create(data [][]byte, fp float64) (key, set []byte, err error) {
	bloom := bloom.NewWithEstimates(uint64(len(data)), fp)
	for _, datum := range data {
		bloom.Add(datum)
	}
	set, err = bloom.Marshal()
	if err != nil {
		return
	}
	key = hasher.Hash256(set)
	if s.Cache != nil {
		s.Cache.Add(key, set)
	}
	s.KeyStore.Set(key, set)
	return
}

// Save saves a set under key directly, after validating they key.
func (s *Set) Save(key []byte, set []byte) error {
	if Validate(key, set) {
		return errors.New("invalid key for set")
	}
	if s.Cache != nil {
		s.Cache.Add(key, set)
	}
	s.KeyStore.Set(key, set)
	return nil
}

// Get returns the set in binary format.
func (s *Set) Get(key types.Key) (set []byte, err error) {
	if s.Cache != nil {
		if c, ok := s.Cache.Get(key); ok {
			set = c
		}
	}
	if set == nil {
		set, err = s.KeyStore.Get(key)
	}
	return
}

// GetBloom returns a bloom filter, intended for cases where multiple tests will
// occur on the same filter
func (s *Set) GetBloom(key types.Key) (*bloom.Bloom, error) {
	set, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	bloom, err := bloom.Unmarshal(set)
	if err != nil {
		return nil, err
	}
	return bloom, nil
}

// Has tests whether datum is in the Set.
func (s *Set) Has(key types.Key, datum []byte) bool {
	bloom, err := s.GetBloom(key)
	if err != nil {
		return false
	}
	return bloom.Has(datum)
}

// HasBatch tests whether each datum in data is in the Set,
// returning an ordered []bool array of results.
func (s *Set) HasBatch(key types.Key, data [][]byte) (ins []bool) {
	bloom, err := s.GetBloom(key)
	if err != nil {
		return
	}
	ins = make([]bool, len(data))
	for i, datum := range data {
		ins[i] = bloom.Has(datum)
	}
	return
}

// Validate validates that a key is valid for a given set.
func Validate(key []byte, set []byte) bool {
	k := hasher.Hash256(set)
	return bytes.Equal(key, k)
}
