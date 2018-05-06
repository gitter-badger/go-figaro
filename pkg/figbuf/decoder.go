// Package figbuf implements Recursive Length Prefix deterministic binary encoding
package figbuf

import (
	"encoding/binary"
	"errors"
	"log"
	"sync"
)

var (
	// DecoderPool is a global, thread-safe, reusable pool of Decoders
	//
	//   Use `DecoderPool.Get()` to get a Decoder, and
	//	`DecoderPool.Put(*Decoder)` to return an Decoder to the pool
	DecoderPool = sync.Pool{
		New: func() interface{} {
			return &Decoder{}
		},
	}

	// ErrInvalidData raised when attemptimg to decode data that is not well-known
	ErrInvalidData = errors.New("figbuf: invalid data for decoding must be well-known type")
)

// DeterministicBinaryUnmarshaler should unmarshal itself using the decoding helpers,
// allowing for deterministic decoding of complex types
type DeterministicBinaryUnmarshaler interface {
	UnmarshalDeterministicBinary(dec *Decoder) ([]byte, error)
}

type rlpType uint

const (
	_ rlpType = iota
	rlpStr
	rlpList
)

type rlpItem struct {
	offset uint
	length uint
	typ    rlpType
}

// Decoder is an RLP decoder
type Decoder struct {
	// Cache of tmp values
	length       uint
	prefix       uint
	lenOfStrLen  uint
	strLen       uint
	lenOfListLen uint
	listLen      uint
	zeros        [8]byte

	// items get reused across decodings
	items [32]rlpItem
	at    uint

	// bs17s do not get reused across decodings,
	// but we allocate them in batches for efficiency
	bs17 [32][17][]byte
	bsat uint
}

// DecodeBytes decodes
//
// Note that the slice returned may use the same
// backing array as `b` for performance
func (dec *Decoder) DecodeBytes(b []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		d = nil
		err = ErrInvalidData
		return
	}
	d = substr(b, item.offset, item.length)
	return
}

// DecodeBytesSlice decodes
//
// Note that the slice returned may use the same
// backing array as `bb` for performance
func (dec *Decoder) DecodeBytesSlice(bb []byte) (dd [][]byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	// 17 is the size of a merkle node, which is
	// the most common use case for a bytes slice
	dd = dec.getByteSlice17()
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, substr(b, item.offset, item.length))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeString decodes
func (dec *Decoder) DecodeString(b []byte) (d string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = ""
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
		return
	}
	d = dec.BytesToString(substr(b, item.offset, item.length))
	return
}

// DecodeInt decodes
func (dec *Decoder) DecodeInt(b []byte) (d int, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
		return
	}
	d = dec.BytesToInt(substr(b, item.offset, item.length))
	return
}

// DecodeInt8 decodes
func (dec *Decoder) DecodeInt8(b []byte) (d int8, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
		return
	}
	d = dec.BytesToInt8(substr(b, item.offset, item.length))
	return
}

// DecodeInt16 decodes
func (dec *Decoder) DecodeInt16(b []byte) (d int16, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
		return
	}
	d = dec.BytesToInt16(substr(b, item.offset, item.length))
	return
}

// DecodeInt32 decodes
func (dec *Decoder) DecodeInt32(b []byte) (d int32, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
		return
	}
	d = dec.BytesToInt32(substr(b, item.offset, item.length))
	return
}

// DecodeInt64 decodes
func (dec *Decoder) DecodeInt64(b []byte) (d int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
		return
	}
	d = dec.BytesToInt64(substr(b, item.offset, item.length))
	return
}

// DecodeUint decodes
func (dec *Decoder) DecodeUint(b []byte) (d uint, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
		return
	}
	d = dec.BytesToUint(substr(b, item.offset, item.length))
	return
}

// DecodeUint8 decodes
func (dec *Decoder) DecodeUint8(b []byte) (d uint8, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
		return
	}
	d = dec.BytesToUint8(substr(b, item.offset, item.length))
	return
}

