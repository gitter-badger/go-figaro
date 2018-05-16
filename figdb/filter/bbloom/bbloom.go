// Copyright 2018 The Figaro Authors.
// <License goes here>
// Based on Andreaas Bries's BBloom https://github.com/AndreasBriese/bbloom
// <Original License>

// The MIT License (MIT)
// Copyright (c) 2014 Andreas Briese, eduToolbox@Bri-C GmbH, Sarstedt

// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package bbloom

import (
	"log"
	"math"
	"sync"
	"unsafe"

	"github.com/figaro-tech/go-figaro/figbuf"
)

// helper
var mask = []uint8{1, 2, 4, 8, 16, 32, 64, 128}

func getSize(ui64 uint64) (size uint64, exponent uint64) {
	if ui64 < uint64(512) {
		ui64 = uint64(512)
	}
	size = uint64(1)
	for size < ui64 {
		size <<= 1
		exponent++
	}
	return size, exponent
}

func calcSizeByWrongPositives(numEntries, wrongs float64) (uint64, uint64) {
	size := -1 * numEntries * math.Log(wrongs) / math.Pow(float64(0.69314718056), 2)
	locs := math.Ceil(float64(0.69314718056) * size / numEntries)
	return uint64(size), uint64(locs)
}

// New returns a new bloomfilter
func New(params ...float64) (bloomfilter *Bloom) {
	var entries, locs uint64
	if len(params) == 2 {
		if params[1] < 1 {
			entries, locs = calcSizeByWrongPositives(params[0], params[1])
		} else {
			entries, locs = uint64(params[0]), uint64(params[1])
		}
	} else {
		log.Fatal("usage: New(float64(number_of_entries), float64(number_of_hashlocations)) i.e. New(float64(1000), float64(3)) or New(float64(number_of_entries), float64(number_of_hashlocations)) i.e. New(float64(1000), float64(0.03))")
	}
	size, exponent := getSize(uint64(entries))
	bloomfilter = &Bloom{
		sizeExp: exponent,
		size:    size - 1,
		setLocs: locs,
		shift:   64 - exponent,
	}
	bloomfilter.Size(size)
	return bloomfilter
}

// NewWithBoolset takes a []byte slice and number of locs per entry
// returns the bloomfilter with a bitset populated according to the input []byte
func NewWithBoolset(bs *[]byte, locs uint64) (bloomfilter *Bloom) {
	bloomfilter = New(float64(len(*bs)<<3), float64(locs))
	ptr := uintptr(unsafe.Pointer(&bloomfilter.bitset[0]))
	for _, b := range *bs {
		*(*uint8)(unsafe.Pointer(ptr)) = b
		ptr++
	}
	return bloomfilter
}

// bloomImExport Im/Export structure used by JSONMarshal / JSONUnmarshal
type bloomImExport struct {
	SetLocs   uint64
	FilterSet []byte
}

// Unmarshal takes figbuf encoded (type bloomImExport) as []bytes
// returns bloom32 / bloom64 object
func Unmarshal(dbData []byte) (bloomfilter *Bloom, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				bloomfilter = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	bloomImEx := bloomImExport{}
	_ = dec.DecodeNextList(dbData, func(buf []byte) {
		bloomImEx.SetLocs, buf = dec.DecodeNextUint64(buf)
		bloomImEx.FilterSet, _ = dec.DecodeNextBytes(buf)
	})
	bf := NewWithBoolset(&bloomImEx.FilterSet, bloomImEx.SetLocs)
	return bf, nil
}

//
// Bloom filter
type Bloom struct {
	Mtx     sync.RWMutex
	ElemNum uint64
	bitset  []uint64
	sizeExp uint64
	size    uint64
	setLocs uint64
	shift   uint64
}

// Update: found sipHash of Jean-Philippe Aumasson & Daniel J. Bernstein to be even faster than absdbm()
// https://131002.net/siphash/
// siphash was implemented for Go by Dmitry Chestnykh https://github.com/dchest/siphash

// Add set the bit(s) for entry; Adds an entry to the Bloom filter
func (bl *Bloom) Add(entry []byte) {
	l, h := bl.sipHash(entry)
	for i := uint64(0); i < (*bl).setLocs; i++ {
		(*bl).Set((h + i*l) & (*bl).size)
		(*bl).ElemNum++
	}
}

// AddTS Thread safe: Mutex.Lock the bloomfilter for the time of processing the entry
func (bl *Bloom) AddTS(entry []byte) {
	bl.Mtx.Lock()
	defer bl.Mtx.Unlock()
	bl.Add(entry[:])
}

// Has check if bit(s) for entry is/are set
// returns true if the entry was added to the Bloom Filter
func (bl *Bloom) Has(entry []byte) bool {
	l, h := bl.sipHash(entry)
	for i := uint64(0); i < bl.setLocs; i++ {
		switch bl.IsSet((h + i*l) & bl.size) {
		case false:
			return false
		}
	}
	return true
}

// HasTS Thread safe: Mutex.Lock the bloomfilter for the time of processing the entry
func (bl *Bloom) HasTS(entry []byte) bool {
	bl.Mtx.RLock()
	defer bl.Mtx.RUnlock()
	return bl.Has(entry[:])
}

// AddIfNotHas Only Add entry if it's not present in the bloomfilter
// returns true if entry was added
// returns false if entry was allready registered in the bloomfilter
func (bl *Bloom) AddIfNotHas(entry []byte) (added bool) {
	if bl.Has(entry[:]) {
		return added
	}
	bl.Add(entry[:])
	return true
}

// AddIfNotHasTS Tread safe: Only Add entry if it's not present in the bloomfilter
// returns true if entry was added
// returns false if entry was allready registered in the bloomfilter
func (bl *Bloom) AddIfNotHasTS(entry []byte) (added bool) {
	bl.Mtx.Lock()
	defer bl.Mtx.Unlock()
	return bl.AddIfNotHas(entry[:])
}

// Size make Bloom filter with as bitset of size sz
func (bl *Bloom) Size(sz uint64) {
	(*bl).bitset = make([]uint64, sz>>6)
}

// Clear resets the Bloom filter
func (bl *Bloom) Clear() {
	for i := range (*bl).bitset {
		(*bl).bitset[i] = 0
	}
}

// Set set the bit[idx] of bitsit
func (bl *Bloom) Set(idx uint64) {
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(&bl.bitset[idx>>6])) + uintptr((idx%64)>>3))
	*(*uint8)(ptr) |= mask[idx%8]
}

// IsSet check if bit[idx] of bitset is set
// returns true/false
func (bl *Bloom) IsSet(idx uint64) bool {
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(&bl.bitset[idx>>6])) + uintptr((idx%64)>>3))
	r := ((*(*uint8)(ptr)) >> (idx % 8)) & 1
	return r == 1
}

// Marshal returns figbuf encoded (type bloomImExport) as []byte
func (bl *Bloom) Marshal() (buf []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				buf = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	bloomImEx := bloomImExport{}
	bloomImEx.SetLocs = uint64(bl.setLocs)
	bloomImEx.FilterSet = make([]byte, len(bl.bitset)<<3)
	ptr := uintptr(unsafe.Pointer(&bl.bitset[0]))
	for i := range bloomImEx.FilterSet {
		bloomImEx.FilterSet[i] = *(*byte)(unsafe.Pointer(ptr))
		ptr++
	}

	buf = enc.EncodeNextList(nil, func(buf []byte) []byte {
		buf = enc.EncodeNextUint64(buf, bloomImEx.SetLocs)
		buf = enc.EncodeNextBytes(buf, bloomImEx.FilterSet)
		return buf
	})
	return buf, nil
}
