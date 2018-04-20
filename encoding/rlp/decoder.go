package rlp

import (
	"errors"
)

const (
	str  = iota
	list = iota
)

var (
	errData = errors.New("rlp: data does not conform to RLP encoding")
	errNil  = errors.New("rlp: data is nil")
)

// Package RLP implement the Recursive Length Prefix encoding protocol
// as outlined here: https://github.com/ethereum/wiki/wiki/RLP

// Decode decodes an RLP encoding []byte into a destination interface
func Decode(dest interface{}, data []byte) error {
	d, _, err := byteDecode(data)
	if err != nil {
		return err
	}
	return deserialize(dest, d)
}

type olt struct {
	Offset int
	Length int
	Type   int
}

// byteDecode creates an arbitrarily nested lists of []byte from RLP encoding
func byteDecode(data []byte) (i interface{}, o *olt, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errData
		}
	}()
	if len(data) == 0 {
		return nil, nil, nil
	}
	dolt, err := decodeLength(data)
	if err != nil {
		return nil, nil, err
	}
	if dolt.Type == str {
		return substr(data, dolt.Offset, dolt.Length), dolt, nil
	}
	if dolt.Type == list {
		ii := []interface{}{}
		l := substr(data, dolt.Offset, dolt.Length)
		off := 0
		for {
			i, o, err := byteDecode(l[off:])
			if err != nil {
				return nil, nil, err
			}
			if i != nil {
				ii = append(ii, i)
				off += o.Offset + o.Length
				continue
			}
			break
		}
		return ii, dolt, nil
	}
	return nil, nil, errData
}

func decodeLength(data []byte) (o *olt, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errData
		}
	}()
	length := len(data)
	if length == 0 {
		return nil, errNil
	}
	prefix := int(data[0])
	if prefix <= 0x7f {
		return &olt{0, 1, str}, nil
	}
	if prefix <= 0xb7 && length > prefix-0x80 {
		strLen := prefix - 0x80
		return &olt{1, strLen, str}, nil
	}
	if prefix <= 0xbf && length > prefix-0xb7 && length > prefix-0xb7+toInt(substr(data, 1, prefix-0xb7)) {
		lenOfStrLen := prefix - 0xb7
		strLen := toInt(substr(data, 1, lenOfStrLen))
		return &olt{1 + lenOfStrLen, strLen, str}, nil
	}
	if prefix <= 0xf7 && length > prefix-0xc0 {
		listLen := prefix - 0xc0
		return &olt{1, listLen, list}, nil
	}
	if prefix <= 0xff && length > prefix-0xf7 && length > prefix-0xf7+toInt(substr(data, 1, prefix-0xf7)) {
		lenOfListLen := prefix - 0xf7
		listLen := toInt(substr(data, 1, lenOfListLen))
		return &olt{1 + lenOfListLen, listLen, list}, nil
	}
	return nil, errData
}

func toInt(b []byte) int {
	l := len(b)
	if l == 0 {
		return 0
	} else if l == 1 {
		return int(b[0])
	} else {
		return int(b[len(b)-1]) + toInt(b[:len(b)-1])*256
	}
}

func substr(b []byte, offset, length int) []byte {
	if offset > len(b)-1 {
		return []byte{}
	}
	if offset+length > len(b) {
		return b[offset:]
	}
	return b[offset : offset+length]
}
