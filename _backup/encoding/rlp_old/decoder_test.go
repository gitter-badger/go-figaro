package rlp_old

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_byteDecode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"The string 'dog'",
			args{[]byte{0x83, 'd', 'o', 'g'}},
			[]byte("dog"),
			false,
		},
		{
			"The list [ 'cat', 'dog' ]",
			args{[]byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'}},
			[]interface{}{[]byte("cat"), []byte("dog")},
			false,
		},
		{
			"The list [ 'cat', 'dog' ] repeated a bunch",
			args{append([]byte{0xf9, 0x04, 0x00}, bytes.Repeat([]byte{0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'}, 128)...)},
			func() []interface{} {
				a := make([]interface{}, 128*2, 128*2)
				for i := range a {
					if i%2 == 0 {
						a[i] = []byte("cat")
					} else {
						a[i] = []byte("dog")
					}
				}
				return a
			}(),
			false,
		},
		{
			"The empty string ('null')",
			args{[]byte{0x80}},
			[]byte(""),
			false,
		},
		{
			"The empty list",
			args{[]byte{0xc0}},
			[]interface{}{},
			false,
		},
		{
			"The integer 0",
			args{[]byte{0x80}},
			[]byte{},
			false,
		},
		{
			"The encoded integer 0",
			args{[]byte{0x00}},
			[]byte{0x00},
			false,
		},
		{
			"The encoded integer 15",
			args{[]byte{0x0f}},
			[]byte("\x0f"),
			false,
		},
		{
			"The encoded integer 1024",
			args{[]byte{0x82, 0x04, 0x00}},
			[]byte("\x04\x00"),
			false,
		},
		{
			"The set theoretical representation of three",
			args{[]byte{0xc7, 0xc0, 0xc1, 0xc0, 0xc3, 0xc0, 0xc1, 0xc0}},
			[]interface{}{
				[]interface{}{},
				[]interface{}{[]interface{}{}},
				[]interface{}{[]interface{}{}, []interface{}{[]interface{}{}}},
			},
			false,
		},
		{
			"The string `Lorem ipsum dolor sit amet, consectetur adipisicing elit`",
			args{[]byte{0xb8, 0x38, 'L', 'o', 'r', 'e', 'm', ' ', 'i', 'p', 's', 'u', 'm', ' ', 'd', 'o', 'l', 'o', 'r', ' ', 's', 'i', 't', ' ', 'a', 'm', 'e', 't', ',', ' ', 'c', 'o', 'n', 's', 'e', 'c', 't', 'e', 't', 'u', 'r', ' ', 'a', 'd', 'i', 'p', 'i', 's', 'i', 'c', 'i', 'n', 'g', ' ', 'e', 'l', 'i', 't'}},
			[]byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := byteDecode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("byteDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("byteDecode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	h := []string{"cat", "dog"}
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
			"The list [ 'cat', 'dog' ]",
			args{&h, []byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'}},
			[]string{"cat", "dog"},
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

func BenchmarkEncode_med(b *testing.B) {
	t := make([]string, 0)
	d := []byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'}
	for i := 0; i < b.N; i++ {
		Decode(&t, d)
	}
}