// DecodeUint16 decodes
func (dec *Decoder) DecodeUint16(b []byte) (d uint16, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
	}
	d = dec.BytesToUint16(substr(b, item.offset, item.length))
	return
}

// DecodeUint32 decodes
func (dec *Decoder) DecodeUint32(b []byte) (d uint32, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
	}
	d = dec.BytesToUint32(substr(b, item.offset, item.length))
	return
}

// DecodeUint64 decodes
func (dec *Decoder) DecodeUint64(b []byte) (d uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	item := dec.nextItem(b)
	if item.typ != rlpStr {
		err = ErrInvalidData
	}
	d = dec.BytesToUint64(substr(b, item.offset, item.length))
	return
}

// DecodeStringSlice decodes
func (dec *Decoder) DecodeStringSlice(bb []byte) (dd []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]string, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToString(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeIntSlice decodes
func (dec *Decoder) DecodeIntSlice(bb []byte) (dd []int, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]int, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToInt(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeInt8Slice decodes
func (dec *Decoder) DecodeInt8Slice(bb []byte) (dd []int8, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]int8, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToInt8(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeInt16Slice decodes
func (dec *Decoder) DecodeInt16Slice(bb []byte) (dd []int16, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]int16, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToInt16(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeInt32Slice decodes
func (dec *Decoder) DecodeInt32Slice(bb []byte) (dd []int32, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]int32, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToInt32(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeInt64Slice decodes
func (dec *Decoder) DecodeInt64Slice(bb []byte) (dd []int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]int64, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		dd = append(dd, dec.BytesToInt64(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeUintSlice decodes
func (dec *Decoder) DecodeUintSlice(bb []byte) (dd []uint, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]uint, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToUint(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeUint8Slice decodes
func (dec *Decoder) DecodeUint8Slice(bb []byte) (dd []uint8, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]uint8, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToUint8(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeUint16Slice decodes
func (dec *Decoder) DecodeUint16Slice(bb []byte) (dd []uint16, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]uint16, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToUint16(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeUint32Slice decodes
func (dec *Decoder) DecodeUint32Slice(bb []byte) (dd []uint32, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]uint32, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToUint32(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// DecodeUint64Slice decodes
func (dec *Decoder) DecodeUint64Slice(bb []byte) (dd []uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec.at = 0
	dd = make([]uint64, 0, 8)
	item := dec.nextItem(bb)
	if item.typ != rlpList {
		dd = nil
		err = ErrInvalidData
		return
	}
	b := substr(bb, item.offset, item.length)
	for {
		item = dec.nextItem(b)
		if item.typ == 0 {
			break
		}
		if item.typ != rlpStr {
			dd = nil
			err = ErrInvalidData
			return
		}
		dd = append(dd, dec.BytesToUint64(substr(b, item.offset, item.length)))
		b = b[item.offset+item.length:]
	}
	return
}

// BytesToString converts
func (dec *Decoder) BytesToString(b []byte) string {
	return string(b)
}

// BytesToInt converts
func (dec *Decoder) BytesToInt(b []byte) int {
	i, n := binary.Varint(b)
	if n < 1 {
		panic(ErrInvalidData)
	}
	return int(i)
}

// BytesToInt8 converts
func (dec *Decoder) BytesToInt8(b []byte) int8 {
	i, n := binary.Varint(b)
	if n < 1 {
		panic(ErrInvalidData)
	}
	return int8(i)
}

// BytesToInt16 converts
func (dec *Decoder) BytesToInt16(b []byte) int16 {
	i, n := binary.Varint(b)
	if n < 1 {
		panic(ErrInvalidData)
	}
	return int16(i)
}

// BytesToInt32 converts
func (dec *Decoder) BytesToInt32(b []byte) int32 {
	i, n := binary.Varint(b)
	if n < 1 {
		panic(ErrInvalidData)
	}
	return int32(i)
}

// BytesToInt64 converts
func (dec *Decoder) BytesToInt64(b []byte) int64 {
	i, n := binary.Varint(b)
	if n < 1 {
		panic(ErrInvalidData)
	}
	return int64(i)
}

// BytesToUint converts
func (dec *Decoder) BytesToUint(b []byte) uint {
	b = append(dec.zeros[:8-len(b)], b...)
	return uint(binary.BigEndian.Uint64(b))
}

// BytesToUint8 converts
func (dec *Decoder) BytesToUint8(b []byte) uint8 {
	if len(b) == 0 {
		return 0
	}
	return b[0]
}

// BytesToUint16 converts
func (dec *Decoder) BytesToUint16(b []byte) uint16 {
	b = append(dec.zeros[:2-len(b)], b...)
	return binary.BigEndian.Uint16(b)
}

// BytesToUint32 converts
func (dec *Decoder) BytesToUint32(b []byte) uint32 {
	b = append(dec.zeros[:4-len(b)], b...)
	return binary.BigEndian.Uint32(b)
}

// BytesToUint64 converts
func (dec *Decoder) BytesToUint64(b []byte) uint64 {
	b = append(dec.zeros[:8-len(b)], b...)
	return binary.BigEndian.Uint64(b)
}

func (dec *Decoder) nextItem(b []byte) *rlpItem {
	dec.length = uint(len(b))
	if dec.length == 0 {
		return dec.getItem(0, 0, 0)
	}
	dec.prefix = uint(b[0])
	if dec.prefix <= 0x7f {
		return dec.getItem(0, 1, rlpStr)
	}
	if dec.prefix <= 0xb7 && dec.length > dec.prefix-0x80 {
		dec.strLen = dec.prefix - 0x80
		return dec.getItem(1, dec.strLen, rlpStr)
	}
	if dec.prefix <= 0xbf && dec.length > dec.prefix-0xb7 && dec.length > dec.prefix-0xb7+dec.BytesToUint(substr(b, 1, dec.prefix-0xb7)) {
		dec.lenOfStrLen = dec.prefix - 0xb7
		dec.strLen = dec.BytesToUint(substr(b, 1, dec.lenOfStrLen))
		return dec.getItem(1+dec.lenOfStrLen, dec.strLen, rlpStr)
	}
	if dec.prefix <= 0xf7 && dec.length > dec.prefix-0xc0 {
		dec.listLen = dec.prefix - 0xc0
		return dec.getItem(1, dec.listLen, rlpList)
	}
	if dec.prefix <= 0xff && dec.length > dec.prefix-0xf7 && dec.length > dec.prefix-0xf7+dec.BytesToUint(substr(b, 1, dec.prefix-0xf7)) {
		dec.lenOfListLen = dec.prefix - 0xf7
		dec.listLen = dec.BytesToUint(substr(b, 1, dec.lenOfListLen))
		return dec.getItem(1+dec.lenOfListLen, dec.listLen, rlpList)
	}
	panic(ErrInvalidData)
}

func (dec *Decoder) getItem(offset uint, length uint, typ rlpType) (i *rlpItem) {
	if dec.at == uint(len(dec.items)) {
		// If we're at the end of the pool,
		// we allocate some more and let the GC
		// collect the old ones
		var fl [32]rlpItem
		dec.items = fl
		dec.at = 0
	}
	// Grab one from the pool and set the values
	i = &dec.items[dec.at]
	i.offset = offset
	i.length = length
	i.typ = typ
	dec.at++
	return
}

func (dec *Decoder) getByteSlice17() [][]byte {
	if dec.bsat == uint(len(dec.bs17)) {
		// If we're at the end of the pool,
		// we allocate some more and let the GC
		// collect the old ones
		var fl [32][17][]byte
		dec.bs17 = fl
		dec.bsat = 0
	}
	// Grab one from the pool and set the values
	bs := &dec.bs17[dec.bsat]
	dec.bsat++
	return bs[:0]
}

func substr(b []byte, o, l uint) []byte {
	if o > uint(len(b)) {
		return b[len(b):]
	}
	if o+l > o+uint(len(b)) {
		return b[o:]
	}
	return b[o : o+l]
}
