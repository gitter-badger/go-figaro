// Package figencoding implements Recursive Length Prefix deterministic binary encoding
package figencoding

import (
	"encoding"
	"encoding/binary"
	"errors"
	"log"
)

// ErrInvalidType raised when attemptimg to encode a type that is not well-known
var ErrInvalidType = errors.New("figencoding: invalid type for encoding must be well-known type")

// DeterministicBinaryMarshaler should marshal itself using the encoding helpers,
// allowing for deterministic encoding of complex types
type DeterministicBinaryMarshaler interface {
	MarshalDeterministicBinary(enc *Encoder) ([]byte, error)
}

// Encoder is an RLP encoder
type Encoder struct {
	buf   [1000]byte
	pad   [1000]byte
	buf10 [10]byte
	buf8  [8]byte
	buf5  [5]byte
	buf4  [4]byte
	buf3  [3]byte
	buf2  [2]byte
}

// Encode encodes well-known types into deterministic byte slices
//
// Returned values from Encoder are only safe to use during the function
// that called them and only until the encoder is used again
func (enc *Encoder) Encode(d interface{}) (b []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	return enc.encode(enc.buf[:0], d), nil
}

func (enc *Encoder) encode(buf []byte, d interface{}) []byte {
	switch d.(type) {
	case []byte:
		return enc.encodeBytes(buf, d.([]byte))
	case [][]byte:
		return enc.encodeBytesSlice(buf, d.([][]byte))
	// case [][][]byte:
	// 	return enc.encodeBytesSliceSlice(buf, d.([][][]byte))
	case []interface{}:
		return enc.encodeList(buf, d.([]interface{}))
	case *[]interface{}:
		return enc.encodeList(buf, *d.(*[]interface{}))
	case string:
		return enc.encodeString(buf, d.(string))
	case int:
		return enc.encodeInt(buf, d.(int))
	case int8:
		return enc.encodeInt8(buf, d.(int8))
	case int16:
		return enc.encodeInt16(buf, d.(int16))
	case int32:
		return enc.encodeInt32(buf, d.(int32))
	case int64:
		return enc.encodeInt64(buf, d.(int64))
	case uint:
		return enc.encodeUint(buf, d.(uint))
	case uint8:
		return enc.encodeUint8(buf, d.(uint8))
	case uint16:
		return enc.encodeUint16(buf, d.(uint16))
	case uint32:
		return enc.encodeUint32(buf, d.(uint32))
	case uint64:
		return enc.encodeUint64(buf, d.(uint64))
	case []string:
		return enc.encodeStringSlice(buf, d.([]string))
	case []int:
		return enc.encodeIntSlice(buf, d.([]int))
	case []int8:
		return enc.encodeInt8Slice(buf, d.([]int8))
	case []int16:
		return enc.encodeInt16Slice(buf, d.([]int16))
	case []int32:
		return enc.encodeInt32Slice(buf, d.([]int32))
	case []int64:
		return enc.encodeInt64Slice(buf, d.([]int64))
	case []uint:
		return enc.encodeUintSlice(buf, d.([]uint))
	case []uint16:
		return enc.encodeUint16Slice(buf, d.([]uint16))
	case []uint32:
		return enc.encodeUint32Slice(buf, d.([]uint32))
	case []uint64:
		return enc.encodeUint64Slice(buf, d.([]uint64))
	case DeterministicBinaryMarshaler:
		m, err := d.(DeterministicBinaryMarshaler).MarshalDeterministicBinary(enc)
		if err != nil {
			panic(err)
		}
		return m
	case encoding.BinaryMarshaler:
		m, err := d.(encoding.BinaryMarshaler).MarshalBinary()
		if err != nil {
			panic(err)
		}
		return enc.encodeBytes(buf, m)
	case encoding.TextMarshaler:
		m, err := d.(encoding.TextMarshaler).MarshalText()
		if err != nil {
			panic(err)
		}
		return enc.encodeBytes(buf, m)
	}
	panic(ErrInvalidType)
}

// EncodeBytes RLP encodes
func (enc *Encoder) EncodeBytes(d []byte) []byte {
	buf := enc.buf[:0]
	return enc.encodeRLPString(buf, d)
}

