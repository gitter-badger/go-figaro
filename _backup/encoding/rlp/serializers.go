package rlp

import (
	"reflect"
)

// Serializers provided here are for convenience
// for serializing complex types to well-known
// types automatically using reflection
//
// It will almost always be faster to do custom
// serialization that avoids reflection

func serializeSlice(s interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(s)
	ii := make([]interface{}, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		ii = append(ii, v.Index(i).Interface())
	}
	return ii, nil
}

func serializeStruct(s interface{}) ([]interface{}, error) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	ii := make([]interface{}, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		// This will be empty for things we can't set
		if ft.PkgPath != "" {
			continue
		}
		if tag, ok := ft.Tag.Lookup("rlp"); ok {
			if tag == "-" {
				continue
			}
		}
		ii = append(ii, v.Field(i).Interface())
	}
	return ii, nil
}
