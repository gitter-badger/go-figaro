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

// Package bloom implements a high performance Bloom filter
package bloom

import (
	"log"
	"math"
	"unsafe"

	"github.com/figaro-tech/go-figaro/figbuf"
)

// Bloom filter
type Bloom struct {
	bitset  []uint64
	sizeExp uint64
	size    uint64
	setLocs uint64
	shift   uint64
}

// New returns a new bloomfilter sized for the
// number of entries and locations
func New(entries, locs uint64) (bloom *Bloom) {
	size, exponent := calcSizeAndExponent(entries)
	bloom = &Bloom{
		sizeExp: exponent,
		size:    size - 1,
		setLocs: locs,
		shift:   64 - exponent,
		bitset:  make([]uint64, size>>6),
	}
	return
}

// NewWithEstimates returns a bloom filter suitable
// for the desire number of entries, n, and false
// positive rate, fp.
func NewWithEstimates(n uint64, fp float64) (bloom *Bloom) {
	size, locs := calcSizeByWrongPositives(n, fp)
	return New(size, locs)
}

// NewWithBitset takes a []uint64 slice and number of locs per entry
// returns the bloomfilter with a bitset populated according to the input
func NewWithBitset(bs []uint64, locs uint64) (bloom *Bloom) {
	size, exponent := calcSizeAndExponent(uint64(len(bs)))
	bloom = &Bloom{
		sizeExp: exponent,
		size:    size - 1,
		setLocs: locs,
		shift:   64 - exponent,
		bitset:  bs,
	}
	return
}

// Unmarshal unmarshals a figbuf encoded Bloom filter into a Bloom filter
func Unmarshal(data []byte) (bloom *Bloom, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				bloom = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	var locs uint64
	var filterset []uint64
	_ = dec.DecodeNextList(data, func(buf []byte) {
		locs, buf = dec.DecodeNextUint64(buf)
		filterset, _ = dec.DecodeNextUint64Slice(buf)
	})
	bloom = NewWithBitset(filterset, locs)
	return
}

// Add adds an entry to the Bloom filter
func (bl *Bloom) Add(entry []byte) {
	l, h := SipHash(entry, bl.shift)
	for i := uint64(0); i < bl.setLocs; i++ {
		bl.set((h + i*l) & bl.size)
	}
}

// Has checks whether the Bloom filter contains the entry
func (bl *Bloom) Has(entry []byte) bool {
	l, h := SipHash(entry, bl.shift)
	for i := uint64(0); i < bl.setLocs; i++ {
		if bl.isSet((h + i*l) & bl.size) {
			continue
		}
		return false
	}
	return true
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

	buf = enc.EncodeNextList(nil, func(buf []byte) []byte {
		buf = enc.EncodeNextUint64(buf, bl.setLocs)
		buf = enc.EncodeNextUint64Slice(buf, bl.bitset)
		return buf
	})
	return buf, nil
}

// Clear resets the Bloom filter
func (bl *Bloom) Clear() {
	for i := range bl.bitset {
		bl.bitset[i] = 0
	}
}

func (bl *Bloom) set(idx uint64) {
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(&bl.bitset[idx>>6])) + uintptr((idx%64)>>3))
	*(*uint8)(ptr) |= mask[idx%8]
}

func (bl *Bloom) isSet(idx uint64) bool {
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(&bl.bitset[idx>>6])) + uintptr((idx%64)>>3))
	r := ((*(*uint8)(ptr)) >> (idx % 8)) & 1
	return r == 1
}

// helper
var mask = []uint8{1, 2, 4, 8, 16, 32, 64, 128}

func calcSizeAndExponent(ui64 uint64) (size uint64, exponent uint64) {
	if ui64 < uint64(512) {
		ui64 = uint64(512)
	}
	size = uint64(1)
	for size < ui64 {
		size <<= 1
		exponent++
	}
	return
}

func calcSizeByWrongPositives(numEntries uint64, wrongs float64) (size uint64, locs uint64) {
	size = uint64(-1 * float64(numEntries) * math.Log(wrongs) / math.Pow(float64(0.69314718056), 2))
	locs = uint64(math.Ceil(float64(0.69314718056) * float64(size) / float64(numEntries)))
	return
}
