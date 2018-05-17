// Package bbloom implements a high performance Bloom filter
package bbloom

import (
	"sync"
)

// BloomTS is a thread-safe Bloom filter
type BloomTS struct {
	lock sync.RWMutex
	Bloom
}

// NewTS returns a new bloomfilter sized for the
// number of entries and locations
func NewTS(entries, locs uint64) *BloomTS {
	bf := New(entries, locs)
	return &BloomTS{Bloom: *bf}
}

// NewTSWithEstimates returns a bloom filter suitable
// for the desire number of entries, n, and false
// positive rate, fp.
func NewTSWithEstimates(n uint64, fp float64) *BloomTS {
	bf := NewWithEstimates(n, fp)
	return &BloomTS{Bloom: *bf}
}

// NewTSWithBitset takes a []uint64 slice and number of locs per entry
// returns the bloomfilter with a bitset populated according to the input
func NewTSWithBitset(bs []uint64, locs uint64) *BloomTS {
	bf := NewWithBitset(bs, locs)
	return &BloomTS{Bloom: *bf}
}

// UnmarshalTS unmarshals a figbuf encoded Bloom filter into a Bloom filter
func UnmarshalTS(data []byte) (*BloomTS, error) {
	bf, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return &BloomTS{Bloom: *bf}, nil
}

// Add adds an entry to the Bloom filter
func (bl *BloomTS) Add(entry []byte) {
	bl.lock.Lock()
	defer bl.lock.Unlock()

	bl.Bloom.Add(entry)
}

// Has checks whether the Bloom filter contains the entry
func (bl *BloomTS) Has(entry []byte) bool {
	bl.lock.RLock()
	defer bl.lock.RUnlock()

	return bl.Bloom.Has(entry)
}

// Clear resets the Bloom filter
func (bl *BloomTS) Clear() {
	bl.lock.Lock()
	defer bl.lock.Unlock()

	bl.Bloom.Clear()
}

// Marshal returns figbuf encoded (type bloomImExport) as []byte
func (bl *BloomTS) Marshal() (buf []byte, err error) {
	bl.lock.RLock()
	defer bl.lock.RUnlock()

	return bl.Bloom.Marshal()
}
