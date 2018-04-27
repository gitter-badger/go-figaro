package rlp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

type rlpType int

const (
	_ rlpType = iota
	rlpString
	rlpList
)

type rlpItem struct {
	offset int
	length int
	typ    rlpType
}

var (
	errInvalidData = errors.New("rlp: invalid data must be encoded rlp string or rlp list")
	errInvalidDest = errors.New("rlp: destination for decoding is not valid for data")
)

// Decode takes RLP data and attempts to ummarshal it into dest
//
// `dest` must be a pointer
func Decode(dest interface{}, data []byte) (err error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		err = errInvalidType
	// 	}
	// }()
	// here we just decode the root, which must either be a list or a single item
	off, l, typ, err := next(data)
	if err != nil {
		return err
	}
	// the root should be the entire data
	if off+l != len(data) {
		return errInvalidData
	}
	switch typ {
	case rlpList:
		ii, err := decodeList(substr(data, off, l))
		if err != nil {
			return err
		}
		// once we have built up a list, we can deserialize it
		switch v := dest.(type) {
		// because it is untyped, we can't deserialize the content
		// of an []interface{} even though we could serialize it
		//
		// dest will be an (possibly arbitrarily nested) []interface{}
		// containing raw []byte; caller really should implement
		// rlp.Deserializer, but we support this just in case
		case *[]interface{}:
			*v = ii
			return nil
		case Deserializer:
			return v.RLPDeserialize(ii)
		}
		// fallback to reflection for slice/struct
		t := reflect.TypeOf(dest)
		if t.Kind() != reflect.Ptr {
			return errInvalidType
		}
		switch t.Elem().Kind() {
		case reflect.Slice:
			return deserializeSlice(dest, ii)
		case reflect.Struct:
			return deserializeStruct(dest, ii)
		}
		return errInvalidDest
	case rlpString:
		return deserializeString(dest, substr(data, off, l))
	default:
		return errWrongData
	}
}

func decodeList(data []byte) ([]interface{}, error) {
	// scan the data substring and build up a []interface{}, recursing as necesssary
	ii := make([]interface{}, 0)
	prevOff := 0
	prevL := 0

	for {
		off, l, typ, err := next(data[prevOff+prevL:])
		if err != nil {
			panic(err)
		}
		if typ == 0 {
			break
		}

		off += prevOff + prevL

		switch typ {
		case rlpList:
			r, err := decodeList(substr(data, off, l))
			if err != nil {
				return nil, err
			}
			ii = append(ii, r)
		case rlpString:
			ii = append(ii, substr(data, off, l))
		}
		prevOff = off
		prevL = l
	}
	return ii, nil
}

func next(data []byte) (off, l int, typ rlpType, err error) {
	length := len(data)
	if length == 0 {
		return 0, 0, 0, nil
	}
	prefix := int(data[0])
	if prefix <= 0x7f {
		return 0, 1, rlpString, nil
	}
	if prefix <= 0xb7 && length > prefix-0x80 {
		strLen := prefix - 0x80
		return 1, strLen, rlpString, nil
	}
	if prefix <= 0xbf && length > prefix-0xb7 && length > prefix-0xb7+toInt(substr(data, 1, prefix-0xb7)) {
		lenOfStrLen := prefix - 0xb7
		strLen := toInt(substr(data, 1, lenOfStrLen))
		return 1 + lenOfStrLen, strLen, rlpString, nil
	}
	if prefix <= 0xf7 && length > prefix-0xc0 {
		listLen := prefix - 0xc0
		return 1, listLen, rlpList, nil
	}
	if prefix <= 0xff && length > prefix-0xf7 && length > prefix-0xf7+toInt(substr(data, 1, prefix-0xf7)) {
		lenOfListLen := prefix - 0xf7
		listLen := toInt(substr(data, 1, lenOfListLen))
		return 1 + lenOfListLen, listLen, rlpList, nil
	}
	return 0, 0, 0, errInvalidData
}

func binaryToInt(b []byte) int {
	i, n := binary.Varint(b)
	if n <= 0 {
		panic(errInvalidData)
	}
	return int(i)
}

func binaryToUint(b []byte) uint {
	b = append(bytes.Repeat([]byte{0x00}, 8-len(b)), b...)
	return uint(binary.BigEndian.Uint64(b))
}

func substr(b []byte, o, l int) []byte {
	if o > len(b) {
		return b[len(b):]
	}
	if o+l > len(b) {
		return b[o:]
	}
	return b[o : o+l]
}

func toInt(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
