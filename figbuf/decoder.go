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
	// ErrInvalidDest raised when attemptimg to decode to a dest that is not well-known
	ErrInvalidDest = errors.New("figbuf: invalid dest for decoding must be well-known type")
)

// RlpType is either a str, a list, or nil
type RlpType uint

const (
	_ RlpType = iota
	// RlpStr is an RLP String ([]byte)
	RlpStr
	// RlpList is an RLP List ([][]byte)
	RlpList
)

// RlpItem represents an encoded item with an offset, length and type
type RlpItem struct {
	Offset uint
	Len    uint
	Typ    RlpType
}

// Decoder is an RLP decoder
type Decoder struct {
	// bs17s do not get reused across decodings,
	// but we allocate them in batches for efficiency
	bs17 [32][17][]byte
	bsat uint

	// Cache of tmp values
	length       uint
	prefix       uint
	lenOfStrLen  uint
	strLen       uint
	lenOfListLen uint
	listLen      uint
	zeros        [8]byte
	buf          [8]byte
}

// Decode decodes an item or list of items
func (dec *Decoder) Decode(b []byte, dest ...interface{}) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	if len(dest) == 0 {
		return
	}
	c := dec.Copy(b)
	var r []byte
	if len(dest) == 1 {
		r = dec.decode(c, dest[0])
	} else {
		r = dec.decode(c, dest)
	}
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeList RLP decodes a list, calling builder for encoding of each item
func (dec *Decoder) DecodeList(bb []byte, builder func([]byte) []byte) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	r = dec.DecodeNextList(c, builder)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeBytes decodes
func (dec *Decoder) DecodeBytes(b []byte) (d []byte, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextBytes(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeBytesSlice decodes
func (dec *Decoder) DecodeBytesSlice(bb []byte) (dd [][]byte, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextBytesSlice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeString decodes
func (dec *Decoder) DecodeString(b []byte) (d string, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = ""
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextString(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeBool decodes
func (dec *Decoder) DecodeBool(b []byte) (d bool, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = false
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextBool(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeInt decodes
func (dec *Decoder) DecodeInt(b []byte) (d int, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextInt(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeInt8 decodes
func (dec *Decoder) DecodeInt8(b []byte) (d int8, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextInt8(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeInt16 decodes
func (dec *Decoder) DecodeInt16(b []byte) (d int16, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextInt16(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeInt32 decodes
func (dec *Decoder) DecodeInt32(b []byte) (d int32, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextInt32(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeInt64 decodes
func (dec *Decoder) DecodeInt64(b []byte) (d int64, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextInt64(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUint decodes
func (dec *Decoder) DecodeUint(b []byte) (d uint, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextUint(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUint8 decodes
func (dec *Decoder) DecodeUint8(b []byte) (d uint8, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextUint8(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUint16 decodes
func (dec *Decoder) DecodeUint16(b []byte) (d uint16, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextUint16(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUint32 decodes
func (dec *Decoder) DecodeUint32(b []byte) (d uint32, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextUint32(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUint64 decodes
func (dec *Decoder) DecodeUint64(b []byte) (d uint64, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				d = 0
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	var r []byte
	d, r = dec.DecodeNextUint64(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeStringSlice decodes
func (dec *Decoder) DecodeStringSlice(bb []byte) (dd []string, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextStringSlice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeIntSlice decodes
func (dec *Decoder) DecodeIntSlice(bb []byte) (dd []int, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextIntSlice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeInt8Slice decodes
func (dec *Decoder) DecodeInt8Slice(bb []byte) (dd []int8, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextInt8Slice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeInt16Slice decodes
func (dec *Decoder) DecodeInt16Slice(bb []byte) (dd []int16, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextInt16Slice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeInt32Slice decodes
func (dec *Decoder) DecodeInt32Slice(bb []byte) (dd []int32, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextInt32Slice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeInt64Slice decodes
func (dec *Decoder) DecodeInt64Slice(bb []byte) (dd []int64, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextInt64Slice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUintSlice decodes
func (dec *Decoder) DecodeUintSlice(bb []byte) (dd []uint, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextUintSlice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUint8Slice decodes
func (dec *Decoder) DecodeUint8Slice(bb []byte) (dd []uint8, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextUint8Slice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUint16Slice decodes
func (dec *Decoder) DecodeUint16Slice(bb []byte) (dd []uint16, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextUint16Slice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUint32Slice decodes
func (dec *Decoder) DecodeUint32Slice(bb []byte) (dd []uint32, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextUint32Slice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeUint64Slice decodes
func (dec *Decoder) DecodeUint64Slice(bb []byte) (dd []uint64, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				dd = nil
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(bb)
	var r []byte
	dd, r = dec.DecodeNextUint64Slice(c)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeBinaryUnmarshaler decodes
func (dec *Decoder) DecodeBinaryUnmarshaler(b []byte, dest encoding.BinaryUnmarshaler) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	r := dec.DecodeNextBinaryUnmarshaler(c, dest)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeTextUnmarshaler decodes
func (dec *Decoder) DecodeTextUnmarshaler(b []byte, dest encoding.TextUnmarshaler) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if re, ok := rec.(error); ok {
				err = re
			} else {
				log.Panic(re)
			}
		}
	}()
	c := dec.Copy(b)
	r := dec.DecodeNextTextUnmarshaler(c, dest)
	if len(r) > 0 {
		err = ErrInvalidData
	}
	return
}

// DecodeNextList gets the next list
func (dec *Decoder) DecodeNextList(b []byte, builder func([]byte) []byte) (r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpList {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	lr := builder(dec.substr(b, item.Offset, item.Len))
	if len(lr) > 0 {
		panic(ErrInvalidData)
	}
	return
}

// DecodeNextBytes decodes
//
// Note that the slice returned may use the same
// backing array as `b` for performance
func (dec *Decoder) DecodeNextBytes(b []byte) (d []byte, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.substr(b, item.Offset, item.Len)
	return
}

// DecodeNextBytesSlice decodes
//
// Note that the slice returned may use the same
// backing array as `bb` for performance
func (dec *Decoder) DecodeNextBytesSlice(bb []byte) (dd [][]byte, r []byte) {
	// 17 is the size of a merkle node, which is
	// the most common use case for a bytes slice
	dd = dec.getByteSlice17()
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, buf)
	})
	return
}

// DecodeNextString decodes
func (dec *Decoder) DecodeNextString(b []byte) (d string, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToString(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextBool decodes
func (dec *Decoder) DecodeNextBool(b []byte) (d bool, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToBool(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextInt decodes
func (dec *Decoder) DecodeNextInt(b []byte) (d int, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToInt(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextInt8 decodes
func (dec *Decoder) DecodeNextInt8(b []byte) (d int8, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToInt8(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextInt16 decodes
func (dec *Decoder) DecodeNextInt16(b []byte) (d int16, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToInt16(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextInt32 decodes
func (dec *Decoder) DecodeNextInt32(b []byte) (d int32, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToInt32(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextInt64 decodes
func (dec *Decoder) DecodeNextInt64(b []byte) (d int64, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToInt64(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextUint decodes
func (dec *Decoder) DecodeNextUint(b []byte) (d uint, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToUint(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextUint8 decodes
func (dec *Decoder) DecodeNextUint8(b []byte) (d uint8, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToUint8(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextUint16 decodes
func (dec *Decoder) DecodeNextUint16(b []byte) (d uint16, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToUint16(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextUint32 decodes
func (dec *Decoder) DecodeNextUint32(b []byte) (d uint32, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToUint32(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextUint64 decodes
func (dec *Decoder) DecodeNextUint64(b []byte) (d uint64, r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	d = dec.BytesToUint64(dec.substr(b, item.Offset, item.Len))
	return
}

// DecodeNextStringSlice decodes
func (dec *Decoder) DecodeNextStringSlice(bb []byte) (dd []string, r []byte) {
	dd = make([]string, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToString(buf))
	})
	return
}

// DecodeNextBoolSlice decodes
func (dec *Decoder) DecodeNextBoolSlice(bb []byte) (dd []bool, r []byte) {
	dd = make([]bool, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToBool(buf))
	})
	return
}

// DecodeNextIntSlice decodes
func (dec *Decoder) DecodeNextIntSlice(bb []byte) (dd []int, r []byte) {
	dd = make([]int, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToInt(buf))
	})
	return
}

// DecodeNextInt8Slice decodes
func (dec *Decoder) DecodeNextInt8Slice(bb []byte) (dd []int8, r []byte) {
	dd = make([]int8, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToInt8(buf))
	})
	return
}

// DecodeNextInt16Slice decodes
func (dec *Decoder) DecodeNextInt16Slice(bb []byte) (dd []int16, r []byte) {
	dd = make([]int16, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToInt16(buf))
	})
	return
}

// DecodeNextInt32Slice decodes
func (dec *Decoder) DecodeNextInt32Slice(bb []byte) (dd []int32, r []byte) {
	dd = make([]int32, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToInt32(buf))
	})
	return
}

// DecodeNextInt64Slice decodes
func (dec *Decoder) DecodeNextInt64Slice(bb []byte) (dd []int64, r []byte) {
	dd = make([]int64, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToInt64(buf))
	})
	return
}

// DecodeNextUintSlice decodes
func (dec *Decoder) DecodeNextUintSlice(bb []byte) (dd []uint, r []byte) {
	dd = make([]uint, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToUint(buf))
	})
	return
}

// DecodeNextUint8Slice decodes
func (dec *Decoder) DecodeNextUint8Slice(bb []byte) (dd []uint8, r []byte) {
	return dec.DecodeNextBytes(bb)
}

// DecodeNextUint16Slice decodes
func (dec *Decoder) DecodeNextUint16Slice(bb []byte) (dd []uint16, r []byte) {
	dd = make([]uint16, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToUint16(buf))
	})
	return
}

// DecodeNextUint32Slice decodes
func (dec *Decoder) DecodeNextUint32Slice(bb []byte) (dd []uint32, r []byte) {
	dd = make([]uint32, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToUint32(buf))
	})
	return
}

// DecodeNextUint64Slice decodes
func (dec *Decoder) DecodeNextUint64Slice(bb []byte) (dd []uint64, r []byte) {
	dd = make([]uint64, 0, 8)
	r = dec.decodeSliceHelper(bb, func(buf []byte) {
		dd = append(dd, dec.BytesToUint64(buf))
	})
	return
}

// DecodeNextBinaryUnmarshaler decodes
func (dec *Decoder) DecodeNextBinaryUnmarshaler(b []byte, dest encoding.BinaryUnmarshaler) (r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	err := dest.UnmarshalBinary(dec.substr(b, item.Offset, item.Len))
	if err != nil {
		panic(err)
	}
	return
}

// DecodeNextTextUnmarshaler decodes
func (dec *Decoder) DecodeNextTextUnmarshaler(b []byte, dest encoding.TextUnmarshaler) (r []byte) {
	item := dec.nextItem(b)
	if item.Typ != RlpStr {
		panic(ErrInvalidData)
	}
	r = b[item.Offset+item.Len:]
	err := dest.UnmarshalText(dec.substr(b, item.Offset, item.Len))
	if err != nil {
		panic(err)
	}
	return
}

// Copy creates an independent copy of the buffer
func (dec *Decoder) Copy(buf []byte) []byte {
	c := make([]byte, len(buf))
	copy(c, buf)
	return c
}

// BytesToString converts
func (dec *Decoder) BytesToString(b []byte) string {
	return string(b)
}

// BytesToBool converts
func (dec *Decoder) BytesToBool(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	if len(b) == 1 {
		if b[0] == 0x00 {
			return false
		}
		return true
	}
	panic(ErrInvalidData)
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
	c := dec.buf[:0]
	c = append(c, dec.zeros[:8-len(b)]...)
	c = append(c, b...)
	return uint(binary.BigEndian.Uint64(c))
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
	c := dec.buf[:0]
	c = append(c, dec.zeros[:2-len(b)]...)
	c = append(c, b...)
	return binary.BigEndian.Uint16(c)
}

// BytesToUint32 converts
func (dec *Decoder) BytesToUint32(b []byte) uint32 {
	c := dec.buf[:0]
	c = append(c, dec.zeros[:4-len(b)]...)
	c = append(c, b...)
	return binary.BigEndian.Uint32(c)
}

// BytesToUint64 converts
func (dec *Decoder) BytesToUint64(b []byte) uint64 {
	c := dec.buf[:0]
	c = append(c, dec.zeros[:8-len(b)]...)
	c = append(c, b...)
	return binary.BigEndian.Uint64(c)
}

func (dec *Decoder) decode(b []byte, dest interface{}) (r []byte) {
	switch dest.(type) {
	case []interface{}:
		item := dec.nextItem(b)
		if item.Typ != RlpList {
			panic(ErrInvalidData)
		}
		r = b[item.Offset+item.Len:]
		lb := dec.substr(b, item.Offset, item.Len)
		for _, d := range dest.([]interface{}) {
			if len(lb) == 0 {
				panic(ErrInvalidData)
			}
			lb = dec.decode(lb, d)
		}
		return
	case *[]byte:
		*dest.(*[]byte), r = dec.DecodeNextBytes(b)
		return
	case *[][]byte:
		*dest.(*[][]byte), r = dec.DecodeNextBytesSlice(b)
		return
	case *string:
		*dest.(*string), r = dec.DecodeNextString(b)
		return
	case *bool:
		*dest.(*bool), r = dec.DecodeNextBool(b)
		return
	case *int:
		*dest.(*int), r = dec.DecodeNextInt(b)
		return
	case *int8:
		*dest.(*int8), r = dec.DecodeNextInt8(b)
		return
	case *int16:
		*dest.(*int16), r = dec.DecodeNextInt16(b)
		return
	case *int32:
		*dest.(*int32), r = dec.DecodeNextInt32(b)
		return
	case *int64:
		*dest.(*int64), r = dec.DecodeNextInt64(b)
		return
	case *uint:
		*dest.(*uint), r = dec.DecodeNextUint(b)
		return
	case *uint8:
		*dest.(*uint8), r = dec.DecodeNextUint8(b)
		return
	case *uint16:
		*dest.(*uint16), r = dec.DecodeNextUint16(b)
		return
	case *uint32:
		*dest.(*uint32), r = dec.DecodeNextUint32(b)
		return
	case *uint64:
		*dest.(*uint64), r = dec.DecodeNextUint64(b)
		return
	case *[]string:
		*dest.(*[]string), r = dec.DecodeNextStringSlice(b)
		return
	case *[]int:
		*dest.(*[]int), r = dec.DecodeNextIntSlice(b)
		return
	case *[]int8:
		*dest.(*[]int8), r = dec.DecodeNextInt8Slice(b)
		return
	case *[]int16:
		*dest.(*[]int16), r = dec.DecodeNextInt16Slice(b)
		return
	case *[]int32:
		*dest.(*[]int32), r = dec.DecodeNextInt32Slice(b)
		return
	case *[]int64:
		*dest.(*[]int64), r = dec.DecodeNextInt64Slice(b)
		return
	case *[]uint:
		*dest.(*[]uint), r = dec.DecodeNextUintSlice(b)
		return
	case *[]uint16:
		*dest.(*[]uint16), r = dec.DecodeNextUint16Slice(b)
		return
	case *[]uint32:
		*dest.(*[]uint32), r = dec.DecodeNextUint32Slice(b)
		return
	case *[]uint64:
		*dest.(*[]uint64), r = dec.DecodeNextUint64Slice(b)
		return
	case encoding.BinaryMarshaler:
		return dec.DecodeNextBinaryUnmarshaler(b, dest.(encoding.BinaryUnmarshaler))
	case encoding.TextMarshaler:
		return dec.DecodeNextTextUnmarshaler(b, dest.(encoding.TextUnmarshaler))
	}
	panic(ErrInvalidDest)
}

func (dec *Decoder) decodeListHelper(bb []byte, iter func(buf []byte)) (r []byte) {
	item := dec.nextItem(bb)
	if item.Typ != RlpList {
		panic(ErrInvalidData)
	}
	r = bb[item.Offset+item.Len:]
	b := dec.substr(bb, item.Offset, item.Len)
	for {
		item = dec.nextItem(b)
		if item.Typ == 0 {
			break
		}
		if item.Typ != RlpStr {
			panic(ErrInvalidData)
		}
		iter(dec.substr(b, item.Offset, item.Len))
		b = b[item.Offset+item.Len:]
	}
	return
}

func (dec *Decoder) decodeSliceHelper(bb []byte, iter func(buf []byte)) (r []byte) {
	item := dec.nextItem(bb)
	if item.Typ != RlpList {
		panic(ErrInvalidData)
	}
	r = bb[item.Offset+item.Len:]
	b := dec.substr(bb, item.Offset, item.Len)
	for {
		item = dec.nextItem(b)
		if item.Typ == 0 {
			break
		}
		if item.Typ != RlpStr {
			panic(ErrInvalidData)
		}
		iter(dec.substr(b, item.Offset, item.Len))
		b = b[item.Offset+item.Len:]
	}
	return
}

// nextItem gets the next RlpItem from the list
func (dec *Decoder) nextItem(b []byte) RlpItem {
	dec.length = uint(len(b))
	if dec.length == 0 {
		return RlpItem{0, 0, 0}
	}
	dec.prefix = uint(b[0])
	if dec.prefix <= 0x7f {
		return RlpItem{0, 1, RlpStr}
	}
	if dec.prefix <= 0xb7 && dec.length > dec.prefix-0x80 {
		dec.strLen = dec.prefix - 0x80
		return RlpItem{1, dec.strLen, RlpStr}
	}
	if dec.prefix <= 0xbf && dec.length > dec.prefix-0xb7 && dec.length > dec.prefix-0xb7+dec.BytesToUint(dec.substr(b, 1, dec.prefix-0xb7)) {
		dec.lenOfStrLen = dec.prefix - 0xb7
		dec.strLen = dec.BytesToUint(dec.substr(b, 1, dec.lenOfStrLen))
		return RlpItem{1 + dec.lenOfStrLen, dec.strLen, RlpStr}
	}
	if dec.prefix <= 0xf7 && dec.length > dec.prefix-0xc0 {
		dec.listLen = dec.prefix - 0xc0
		return RlpItem{1, dec.listLen, RlpList}
	}
	if dec.prefix <= 0xff && dec.length > dec.prefix-0xf7 && dec.length > dec.prefix-0xf7+dec.BytesToUint(dec.substr(b, 1, dec.prefix-0xf7)) {
		dec.lenOfListLen = dec.prefix - 0xf7
		dec.listLen = dec.BytesToUint(dec.substr(b, 1, dec.lenOfListLen))
		return RlpItem{1 + dec.lenOfListLen, dec.listLen, RlpList}
	}
	panic(ErrInvalidData)
}

func (dec *Decoder) substr(b []byte, offset, length uint) []byte {
	if offset > uint(len(b)) {
		return b[len(b):]
	}
	if offset+length > offset+uint(len(b)) {
		return b[offset:]
	}
	return b[offset : offset+length]
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
	bs := dec.bs17[dec.bsat]
	dec.bsat++
	return bs[:0]
}