// EncodeBytesSlice RLP encodes slice
func (enc *Encoder) EncodeBytesSlice(dd [][]byte) []byte {
	buf := enc.buf[:0]
	return enc.encodeBytesSlice(buf, dd)
}

// // EncodeBytesSliceSlice RLP encodes
// func (enc *Encoder) EncodeBytesSliceSlice(ddd [][][]byte) []byte {
// 	buf := enc.buf[:0]
// 	return enc.encodeBytesSlicSlice(buf, dd)
// }

// EncodeList RLP encodes
func (enc *Encoder) EncodeList(d ...interface{}) []byte {
	buf := enc.buf[:0]
	return enc.encodeRLPList(buf, d)
}

// EncodeString RLP encodes
func (enc *Encoder) EncodeString(d string) []byte {
	buf := enc.buf[:0]
	return enc.encodeRLPString(buf, []byte(d))
}

// EncodeInt RLP encodes
func (enc *Encoder) EncodeInt(d int) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf10[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeInt8 RLP encodes
func (enc *Encoder) EncodeInt8(d int8) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf2[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeInt16 RLP encodes
func (enc *Encoder) EncodeInt16(d int16) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf3[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeInt32 RLP encodes
func (enc *Encoder) EncodeInt32(d int32) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf5[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeInt64 RLP encodes
func (enc *Encoder) EncodeInt64(d int64) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf10[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeUint RLP encodes
func (enc *Encoder) EncodeUint(d uint) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf8[:]
	binary.BigEndian.PutUint64(b, uint64(d))
	bl := binaryLen(d)
	return enc.encodeRLPString(buf, b[len(b)-int(bl):])
}

// EncodeUint8 RLP encodes
func (enc *Encoder) EncodeUint8(d uint8) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	return enc.encodeRLPString(buf, []byte{d})
}

// EncodeUint16 RLP encodes
func (enc *Encoder) EncodeUint16(d uint16) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf2[:]
	binary.BigEndian.PutUint16(b, uint16(d))
	bl := binaryLen(uint(d))
	return enc.encodeRLPString(buf, b[len(b)-int(bl):])
}

// EncodeUint32 RLP encodes
func (enc *Encoder) EncodeUint32(d uint32) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf4[:]
	binary.BigEndian.PutUint32(b, uint32(d))
	bl := binaryLen(uint(d))
	return enc.encodeRLPString(buf, b[len(b)-int(bl):])
}

// EncodeUint64 RLP encodes
func (enc *Encoder) EncodeUint64(d uint64) []byte {
	buf := enc.buf[:0]
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf8[:]
	binary.BigEndian.PutUint64(b, uint64(d))
	bl := binaryLen(uint(d))
	return enc.encodeRLPString(buf, b[len(b)-int(bl):])
}

// EncodeStringSlice RLP encodes
func (enc *Encoder) EncodeStringSlice(dd []string) []byte {
	buf := enc.buf[:0]
	return enc.encodeStringSlice(buf, dd)
}

// EncodeIntSlice RLP encodes
func (enc *Encoder) EncodeIntSlice(dd []int) []byte {
	buf := enc.buf[:0]
	return enc.encodeIntSlice(buf, dd)
}

// EncodeInt8Slice RLP encodes
func (enc *Encoder) EncodeInt8Slice(dd []int8) []byte {
	buf := enc.buf[:0]
	return enc.encodeInt8Slice(buf, dd)
}

// EncodeInt16Slice RLP encodes
func (enc *Encoder) EncodeInt16Slice(dd []int16) []byte {
	buf := enc.buf[:0]
	return enc.encodeInt16Slice(buf, dd)
}

// EncodeInt32Slice RLP encodes
func (enc *Encoder) EncodeInt32Slice(dd []int32) []byte {
	buf := enc.buf[:0]
	return enc.encodeInt32Slice(buf, dd)
}

// EncodeInt64Slice RLP encodes
func (enc *Encoder) EncodeInt64Slice(dd []int64) []byte {
	buf := enc.buf[:0]
	return enc.encodeInt64Slice(buf, dd)
}

// EncodeUintSlice RLP encodes
func (enc *Encoder) EncodeUintSlice(dd []uint) []byte {
	buf := enc.buf[:0]
	return enc.encodeUintSlice(buf, dd)
}

// EncodeUint8Slice RLP encodes
func (enc *Encoder) EncodeUint8Slice(dd []uint8) []byte {
	buf := enc.buf[:0]
	return enc.encodeRLPString(buf, dd)
}

// EncodeUint16Slice RLP encodes
func (enc *Encoder) EncodeUint16Slice(dd []uint16) []byte {
	buf := enc.buf[:0]
	return enc.encodeUint16Slice(buf, dd)
}

// EncodeUint32Slice RLP encodes
func (enc *Encoder) EncodeUint32Slice(dd []uint32) []byte {
	buf := enc.buf[:0]
	return enc.encodeUint32Slice(buf, dd)
}

// EncodeUint64Slice RLP encodes
func (enc *Encoder) EncodeUint64Slice(dd []uint64) []byte {
	buf := enc.buf[:0]
	return enc.encodeUint64Slice(buf, dd)
}

// EncodeBytes RLP encodes
func (enc *Encoder) encodeBytes(buf []byte, d []byte) []byte {
	return enc.encodeRLPString(buf, d)
}

// EncodeBytesSlice RLP encodes slice
func (enc *Encoder) encodeBytesSlice(buf []byte, dd [][]byte) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeBytes(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// // EncodeBytesSliceSlice RLP encodes
// func (enc *Encoder) encodeBytesSliceSlice(buf []byte, ddd [][][]byte) []byte {
// 	idx := uint(len(buf))
// 	for _, dd := range ddd {
// 		buf = enc.encodeBytesSlice(buf, dd)
// 	}
// 	l := uint(len(buf)) - idx
// 	var pad uint
// 	if l < 56 {
// 		pad = 1
// 	} else {
// 		pad = 1 + binaryLen(l)
// 	}
// 	padding := enc.pad[:pad]
// 	buf = append(buf, padding...)
// 	copy(buf[idx+pad:], buf[idx:])
// 	buf = buf[:idx]
// 	buf = enc.encodeRLPLength(buf, l, 0xc0)
// 	buf = buf[:len(buf)+int(l)]
// 	return buf
// }

// EncodeList RLP encodes
func (enc *Encoder) encodeList(buf []byte, dd []interface{}) []byte {
	return enc.encodeRLPList(buf, dd)
}

// EncodeString RLP encodes
func (enc *Encoder) encodeString(buf []byte, d string) []byte {
	return enc.encodeRLPString(buf, []byte(d))
}

// EncodeInt RLP encodes
func (enc *Encoder) encodeInt(buf []byte, d int) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf10[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeInt8 RLP encodes
func (enc *Encoder) encodeInt8(buf []byte, d int8) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf2[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeInt16 RLP encodes
func (enc *Encoder) encodeInt16(buf []byte, d int16) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf3[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeInt32 RLP encodes
func (enc *Encoder) encodeInt32(buf []byte, d int32) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf5[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeInt64 RLP encodes
func (enc *Encoder) encodeInt64(buf []byte, d int64) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf10[:]
	n := binary.PutVarint(b, int64(d))
	return enc.encodeRLPString(buf, b[:n])
}

// EncodeUint RLP encodes
func (enc *Encoder) encodeUint(buf []byte, d uint) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf8[:]
	binary.BigEndian.PutUint64(b, uint64(d))
	bl := binaryLen(d)
	return enc.encodeRLPString(buf, b[len(b)-int(bl):])
}

// EncodeUint8 RLP encodes
func (enc *Encoder) encodeUint8(buf []byte, d uint8) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	return enc.encodeRLPString(buf, []byte{d})
}

// EncodeUint16 RLP encodes
func (enc *Encoder) encodeUint16(buf []byte, d uint16) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf2[:]
	binary.BigEndian.PutUint16(b, uint16(d))
	bl := binaryLen(uint(d))
	return enc.encodeRLPString(buf, b[len(b)-int(bl):])
}

// EncodeUint32 RLP encodes
func (enc *Encoder) encodeUint32(buf []byte, d uint32) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf4[:]
	binary.BigEndian.PutUint32(b, uint32(d))
	bl := binaryLen(uint(d))
	return enc.encodeRLPString(buf, b[len(b)-int(bl):])
}

// EncodeUint64 RLP encodes
func (enc *Encoder) encodeUint64(buf []byte, d uint64) []byte {
	if d == 0 {
		return enc.encodeRLPString(buf, nil)
	}
	b := enc.buf8[:]
	binary.BigEndian.PutUint64(b, uint64(d))
	bl := binaryLen(uint(d))
	return enc.encodeRLPString(buf, b[len(b)-int(bl):])
}

// EncodeStringSlice RLP encodes
func (enc *Encoder) encodeStringSlice(buf []byte, dd []string) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeString(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// EncodeIntSlice RLP encodes
func (enc *Encoder) encodeIntSlice(buf []byte, dd []int) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeInt(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// EncodeInt8Slice RLP encodes
func (enc *Encoder) encodeInt8Slice(buf []byte, dd []int8) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeInt8(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// EncodeInt16Slice RLP encodes
func (enc *Encoder) encodeInt16Slice(buf []byte, dd []int16) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeInt16(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// EncodeInt32Slice RLP encodes
func (enc *Encoder) encodeInt32Slice(buf []byte, dd []int32) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeInt32(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// EncodeInt64Slice RLP encodes
func (enc *Encoder) encodeInt64Slice(buf []byte, dd []int64) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeInt64(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// EncodeUintSlice RLP encodes
func (enc *Encoder) encodeUintSlice(buf []byte, dd []uint) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeUint(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// EncodeUint8Slice RLP encodes
func (enc *Encoder) encodeUint8Slice(buf []byte, dd []uint8) []byte {
	return enc.encodeBytes(buf, dd)
}

// EncodeUint16Slice RLP encodes
func (enc *Encoder) encodeUint16Slice(buf []byte, dd []uint16) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeUint16(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// EncodeUint32Slice RLP encodes
func (enc *Encoder) encodeUint32Slice(buf []byte, dd []uint32) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeUint32(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

// EncodeUint64Slice RLP encodes
func (enc *Encoder) encodeUint64Slice(buf []byte, dd []uint64) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeUint64(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

func (enc *Encoder) encodeRLPString(buf []byte, s []byte) []byte {
	if len(s) == 1 && s[0] < 0x80 {
		buf = append(buf, s[0])
		return buf
	}
	buf = enc.encodeRLPLength(buf, uint(len(s)), 0x80)
	buf = append(buf, s...)
	return buf
}

func (enc *Encoder) encodeRLPList(buf []byte, dd []interface{}) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encode(buf, d)
	}
	return enc.encodeListHelper(buf, idx)
}

func (enc *Encoder) encodeListHelper(buf []byte, idx uint) []byte {
	l := uint(len(buf)) - idx
	var pad uint
	if l < 56 {
		pad = 1
	} else {
		pad = 1 + binaryLen(l)
	}
	padding := enc.pad[:pad]
	buf = append(buf, padding...)
	copy(buf[idx+pad:], buf[idx:])
	buf = buf[:idx]
	buf = enc.encodeRLPLength(buf, l, 0xc0)
	buf = buf[:len(buf)+int(l)]
	return buf
}

func (enc *Encoder) encodeRLPLength(buf []byte, l uint, off uint) []byte {
	if l < 56 {
		buf = append(buf, byte(l+off))
		return buf
	}
	bl := binaryLen(l)
	var b []byte
	if bl+1 > 4 {
		b := enc.buf8[:]
		binary.BigEndian.PutUint64(b, uint64(l))
	}
	if bl+1 > 2 {
		b := enc.buf4[:]
		binary.BigEndian.PutUint32(b, uint32(l))

	} else {
		b := enc.buf2[:]
		binary.BigEndian.PutUint16(b, uint16(l))
	}
	buf = append(buf, byte(bl+off+55))
	buf = enc.encodeUint(buf, l)
	buf = append(buf, b...)
	return buf
}

var tab32 = [32]uint{
	0, 9, 1, 10, 13, 21, 2, 29, 11, 14, 16, 18, 22, 25, 3, 30,
	8, 12, 20, 28, 15, 17, 24, 7, 19, 27, 23, 6, 26, 5, 4, 31}

func binaryLen(u uint) uint {
	if u == 0 {
		return 0
	}
	v := uint32(u)
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	nbits := tab32[(v*0x07C4ACDD)>>27] + 1
	if nbits%8 == 0 {
		return nbits / 8
	}
	return nbits/8 + 1
}
