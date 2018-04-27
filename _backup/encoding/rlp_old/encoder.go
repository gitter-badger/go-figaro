package rlp_old

import (
	"errors"
)

var (
	errType = errors.New("rlp: data must be []byte or nested list of []byte")
)

// Package RLP implement the Recursive Length Prefix encoding protocol
// as outlined here: https://github.com/ethereum/wiki/wiki/RLP

// Encode RLP encodes arbitrary data into a []byte
func Encode(data interface{}) ([]byte, error) {
	s, err := serialize(data)
	if err != nil {
		return nil, err
	}
	return byteEncode(s)
}

// byteEncode implements RLP encoding of arbitrarily nested lists of []byte
func byteEncode(data interface{}) ([]byte, error) {
	if data == nil {
		return []byte{}, nil
	}
	// if data is an RLP string
	if d, ok := data.([]byte); ok {
		if len(d) == 1 && d[0] < 0x80 {
			return d, nil
		}
		el := encodeLength(len(d), 0x80)
		return append(el, d...), nil
	}
	// if data is an RLP list
	if ii, ok := data.([]interface{}); ok {
		d := []byte{}
		for _, i := range ii {
			ed, err := byteEncode(i)
			if err != nil {
				return nil, err
			}
			d = append(d, ed...)
		}
		el := encodeLength(len(d), 0xc0)
		return append(el, d...), nil
	}
	return nil, errType
}

// RLP length encoding schema
func encodeLength(l int, offset int) []byte {
	if l < 56 {
		return toBinary(l + offset)
	}
	bl := len(toBinary(l))
	return append(toBinary(bl+offset+55), toBinary(l)...)
}

// Converts an int to a []byte for convenience
func toBinary(i int) []byte {
	if i == 0 {
		return []byte{}
	}
	return append(toBinary(int(i/256)), byte(i%256))
}
