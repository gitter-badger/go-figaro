// Package figbuf implements Recursive Length Prefix deterministic binary encoding
package figbuf

import (
	"encoding"
	"encoding/binary"
	"errors"
	"log"
	"sync"
)

var (
	// EncoderPool is a global, thread-safe, reusable pool of Encoders
	//
	//   Use `EncoderPool.Get()` to get an Encoder, and
	//	`EncoderPool.Put(*Encoder)` to return an Encoder to the pool
	EncoderPool = sync.Pool{
		New: func() interface{} {
			return new(Encoder)
		},
	}

	// ErrInvalidType raised when attemptimg to encode a type that is not well-known
	ErrInvalidType = errors.New("figbuf: invalid type for encoding must be well-known type")
)

// Encoder is an RLP encoder
type Encoder struct {
	buf   [4096]byte
	lbuf  [256][]byte
	pad   [9]byte
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
func (enc *Encoder) Encode(d ...interface{}) ([]byte, error) {
	buf := enc.buf[:0]
	b, err := enc.EncodeNext(buf, d...)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// EncodeBytes RLP encodes
func (enc *Encoder) EncodeBytes(d []byte) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextBytes(buf, d)
}

// EncodeBytesSlice RLP encodes slice
func (enc *Encoder) EncodeBytesSlice(dd [][]byte) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextBytesSlice(buf, dd)
}

// EncodeString RLP encodes
func (enc *Encoder) EncodeString(d string) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextString(buf, d)
}

// EncodeBool RLP encodes
func (enc *Encoder) EncodeBool(d bool) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextBool(buf, d)
}

// EncodeInt RLP encodes
func (enc *Encoder) EncodeInt(d int) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextInt(buf, d)
}

// EncodeInt8 RLP encodes
func (enc *Encoder) EncodeInt8(d int8) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextInt8(buf, d)
}

// EncodeInt16 RLP encodes
func (enc *Encoder) EncodeInt16(d int16) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextInt16(buf, d)
}

// EncodeInt32 RLP encodes
func (enc *Encoder) EncodeInt32(d int32) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextInt32(buf, d)
}

// EncodeInt64 RLP encodes
func (enc *Encoder) EncodeInt64(d int64) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextInt64(buf, d)
}

// EncodeUint RLP encodes
func (enc *Encoder) EncodeUint(d uint) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUint(buf, d)
}

// EncodeUint8 RLP encodes
func (enc *Encoder) EncodeUint8(d uint8) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUint8(buf, d)
}

// EncodeUint16 RLP encodes
func (enc *Encoder) EncodeUint16(d uint16) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUint16(buf, d)
}

// EncodeUint32 RLP encodes
func (enc *Encoder) EncodeUint32(d uint32) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUint32(buf, d)
}

// EncodeUint64 RLP encodes
func (enc *Encoder) EncodeUint64(d uint64) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUint64(buf, d)
}

// EncodeStringSlice RLP encodes
func (enc *Encoder) EncodeStringSlice(dd []string) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextStringSlice(buf, dd)
}

// EncodeIntSlice RLP encodes
func (enc *Encoder) EncodeIntSlice(dd []int) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextIntSlice(buf, dd)
}

// EncodeInt8Slice RLP encodes
func (enc *Encoder) EncodeInt8Slice(dd []int8) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextInt8Slice(buf, dd)
}

// EncodeInt16Slice RLP encodes
func (enc *Encoder) EncodeInt16Slice(dd []int16) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextInt16Slice(buf, dd)
}

// EncodeInt32Slice RLP encodes
func (enc *Encoder) EncodeInt32Slice(dd []int32) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextInt32Slice(buf, dd)
}

// EncodeInt64Slice RLP encodes
func (enc *Encoder) EncodeInt64Slice(dd []int64) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextInt64Slice(buf, dd)
}

// EncodeUintSlice RLP encodes
func (enc *Encoder) EncodeUintSlice(dd []uint) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUintSlice(buf, dd)
}

// EncodeUint8Slice RLP encodes
func (enc *Encoder) EncodeUint8Slice(dd []uint8) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUint8Slice(buf, dd)
}

