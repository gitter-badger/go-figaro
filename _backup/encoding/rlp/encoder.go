package rlp

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"reflect"
	"sync"
)

var (
	errInvalidType = errors.New("rlp: invalid type for encoding must be convertible to rlp string or rlp list")
	encPool        = sync.Pool{
		New: func() interface{} {
			return newEncoder()
		},
	}
)

// Encode intializes a new Encoder from an Encoder pool
// and RLP encodes the data, returning the result
//
// Unless it is desired to work with an existing Encoder,
// this is the easiest way to RLP encode arbitrary data
func Encode(e interface{}) ([]byte, error) {
	enc := getEncoder()
	defer putEncoder(enc)
	if err := enc.Encode(e); err != nil {
		return nil, err
	}
	return enc.bytes(), nil
}

type encoder struct {
	buf  []byte
	vbuf []byte
	lbuf []byte
}

func newEncoder() *encoder {
	return &encoder{
		buf:  make([]byte, 0, 8),
		vbuf: make([]byte, 8),
		lbuf: make([]byte, 8)}
}

func getEncoder() *encoder {
	return encPool.Get().(*encoder)
}

func putEncoder(enc *encoder) {
	enc.reset()
	encPool.Put(enc)
}

func (enc *encoder) reset() {
	enc.buf = enc.buf[:0]
}

func (enc *encoder) bytes() []byte {
	new := make([]byte, len(enc.buf), len(enc.buf))
	copy(new, enc.buf)
	return new
}

// Encode RLP encodes an interface{}
func (enc *encoder) Encode(e interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidType
		}
	}()
	switch e.(type) {
	case []byte:
		enc.encodeString(e.([]byte))
		return nil
	case []interface{}:
		enc.encodeList(e.([]interface{}))
		return nil
	case string:
		enc.encodeString([]byte(e.(string)))
		return nil
	case uint:
		enc.uintToBinary(e.(uint))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case uint8:
		enc.uintToBinary(uint(e.(uint8)))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case uint16:
		enc.uintToBinary(uint(e.(uint16)))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case uint32:
		enc.uintToBinary(uint(e.(uint32)))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case uint64:
		enc.uintToBinary(uint(e.(uint64)))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case int:
		n := enc.intToBinary(e.(int))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	case int8:
		n := enc.intToBinary(int(e.(int8)))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	case int16:
		n := enc.intToBinary(int(e.(int16)))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	case int32:
		n := enc.intToBinary(int(e.(int32)))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	case int64:
		n := enc.intToBinary(int(e.(int64)))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	case *string:
		enc.encodeString([]byte(*e.(*string)))
		return nil
	case *uint:
		enc.uintToBinary(*e.(*uint))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case *uint8:
		enc.uintToBinary(*e.(*uint))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case *uint16:
		enc.uintToBinary(*e.(*uint))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case *uint32:
		enc.uintToBinary(*e.(*uint))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case *uint64:
		enc.uintToBinary(*e.(*uint))
		enc.encodeString(bytes.TrimLeft(enc.vbuf, "\x00"))
		return nil
	case *int:
		n := enc.intToBinary(*e.(*int))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	case *int8:
		n := enc.intToBinary(*e.(*int))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	case *int16:
		n := enc.intToBinary(*e.(*int))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	case *int32:
		n := enc.intToBinary(*e.(*int))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	case *int64:
		n := enc.intToBinary(*e.(*int))
		enc.encodeString(bytes.TrimLeft(enc.vbuf[:n], "\x00"))
		return nil
	// rlp.Serializer will convert itself to a
	// well-known RLP type
	case Serializer:
		sz, err := e.(Serializer).RLPSerialize()
		if err != nil {
			return err
		}
		return enc.Encode(sz)
	// binary and text marshalers can be RLP
	// encoding as well-known type []byte
	case encoding.BinaryMarshaler:
		sz, err := e.(encoding.BinaryMarshaler).MarshalBinary()
		if err != nil {
			return err
		}
		enc.encodeString(sz)
		return nil
	case encoding.TextMarshaler:
		sz, err := e.(encoding.TextMarshaler).MarshalText()
		if err != nil {
			return err
		}
		enc.encodeString(sz)
		return nil
	}

	// fallback to reflection for pointers, slice, map and struct
	t := reflect.TypeOf(e)
	switch t.Kind() {
	case reflect.Ptr:
		return enc.Encode(reflect.Indirect(reflect.ValueOf(e)).Interface())
	case reflect.Slice:
		sz, err := serializeSlice(e)
		if err != nil {
			return err
		}
		return enc.Encode(sz)
	case reflect.Struct:
		sz, err := serializeStruct(e)
		if err != nil {
			return err
		}
		return enc.Encode(sz)
	}
	return errInvalidType
}

func (enc *encoder) encodeString(e []byte) {
	if len(e) == 1 && e[0] < 0x80 {
		enc.buf = append(enc.buf, e...)
		return
	}
	enc.encodeLength(len(e), 0x80)
	enc.buf = append(enc.buf, e...)
}

func (enc *encoder) encodeList(e []interface{}) {
	// mark our place
	bIndex := len(enc.buf)
	for _, item := range e {
		err := enc.Encode(item)
		if err != nil {
			panic(err)
		}
	}
	// we shift things a bit so we can insert the length at the proper place
	// TODO: explore faster ways to do solve this problem
	l := len(enc.buf) - bIndex
	var pad int
	if l < 56 {
		pad = 1
	} else {
		pad = 1 + binaryLen(l)
	}
	padding := make([]byte, pad, pad)
	enc.buf = append(enc.buf, padding...)
	copy(enc.buf[bIndex+pad:], enc.buf[bIndex:])
	enc.buf = enc.buf[:bIndex]
	enc.encodeLength(l, 0xc0)
	enc.buf = enc.buf[:len(enc.buf)+l]
}

func (enc *encoder) encodeLength(length int, offset int) {
	if length < 56 {
		binary.BigEndian.PutUint64(enc.lbuf, uint64(length+offset))
		enc.buf = append(enc.buf, bytes.TrimLeft(enc.lbuf, "\x00")...)
	} else {
		bl := binaryLen(length)
		binary.BigEndian.PutUint64(enc.lbuf, uint64(bl+offset+55))
		enc.buf = append(enc.buf, bytes.TrimLeft(enc.lbuf, "\x00")...)
		binary.BigEndian.PutUint64(enc.lbuf, uint64(length))
		enc.buf = append(enc.buf, bytes.TrimLeft(enc.lbuf, "\x00")...)
	}
}

func (enc *encoder) intToBinary(i int) int {
	return binary.PutVarint(enc.vbuf, int64(i))
}

func (enc *encoder) uintToBinary(u uint) {
	binary.BigEndian.PutUint64(enc.vbuf, uint64(u))
}

func binaryLen(u int) int {
	if u%256 == 0 {
		return u / 256
	}
	return u/256 + 1
}
