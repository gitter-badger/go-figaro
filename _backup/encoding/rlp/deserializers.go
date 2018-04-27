package rlp

import (
	"encoding"
	"errors"
	"reflect"
)

// Deserializers provided here are for convenience
// for to reverse the convenience serializer functions
//
// It will almost always be faster to do custom
// de/serialization that avoids reflection

var (
	errWrongData = errors.New("rlp: wrong data type supplied for deserialization type")
	errWrongType = errors.New("rlp: wrong data type supplied for serialization type")
)

func DeserializeString(dest interface{}, data []byte) error {
	return deserializeString(dest, data)
}

func deserializeStruct(dest interface{}, s []interface{}) error {
	return nil
}

func deserializeSlice(dest interface{}, s []interface{}) error {
	// dest should be a pointer to a slice
	t := reflect.TypeOf(dest).Elem()
	v := reflect.ValueOf(dest).Elem()
	v.Set(reflect.MakeSlice(t, len(s), len(s)))
	for i, elem := range s {
		e := t.Elem()
		k := e.Kind()
		if k == reflect.Ptr {
			p := reflect.New(e.Elem())
			err := deserialze(k, p, elem)
			if err != nil {
				return err
			}
			v.Index(i).Set(reflect.ValueOf(p))
		} else {
			p := reflect.New(e).Interface()
			err := deserialze(k, p, elem)
			if err != nil {
				return err
			}
			v.Index(i).Set(reflect.Indirect(reflect.ValueOf(p)))
		}
	}
	return nil
}

func deserializeString(dest interface{}, data []byte) error {
	// decode RLP string into appropriate type or error
	switch v := dest.(type) {
	case *[]byte:
		*v = append(*v, data...)
	case *string:
		*v = string(data)
	case *uint:
		*v = binaryToUint(data)
	case *uint8:
		*v = uint8(binaryToUint(data))
	case *uint16:
		*v = uint16(binaryToUint(data))
	case *uint32:
		*v = uint32(binaryToUint(data))
	case *uint64:
		*v = uint64(binaryToUint(data))
	case *int:
		*v = binaryToInt(data)
	case *int8:
		*v = int8(binaryToInt(data))
	case *int16:
		*v = int16(binaryToInt(data))
	case *int32:
		*v = int32(binaryToInt(data))
	case *int64:
		*v = int64(binaryToInt(data))
	case encoding.BinaryUnmarshaler:
		v.UnmarshalBinary(data)
	case encoding.TextUnmarshaler:
		v.UnmarshalText(data)
	default:
		return errInvalidDest
	}
	return nil
}

func deserialze(kind reflect.Kind, dest, data interface{}) error {
	switch kind {
	case reflect.Slice:
		if ii, ok := data.([]interface{}); ok {
			return deserializeSlice(dest, ii)
		}
		return errInvalidData
	case reflect.Struct:
		if ii, ok := data.([]interface{}); ok {
			return deserializeStruct(dest, ii)
		}
		return errInvalidData
	default:
		if b, ok := data.([]byte); ok {
			return deserializeString(dest, b)
		}
		return errInvalidData
	}
}
