package rlp_old

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	"sort"
	"strconv"
)

// complex64 and complex128 NOT SUPPORTED

var (
	errSerializationNotSupported = errors.New("rlp: unsupported type serialization")
	errWrongType                 = errors.New("rlp: wrong data type supplied for serialization type")
	iSerializer                  = reflect.TypeOf((*Serializer)(nil)).Elem()
	iTextMarshaler               = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

// Serializer implementation Serializes itself into a []byte or
// an aribrarily nested lists of []byte
//		Ex:	[]interface{}{[]byte, []interface{}{[]byte, []byte}}
type Serializer interface {
	SerializeRLP() (interface{}, error)
}

// SerializeUint prepares for RLP encoding
func SerializeUint(d uint) (interface{}, error) {
	var b []byte
	if d > math.MaxUint32 {
		b = make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(d))
	} else if d > math.MaxUint16 {
		b = make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(d))
	} else if d > math.MaxUint8 {
		b = make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(d))
	} else {
		b = []byte{byte(d)}
	}
	return bytes.TrimLeft(b, "\x00"), nil
}

// SerializeUint8 prepares for RLP encoding
func SerializeUint8(d uint8) (interface{}, error) {
	return bytes.TrimLeft([]byte{d}, "\x00"), nil
}

// SerializeUint16 prepares for RLP encoding
func SerializeUint16(d uint16) (interface{}, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(d))
	return bytes.TrimLeft(b, "\x00"), nil
}

// SerializeUint32 prepares for RLP encoding
func SerializeUint32(d uint32) (interface{}, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(d))
	return bytes.TrimLeft(b, "\x00"), nil
}

// SerializeUint64 prepares for RLP encoding
func SerializeUint64(d uint64) (interface{}, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(d))
	return bytes.TrimLeft(b, "\x00"), nil
}

// SerializeInt prepares for RLP encoding
func SerializeInt(d int) (interface{}, error) {
	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(b, int64(d))
	return bytes.TrimLeft(b[:n], "\x00"), nil
}

// SerializeInt8 prepares for RLP encoding
func SerializeInt8(d int8) (interface{}, error) {
	b := make([]byte, binary.MaxVarintLen16)
	n := binary.PutVarint(b, int64(d))
	return bytes.TrimLeft(b[:n], "\x00"), nil
}

// SerializeInt16 prepares for RLP encoding
func SerializeInt16(d int16) (interface{}, error) {
	b := make([]byte, binary.MaxVarintLen16)
	n := binary.PutVarint(b, int64(d))
	return bytes.TrimLeft(b[:n], "\x00"), nil
}

// SerializeInt32 prepares for RLP encoding
func SerializeInt32(d int32) (interface{}, error) {
	b := make([]byte, binary.MaxVarintLen32)
	n := binary.PutVarint(b, int64(d))
	return bytes.TrimLeft(b[:n], "\x00"), nil
}

// SerializeInt64 prepares for RLP encoding
func SerializeInt64(d int64) (interface{}, error) {
	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(b, int64(d))
	return bytes.TrimLeft(b[:n], "\x00"), nil
}

// SerializeFloat32 prepares for RLP encoding
func SerializeFloat32(d float32) (interface{}, error) {
	return SerializeString(strconv.FormatFloat(float64(d), 'f', -1, 32))
}

// SerializeFloat64 prepares for RLP encoding
func SerializeFloat64(d float64) (interface{}, error) {
	return SerializeString(strconv.FormatFloat(d, 'f', -1, 64))
}

// SerializeBytes prepares for RLP encoding
func SerializeBytes(d []byte) (interface{}, error) {
	return d, nil
}

// SerializeString prepares for RLP encoding
func SerializeString(d string) (interface{}, error) {
	return []byte(d), nil
}

