package rlp

import (
	"errors"
	"reflect"
	"testing"
)

func BenchmarkDecode_simple(b *testing.B) {
	var h string
	t := []byte{0x83, 'c', 'a', 't'}
	for i := 0; i < b.N; i++ {
		Decode(&h, t)
	}
}

func BenchmarkDecode_med(b *testing.B) {
	h := []string{}
	t := []byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'}
	for i := 0; i < b.N; i++ {
		Decode(&h, t)
	}
}

type CustomDeserializer struct {
	Name      string
	Age       uint
	Pointless string
}

func (c *CustomDeserializer) RLPDeserialize(data interface{}) error {
	ii, ok := data.([]interface{})
	if !ok {
		return errors.New("this isn't my data")
	}
	if len(ii) != 3 {
		return errors.New("this isn't my data")
	}
	deserializeString(&c.Name, ii[0].([]byte))
	deserializeString(&c.Age, ii[1].([]byte))
	deserializeString(&c.Pointless, ii[2].([]byte))
	return nil
}

func BenchmarkDecode_complex_custom(b *testing.B) {
	var h CustomDeserializer
	t := []byte{0xd4, 0x83, 0x42, 0x6f, 0x62, 0x81, 0x96, 0x8d, 0x61, 0x6c, 0x6c, 0x20, 0x79, 0x6f, 0x75, 0x72, 0x20, 0x62, 0x61, 0x73, 0x65}
	for i := 0; i < b.N; i++ {
		Decode(&h, t)
	}
}

func TestDecode(t *testing.T) {
	var h string
	var h2 []string
	var h3 CustomDeserializer
	type args struct {
		dest interface{}
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"A string",
			args{&h, []byte{0x83, 'c', 'a', 't'}},
			"cat",
			false,
		},
		{
			"The list [ 'cat', 'dog' ]",
			args{&h2, []byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'}},
			[]string{"cat", "dog"},
			false,
		},
		{
			"A custom deserializer",
			args{&h3, []byte{0xd4, 0x83, 0x42, 0x6f, 0x62, 0x81, 0x96, 0x8d, 0x61, 0x6c, 0x6c, 0x20, 0x79, 0x6f, 0x75, 0x72, 0x20, 0x62, 0x61, 0x73, 0x65}},
			CustomDeserializer{"Bob", 150, "all your base"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.dest, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(reflect.Indirect(reflect.ValueOf(tt.args.dest)).Interface(), tt.want) {
				t.Errorf("Decode() = %v, want %v", reflect.Indirect(reflect.ValueOf(tt.args.dest)), tt.want)
			}
		})
	}
}
