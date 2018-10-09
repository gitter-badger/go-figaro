// Package bloom implements a high performance Bloom filter
package bloom

import (
	"sync"
)

// ThreadSafe is a thread-ThreadSafe Bloom filter
type ThreadSafe struct {
	lock sync.RWMutex
	Bloom
}

// NewThreadSafe returns a new bloomfilter sized for the
// number of entries and locations
func NewThreadSafe(entries, locs uint64) *ThreadSafe {
	bf := New(entries, locs)
	return &ThreadSafe{Bloom: *bf}
}

// NewThreadSafeWithEstimates returns a bloom filter suitable
// for the desire number of entries, n, and false
// positive rate, fp.
func NewThreadSafeWithEstimates(n uint64, fp float64) *ThreadSafe {
	bf := NewWithEstimates(n, fp)
	return &ThreadSafe{Bloom: *bf}
}

// NewThreadSafeWithBitset takes a []uint64 slice and number of locs per entry
// returns the bloomfilter with a bitset populated according to the input
func NewThreadSafeWithBitset(bs []uint64, locs uint64) *ThreadSafe {
	bf := NewWithBitset(bs, locs)
	return &ThreadSafe{Bloom: *bf}
}

// UnmarshalToThreadSafe unmarshals a figbuf encoded Bloom filter into a Thread-Safe Bloom filter
func UnmarshalToThreadSafe(data []byte) (*ThreadSafe, error) {
	bf, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return &ThreadSafe{Bloom: *bf}, nil
}

// Add adds an entry to the Bloom filter
func (bl *ThreadSafe) Add(entry []byte) {
	bl.lock.Lock()
	defer bl.lock.Unlock()

	bl.Bloom.Add(entry)
}

// Has checks whether the Bloom filter contains the entry
func (bl *ThreadSafe) Has(entry []byte) bool {
	bl.lock.RLock()
	defer bl.lock.RUnlock()

	return bl.Bloom.Has(entry)
}

// Clear resets the Bloom filter
func (bl *ThreadSafe) Clear() {
	bl.lock.Lock()
	defer bl.lock.Unlock()

	bl.Bloom.Clear()
}

// Marshal returns figbuf encoded (type bloomImExport) as []byte
func (bl *ThreadSafe) Marshal() (buf []byte, err error) {
	bl.lock.RLock()
	defer bl.lock.RUnlock()

	return bl.Bloom.Marshal()
}