// SerializeBool prepares for RLP encoding
func SerializeBool(d bool) (interface{}, error) {
	if d {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// SerializeSlice prepares for RLP encoding
func SerializeSlice(d interface{}) (interface{}, error) {
	t := reflect.TypeOf(d)
	if t.Kind() != reflect.Slice {
		return nil, errWrongType
	}
	s := reflect.ValueOf(d)
	ii := make([]interface{}, s.Len(), s.Len())
	for i := 0; i < s.Len(); i++ {
		v, err := serialize(s.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		ii[i] = v
	}
	return ii, nil
}

// SerializeMap prepares for RLP encoding
func SerializeMap(d interface{}) (interface{}, error) {
	var sortFn func(int, int) bool
	t := reflect.TypeOf(d)
	if t.Kind() != reflect.Map {
		return nil, errWrongType
	}
	m := reflect.ValueOf(d)
	keys := m.MapKeys()
	switch t.Key().Kind() {
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		sortFn = func(i, j int) bool {
			return keys[i].Interface().(int) < keys[j].Interface().(int)
		}
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		sortFn = func(i, j int) bool {
			return keys[i].Interface().(float64) < keys[j].Interface().(float64)
		}
	case reflect.String:
		sortFn = func(i, j int) bool {
			return keys[i].Interface().(string) < keys[j].Interface().(string)
		}
	default:
		return nil, errSerializationNotSupported
	}
	sort.Slice(keys, sortFn)
	ii := make([]interface{}, m.Len(), m.Len())
	for i, kv := range keys {
		k, err := serialize(kv.Interface())
		if err != nil {
			return nil, err
		}
		v, err := serialize(m.MapIndex(kv).Interface())
		if err != nil {
			return nil, err
		}
		ii[i] = []interface{}{k, v}
	}
	return ii, nil
}

// SerializeStruct prepares for RLP encoding
func SerializeStruct(d interface{}) (interface{}, error) {
	t := reflect.TypeOf(d)
	if t.Kind() != reflect.Struct {
		return nil, errWrongType
	}
	v := reflect.ValueOf(d)
	ii := make([]interface{}, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		// PkgPath is empty for exported fields
		if ft.PkgPath != "" {
			continue
		}
		if tag, ok := ft.Tag.Lookup("rlp"); ok {
			if tag == "-" {
				continue
			}
		}
		k, err := SerializeString(ft.Name)
		if err != nil {
			return nil, err
		}
		v, err := serialize(v.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		ii = append(ii, []interface{}{k, v})
	}

	return ii, nil
}

// Serialize automatically serializes data deterministically for RLP encoding
func serialize(data interface{}) (interface{}, error) {
	// Check for custom serialization
	if d, ok := data.(Serializer); ok {
		return d.SerializeRLP()
	}
	if d, ok := data.(encoding.TextMarshaler); ok {
		return d.MarshalText()
	}
	if d, ok := data.([]byte); ok {
		return SerializeBytes(d)
	}
	if d, ok := data.(uint); ok {
		return SerializeUint(d)
	}
	if d, ok := data.(uint8); ok {
		return SerializeUint8(d)
	}
	if d, ok := data.(uint16); ok {
		return SerializeUint16(d)
	}
	if d, ok := data.(uint32); ok {
		return SerializeUint32(d)
	}
	if d, ok := data.(uint64); ok {
		return SerializeUint64(d)
	}
	if d, ok := data.(int); ok {
		return SerializeInt(d)
	}
	if d, ok := data.(int8); ok {
		return SerializeInt8(d)
	}
	if d, ok := data.(int16); ok {
		return SerializeInt16(d)
	}
	if d, ok := data.(int32); ok {
		return SerializeInt32(d)
	}
	if d, ok := data.(int64); ok {
		return SerializeInt64(d)
	}
	if d, ok := data.(float32); ok {
		return SerializeFloat32(d)
	}
	if d, ok := data.(float64); ok {
		return SerializeFloat64(d)
	}
	if d, ok := data.(string); ok {
		return SerializeString(d)
	}
	if d, ok := data.(bool); ok {
		return SerializeBool(d)
	}

	// Fallback to reflection for complex types
	t := reflect.TypeOf(data)
	switch t.Kind() {
	case reflect.Ptr:
		return serialize(reflect.ValueOf(data).Elem().Interface())
	case reflect.Slice:
		return SerializeSlice(data)
	case reflect.Map:
		return SerializeMap(data)
	case reflect.Struct:
		return SerializeStruct(data)
	}
	return nil, errSerializationNotSupported
}
