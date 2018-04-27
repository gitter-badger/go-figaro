package rlp_old

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	errInvalidData   = errors.New("rlp: invalid data format for specified type")
	errInvalidDest   = errors.New("rlp: invalid dest must be pointer of supported type")
	iDeserializer    = reflect.TypeOf((*Deserializer)(nil)).Elem()
	iTextUnmarshaler = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

// Deserializer implementation Deserializes itself from a []byte or
// an aribrarily nested lists of []byte
//		Ex:	[]interface{}{[]byte, []interface{}{[]byte, []byte}}
type Deserializer interface {
	DeserializeRLP(interface{}) error
}

// DeserializeUint restores from RLP encoding
func DeserializeUint(dest *uint, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	if d[0] == '\x00' {
		return errInvalidData
	}
	d = append(bytes.Repeat([]byte{0x00}, 8-len(d)), d...)
	i := binary.BigEndian.Uint64(d)
	*dest = uint(i)
	return nil
}

// DeserializeUint8 restores from RLP encoding
func DeserializeUint8(dest *uint8, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	if len(d) > 1 {
		panic(errInvalidData)
	}
	if d[0] == '\x00' {
		return errInvalidData
	}
	*dest = uint8(d[0])
	return nil
}

// DeserializeUint16 restores from RLP encoding
func DeserializeUint16(dest *uint16, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	if d[0] == '\x00' {
		return errInvalidData
	}
	d = append(bytes.Repeat([]byte{0x00}, 2-len(d)), d...)
	i := binary.BigEndian.Uint16(d)
	*dest = uint16(i)
	return nil
}

// DeserializeUint32 restores from RLP encoding
func DeserializeUint32(dest *uint32, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	if d[0] == '\x00' {
		return errInvalidData
	}
	d = append(bytes.Repeat([]byte{0x00}, 4-len(d)), d...)
	i := binary.BigEndian.Uint32(d)
	*dest = uint32(i)
	return nil
}

// DeserializeUint64 restores from RLP encoding
func DeserializeUint64(dest *uint64, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	if d[0] == '\x00' {
		return errInvalidData
	}
	d = append(bytes.Repeat([]byte{0x00}, 8-len(d)), d...)
	*dest = binary.BigEndian.Uint64(d)
	return nil
}

// DeserializeInt restores from RLP encoding
func DeserializeInt(dest *int, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	b, n := binary.Varint(d)
	if n <= 0 {
		return errInvalidData
	}
	*dest = int(b)
	return nil
}

// DeserializeInt8 restores from RLP encoding
func DeserializeInt8(dest *int8, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	b, n := binary.Varint(d)
	if n <= 0 {
		return errInvalidData
	}
	*dest = int8(b)
	return nil
}

// DeserializeInt16 restores from RLP encoding
func DeserializeInt16(dest *int16, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	b, n := binary.Varint(d)
	if n <= 0 {
		return errInvalidData
	}
	*dest = int16(b)
	return nil
}

// DeserializeInt32 restores from RLP encoding
func DeserializeInt32(dest *int32, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	b, n := binary.Varint(d)
	if n <= 0 {
		return errInvalidData
	}
	*dest = int32(b)
	return nil
}

// DeserializeInt64 restores from RLP encoding
func DeserializeInt64(dest *int64, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	if len(d) == 0 {
		*dest = 0
		return nil
	}
	b, n := binary.Varint(d)
	if n <= 0 {
		return errInvalidData
	}
	*dest = b
	return nil
}

// DeserializeFloat32 restores from RLP encoding
func DeserializeFloat32(dest *float32, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	f, err := strconv.ParseFloat(string(d), 32)
	fmt.Println(f, err, err != nil)
	if err != nil {
		*dest = 0
		return err
	}
	*dest = float32(f)
	return nil
}

// DeserializeFloat64 restores from RLP encoding
func DeserializeFloat64(dest *float64, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	f, err := strconv.ParseFloat(string(d), 64)
	if err != nil {
		*dest = 0
		return err
	}
	*dest = float64(f)
	return nil
}

// DeserializeBytes restores from RLP encoding
func DeserializeBytes(dest *[]byte, d []byte) (err error) {
	*dest = d
	return nil
}

// DeserializeString restores from RLP encoding
func DeserializeString(dest *string, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	*dest = string(d)
	return nil
}

// DeserializeBool restores from RLP encoding
func DeserializeBool(dest *bool, d []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	switch d[0] {
	case 1:
		*dest = true
		return nil
	case 0:
		*dest = false
		return nil
	default:
		return errInvalidData
	}
}

// DeserializeSlice restores from RLP encoding
func DeserializeSlice(dest interface{}, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	pt := reflect.TypeOf(dest)
	if pt.Kind() != reflect.Ptr {
		return errInvalidDest
	}
	t := pt.Elem()
	v := reflect.ValueOf(dest).Elem()
	ii, ok := data.([]interface{})
	if !ok {
		return errInvalidData
	}
	v.Set(reflect.MakeSlice(t, len(ii), len(ii)))
	for i, elem := range ii {
		if t.Elem().Kind() == reflect.Ptr {
			p := reflect.New(t.Elem().Elem()).Interface()
			err := deserialize(p, elem)
			if err != nil {
				return err
			}
			v.Index(i).Set(reflect.ValueOf(p))
		} else {
			p := reflect.New(t.Elem()).Interface()
			err := deserialize(p, elem)
			if err != nil {
				return err
			}
			v.Index(i).Set(reflect.Indirect(reflect.ValueOf(p)))
		}

	}
	return nil
}

