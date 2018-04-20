package rlp

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// These test cases come from https://github.com/ethereum/wiki/wiki/RLP
		{
			"The string 'dog'",
			args{"dog"},
			[]byte{0x83, 'd', 'o', 'g'},
			false,
		},
		{
			"The list [ 'cat', 'dog' ]",
			args{[]interface{}{"cat", "dog"}},
			[]byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'},
			false,
		},
		{
			"The empty string ('null')",
			args{""},
			[]byte{0x80},
			false,
		},
		{
			"The empty list",
			args{[]interface{}{}},
			[]byte{0xc0},
			false,
		},
		{
			"The integer 0",
			args{0},
			[]byte{0x80},
			false,
		},
		{
			"The encoded integer 0",
			args{"\x00"},
			[]byte{0x00},
			false,
		},
		{
			"The encoded integer 15",
			args{"\x0f"},
			[]byte{0x0f},
			false,
		},
		{
			"The encoded integer 1024",
			args{"\x04\x00"},
			[]byte{0x82, 0x04, 0x00},
			false,
		},
		{
			"The set theoretical representation of three",
			args{
				[]interface{}{
					[]interface{}{},
					[]interface{}{[]interface{}{}},
					[]interface{}{[]interface{}{}, []interface{}{[]interface{}{}}},
				},
			},
			[]byte{0xc7, 0xc0, 0xc1, 0xc0, 0xc3, 0xc0, 0xc1, 0xc0},
			false,
		},
		{
			"The string `Lorem ipsum dolor sit amet, consectetur adipisicing elit`",
			args{"Lorem ipsum dolor sit amet, consectetur adipisicing elit"},
			[]byte{0xb8, 0x38, 'L', 'o', 'r', 'e', 'm', ' ', 'i', 'p', 's', 'u', 'm', ' ', 'd', 'o', 'l', 'o', 'r', ' ', 's', 'i', 't', ' ', 'a', 'm', 'e', 't', ',', ' ', 'c', 'o', 'n', 's', 'e', 'c', 't', 'e', 't', 'u', 'r', ' ', 'a', 'd', 'i', 'p', 'i', 's', 'i', 'c', 'i', 'n', 'g', ' ', 'e', 'l', 'i', 't'},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkEncode_complex(b *testing.B) {
	t := struct {
		Name      string
		Age       uint
		Owed      map[string]uint
		Pointless []string
	}{"Bob", 150, map[string]uint{"Sue": 100, "Jim": 1000}, []string{"all", "your", "base"}}

	for i := 0; i < b.N; i++ {
		Encode(&t)
	}
}

func BenchmarkEncode_simple(b *testing.B) {
	t := "hello my name is inigo montoya"
	for i := 0; i < b.N; i++ {
		Encode(&t)
	}
}