// EncodeUint16Slice RLP encodes
func (enc *Encoder) EncodeUint16Slice(dd []uint16) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUint16Slice(buf, dd)
}

// EncodeUint32Slice RLP encodes
func (enc *Encoder) EncodeUint32Slice(dd []uint32) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUint32Slice(buf, dd)
}

// EncodeUint64Slice RLP encodes
func (enc *Encoder) EncodeUint64Slice(dd []uint64) []byte {
	buf := enc.buf[:0]
	return enc.EncodeNextUint64Slice(buf, dd)
}

// EncodeBinaryMarshaler RLP encodes
func (enc *Encoder) EncodeBinaryMarshaler(d encoding.BinaryMarshaler) ([]byte, error) {
	buf := enc.buf[:]
	return enc.EncodeNextBinaryMarshaler(buf, d)
}

// EncodeTextMarshaler RLP encodes
func (enc *Encoder) EncodeTextMarshaler(d encoding.TextMarshaler) ([]byte, error) {
	buf := enc.buf[:]
	return enc.EncodeNextTextMarshaler(buf, d)
}

// EncodeNext RLP encodes
func (enc *Encoder) EncodeNext(buf []byte, d ...interface{}) (b []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	if len(d) == 0 {
		return nil, nil
	}
	if len(d) == 1 {
		return enc.encode(buf, d[0]), nil
	}
	return enc.encodeList(buf, d), nil
}

// EncodeNextList RLP encodes
func (enc *Encoder) EncodeNextList(buf []byte, idx uint) []byte {
	return enc.encodeListHelper(buf, idx)
}

// EncodeNextBytes RLP encodes
func (enc *Encoder) EncodeNextBytes(buf []byte, d []byte) []byte {
	return enc.encodeRLPString(buf, d)
}

// EncodeNextBytesSlice RLP encodes slice
func (enc *Encoder) EncodeNextBytesSlice(buf []byte, dd [][]byte) []byte {
	return enc.encodeRLPList(buf, dd)
}