// DeserializeMap restores from RLP encoding
func DeserializeMap(dest interface{}, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	pt := reflect.TypeOf(dest)
	if pt.Kind() != reflect.Ptr {
		return errInvalidDest
	}
	t := reflect.TypeOf(dest).Elem()
	v := reflect.ValueOf(dest).Elem()
	ii, ok := data.([]interface{})
	if !ok {
		return errInvalidData
	}
	v.Set(reflect.MakeMap(t))
	for _, elem := range ii {
		var err error
		kv, ok := elem.([]interface{})
		if !ok || len(kv) != 2 {
			return errInvalidData
		}
		kd, vd := kv[0], kv[1]
		kp := reflect.New(t.Key()).Interface()
		err = deserialize(kp, kd)
		if err != nil {
			return err
		}
		if t.Elem().Kind() == reflect.Ptr {
			vp := reflect.New(t.Elem().Elem()).Interface()
			err = deserialize(vp, vd)
			if err != nil {
				return err
			}
			v.SetMapIndex(reflect.Indirect(reflect.ValueOf(kp)), reflect.ValueOf(vp))
		} else {
			vp := reflect.New(t.Elem()).Interface()
			err = deserialize(vp, vd)
			if err != nil {
				return err
			}
			v.SetMapIndex(reflect.Indirect(reflect.ValueOf(kp)), reflect.Indirect(reflect.ValueOf(vp)))
		}
	}
	return nil
}

// DeserializeStruct restores from RLP encoding
func DeserializeStruct(dest interface{}, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errInvalidData
		}
	}()
	pt := reflect.TypeOf(dest)
	if pt.Kind() != reflect.Ptr {
		return errInvalidDest
	}
	t := reflect.TypeOf(dest).Elem()
	v := reflect.ValueOf(dest).Elem()
	ii, ok := data.([]interface{})
	if !ok {
		v.Set(reflect.Zero(v.Type()))
		return errInvalidData
	}
	for _, elem := range ii {
		var err error
		fvs, ok := elem.([]interface{})
		if !ok || len(fvs) != 2 {
			v.Set(reflect.Zero(v.Type()))
			return errInvalidData
		}
		fd, vd := fvs[0], fvs[1]
		fp := reflect.New(reflect.TypeOf("")).Interface()
		err = deserialize(fp, fd)
		if err != nil {
			v.Set(reflect.Zero(v.Type()))
			return err
		}
		ft, ok := t.FieldByName(reflect.Indirect(reflect.ValueOf(fp)).Interface().(string))
		if !ok {
			v.Set(reflect.Zero(v.Type()))
			return errInvalidData
		}
		if tag, ok := ft.Tag.Lookup("rlp"); ok {
			if tag == "-" {
				v.Set(reflect.Zero(v.Type()))
				return errInvalidData
			}
		}
		fv := v.FieldByName(reflect.Indirect(reflect.ValueOf(fp)).Interface().(string))
		if !fv.IsValid() || !fv.CanSet() {
			v.Set(reflect.Zero(v.Type()))
			return errInvalidDest
		}
		if ft.Type.Kind() == reflect.Ptr {
			vp := reflect.New(ft.Type.Elem()).Interface()
			err = deserialize(vp, vd)
			if err != nil {
				v.Set(reflect.Zero(v.Type()))
				return err
			}
			fv.Set(reflect.ValueOf(vp))
		} else {
			vp := reflect.New(ft.Type).Interface()
			err = deserialize(vp, vd)
			if err != nil {
				v.Set(reflect.Zero(v.Type()))
				return err
			}
			fv.Set(reflect.Indirect(reflect.ValueOf(vp)))
		}
	}
	return nil
}

// Deserialize takes data returned from RLP decoding and converts it to
// the destination type, setting the pointer value
func deserialize(dest interface{}, data interface{}) (err error) {
	// Check for custom deserialization
	if d, ok := dest.(Deserializer); ok {
		return d.DeserializeRLP(data)
	}
	if d, ok := dest.(encoding.TextUnmarshaler); ok {
		if b, ok := data.([]byte); ok {
			return d.UnmarshalText(b)
		}
		return errInvalidData
	}
	if d, ok := dest.(*[]byte); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeBytes(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*uint); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeUint(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*uint8); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeUint8(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*uint16); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeUint16(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*uint32); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeUint32(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*uint64); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeUint64(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*int); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeInt(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*int8); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeInt8(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*int16); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeInt16(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*int32); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeInt32(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*int64); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeInt64(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*float32); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeFloat32(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*float64); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeFloat64(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*string); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeString(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}
	if d, ok := dest.(*bool); ok {
		if b, ok := data.([]byte); ok {
			err := DeserializeBool(d, b)
			if err != nil {
				return err
			}
			return nil
		}
		return errInvalidData
	}

	// Fallback to reflection for complex types
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr {
		return errInvalidDest
	}
	switch t.Elem().Kind() {
	case reflect.Slice:
		return DeserializeSlice(dest, data)
	case reflect.Map:
		return DeserializeMap(dest, data)
	case reflect.Struct:
		return DeserializeStruct(dest, data)
	default:
		return errInvalidDest
	}
}