// EncodeNextString RLP encodes
func (enc *Encoder) EncodeNextString(buf []byte, d string) []byte {
	b := enc.StringToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextBool RLP encodes
func (enc *Encoder) EncodeNextBool(buf []byte, d bool) []byte {
	b := enc.BoolToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextInt RLP encodes
func (enc *Encoder) EncodeNextInt(buf []byte, d int) []byte {
	b := enc.IntToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextInt8 RLP encodes
func (enc *Encoder) EncodeNextInt8(buf []byte, d int8) []byte {
	b := enc.Int8ToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextInt16 RLP encodes
func (enc *Encoder) EncodeNextInt16(buf []byte, d int16) []byte {
	b := enc.Int16ToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextInt32 RLP encodes
func (enc *Encoder) EncodeNextInt32(buf []byte, d int32) []byte {
	b := enc.Int32ToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextInt64 RLP encodes
func (enc *Encoder) EncodeNextInt64(buf []byte, d int64) []byte {
	b := enc.Int64ToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextUint RLP encodes
func (enc *Encoder) EncodeNextUint(buf []byte, d uint) []byte {
	b := enc.UintToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextUint8 RLP encodes
func (enc *Encoder) EncodeNextUint8(buf []byte, d uint8) []byte {
	b := enc.Uint8ToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextUint16 RLP encodes
func (enc *Encoder) EncodeNextUint16(buf []byte, d uint16) []byte {
	b := enc.Uint16ToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextUint32 RLP encodes
func (enc *Encoder) EncodeNextUint32(buf []byte, d uint32) []byte {
	b := enc.Uint32ToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextUint64 RLP encodes
func (enc *Encoder) EncodeNextUint64(buf []byte, d uint64) []byte {
	b := enc.Uint64ToBytes(d)
	return enc.encodeRLPString(buf, b)
}

// EncodeNextStringSlice RLP encodes
func (enc *Encoder) EncodeNextStringSlice(buf []byte, dd []string) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.StringToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextIntSlice RLP encodes
func (enc *Encoder) EncodeNextIntSlice(buf []byte, dd []int) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.IntToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextInt8Slice RLP encodes
func (enc *Encoder) EncodeNextInt8Slice(buf []byte, dd []int8) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.Int8ToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextInt16Slice RLP encodes
func (enc *Encoder) EncodeNextInt16Slice(buf []byte, dd []int16) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.Int16ToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextInt32Slice RLP encodes
func (enc *Encoder) EncodeNextInt32Slice(buf []byte, dd []int32) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.Int32ToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextInt64Slice RLP encodes
func (enc *Encoder) EncodeNextInt64Slice(buf []byte, dd []int64) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.Int64ToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextUintSlice RLP encodes
func (enc *Encoder) EncodeNextUintSlice(buf []byte, dd []uint) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.UintToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextUint8Slice RLP encodes
func (enc *Encoder) EncodeNextUint8Slice(buf []byte, dd []uint8) []byte {
	return enc.EncodeNextBytes(buf, dd)
}

// EncodeNextUint16Slice RLP encodes
func (enc *Encoder) EncodeNextUint16Slice(buf []byte, dd []uint16) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.Uint16ToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextUint32Slice RLP encodes
func (enc *Encoder) EncodeNextUint32Slice(buf []byte, dd []uint32) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.Uint32ToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextUint64Slice RLP encodes
func (enc *Encoder) EncodeNextUint64Slice(buf []byte, dd []uint64) []byte {
	bb := enc.lbuf[:0]
	for _, d := range dd {
		bb = append(bb, enc.Uint64ToBytes(d))
	}
	return enc.encodeRLPList(buf, bb)
}

// EncodeNextBinaryMarshaler RLP encodes
func (enc *Encoder) EncodeNextBinaryMarshaler(buf []byte, d encoding.BinaryMarshaler) ([]byte, error) {
	m, err := d.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return enc.EncodeNextBytes(buf, m), nil
}

// EncodeNextTextMarshaler RLP encodes
func (enc *Encoder) EncodeNextTextMarshaler(buf []byte, d encoding.TextMarshaler) ([]byte, error) {
	m, err := d.MarshalText()
	if err != nil {
		return nil, err
	}
	return enc.EncodeNextBytes(buf, m), nil
}

// Copy creates an independent copy of the buffer
func (enc *Encoder) Copy(buf []byte) []byte {
	c := make([]byte, len(buf))
	copy(c, buf)
	return c
}

// StringToBytes converts
func (enc *Encoder) StringToBytes(d string) []byte {
	return []byte(d)
}

// BoolToBytes converts
func (enc *Encoder) BoolToBytes(d bool) []byte {
	if d == false {
		return []byte{0x00}
	}
	return []byte{0x01}
}

// IntToBytes converts
func (enc *Encoder) IntToBytes(d int) []byte {
	if d == 0 {
		return nil
	}
	b := enc.buf10[:]
	n := binary.PutVarint(b, int64(d))
	return b[:n]
}

// Int8ToBytes converts
func (enc *Encoder) Int8ToBytes(d int8) []byte {
	if d == 0 {
		return nil
	}
	b := enc.buf2[:]
	n := binary.PutVarint(b, int64(d))
	return b[:n]
}

// Int16ToBytes converts
func (enc *Encoder) Int16ToBytes(d int16) []byte {
	if d == 0 {
		return nil
	}
	b := enc.buf3[:]
	n := binary.PutVarint(b, int64(d))
	return b[:n]
}

// Int32ToBytes converts
func (enc *Encoder) Int32ToBytes(d int32) []byte {
	if d == 0 {
		return nil
	}
	b := enc.buf5[:]
	n := binary.PutVarint(b, int64(d))
	return b[:n]
}

// Int64ToBytes converts
func (enc *Encoder) Int64ToBytes(d int64) []byte {
	if d == 0 {
		return nil
	}
	b := enc.buf10[:]
	n := binary.PutVarint(b, int64(d))
	return b[:n]
}

// UintToBytes converts
func (enc *Encoder) UintToBytes(d uint) []byte {
	if d == 0 {
		return nil
	}
	b := enc.buf8[:]
	binary.BigEndian.PutUint64(b, uint64(d))
	bl := binaryLen(d)
	return b[len(b)-int(bl):]
}

// Uint8ToBytes converts
func (enc *Encoder) Uint8ToBytes(d uint8) []byte {
	if d == 0 {
		return nil
	}
	return []byte{d}
}

// Uint16ToBytes converts
func (enc *Encoder) Uint16ToBytes(d uint16) []byte {
	if d == 0 {
		return nil
	}
	b := enc.buf2[:]
	binary.BigEndian.PutUint16(b, uint16(d))
	bl := binaryLen(uint(d))
	return b[len(b)-int(bl):]
}

// Uint32ToBytes converts
func (enc *Encoder) Uint32ToBytes(d uint32) []byte {
	if d == 0 {
		return nil
	}
	b := enc.buf4[:]
	binary.BigEndian.PutUint32(b, uint32(d))
	bl := binaryLen(uint(d))
	return b[len(b)-int(bl):]
}

// Uint64ToBytes converts
func (enc *Encoder) Uint64ToBytes(d uint64) []byte {
	if d == 0 {
		return nil
	}
	b := enc.buf8[:]
	binary.BigEndian.PutUint64(b, uint64(d))
	bl := binaryLen(uint(d))
	return b[len(b)-int(bl):]
}

func (enc *Encoder) encode(buf []byte, d interface{}) []byte {
	switch d.(type) {
	case []interface{}:
		return enc.encodeList(buf, d.([]interface{}))
	case []byte:
		return enc.EncodeNextBytes(buf, d.([]byte))
	case [][]byte:
		return enc.EncodeNextBytesSlice(buf, d.([][]byte))
	case string:
		return enc.EncodeNextString(buf, d.(string))
	case bool:
		return enc.EncodeNextBool(buf, d.(bool))
	case int:
		return enc.EncodeNextInt(buf, d.(int))
	case int8:
		return enc.EncodeNextInt8(buf, d.(int8))
	case int16:
		return enc.EncodeNextInt16(buf, d.(int16))
	case int32:
		return enc.EncodeNextInt32(buf, d.(int32))
	case int64:
		return enc.EncodeNextInt64(buf, d.(int64))
	case uint:
		return enc.EncodeNextUint(buf, d.(uint))
	case uint8:
		return enc.EncodeNextUint8(buf, d.(uint8))
	case uint16:
		return enc.EncodeNextUint16(buf, d.(uint16))
	case uint32:
		return enc.EncodeNextUint32(buf, d.(uint32))
	case uint64:
		return enc.EncodeNextUint64(buf, d.(uint64))
	case []string:
		return enc.EncodeNextStringSlice(buf, d.([]string))
	case []int:
		return enc.EncodeNextIntSlice(buf, d.([]int))
	case []int8:
		return enc.EncodeNextInt8Slice(buf, d.([]int8))
	case []int16:
		return enc.EncodeNextInt16Slice(buf, d.([]int16))
	case []int32:
		return enc.EncodeNextInt32Slice(buf, d.([]int32))
	case []int64:
		return enc.EncodeNextInt64Slice(buf, d.([]int64))
	case []uint:
		return enc.EncodeNextUintSlice(buf, d.([]uint))
	case []uint16:
		return enc.EncodeNextUint16Slice(buf, d.([]uint16))
	case []uint32:
		return enc.EncodeNextUint32Slice(buf, d.([]uint32))
	case []uint64:
		return enc.EncodeNextUint64Slice(buf, d.([]uint64))
	case encoding.BinaryMarshaler:
		m, err := enc.EncodeNextBinaryMarshaler(buf, d.(encoding.BinaryMarshaler))
		if err != nil {
			panic(err)
		}
		return m
	case encoding.TextMarshaler:
		m, err := enc.EncodeNextTextMarshaler(buf, d.(encoding.TextMarshaler))
		if err != nil {
			panic(err)
		}
		return m
	}
	panic(ErrInvalidType)
}

func (enc *Encoder) encodeList(buf []byte, dd []interface{}) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encode(buf, d)
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

func (enc *Encoder) encodeRLPList(buf []byte, dd [][]byte) []byte {
	idx := uint(len(buf))
	for _, d := range dd {
		buf = enc.encodeRLPString(buf, d)
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
		b = enc.Uint64ToBytes(uint64(l))
	}
	if bl+1 > 2 {
		b = enc.Uint32ToBytes(uint32(l))
	} else {
		b = enc.Uint16ToBytes(uint16(l))
	}
	buf = append(buf, byte(bl+off+55))
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