package rlp

import (
	"bytes"
	"reflect"
	"testing"
)

func BenchmarkEncoder_EncodeNode(b *testing.B) {
	t := make([][]byte, 0, 17)
	for i := 0; i < 17; i++ {
		t = append(t, bytes.Repeat([]byte{0xff}, 4))
	}
	for i := 0; i < b.N; i++ {
		Encode(t)
	}
}

func BenchmarkEncode_simple(b *testing.B) {
	t := "hello my name is inigo montoya"
	for i := 0; i < b.N; i++ {
		Encode(t)
	}
}

func BenchmarkEncode_simple_1(b *testing.B) {
	t := uint32(1024 * 1024)
	for i := 0; i < b.N; i++ {
		Encode(t)
	}
}

func BenchmarkEncode_simple_2(b *testing.B) {
	t := int32(-1024 * 1024)
	for i := 0; i < b.N; i++ {
		Encode(t)
	}
}

func BenchmarkEncode_med(b *testing.B) {
	t := []interface{}{"cat", "dog"}
	for i := 0; i < b.N; i++ {
		Encode(t)
	}
}

func BenchmarkEncode_complex(b *testing.B) {
	t := []interface{}{"cat", "dog", 77, []interface{}{"sheep", 666}}
	for i := 0; i < b.N; i++ {
		Encode(t)
	}
}

func BenchmarkEncode_complex_1(b *testing.B) {
	t := []string{"cat", "dog"}
	for i := 0; i < b.N; i++ {
		Encode(t)
	}
}

func BenchmarkEncode_complex_2(b *testing.B) {
	t := struct {
		Name      string
		Age       uint
		Pointless string
	}{
		"Bob",
		150,
		"all your base"}

	for i := 0; i < b.N; i++ {
		Encode(t)
	}
}

type CustomSerializer struct {
	Name      string
	Age       uint
	Pointless string
}

func (c *CustomSerializer) RLPSerialize() (interface{}, error) {
	return []interface{}{c.Name, c.Age, c.Pointless}, nil
}

func BenchmarkEncode_complex_custom(b *testing.B) {
	t := CustomSerializer{
		"Bob",
		150,
		"all your base"}

	for i := 0; i < b.N; i++ {
		Encode(&t)
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		e interface{}
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
		// Test cases below here are made up
		{
			"The integer 1024",
			args{uint(1024)},
			[]byte{0x82, 0x04, 0x00},
			false,
		},
		{
			"The integer 512",
			args{uint(512)},
			[]byte{0x82, 0x02, 0x00},
			false,
		},
		{
			"The string slice [ 'cat', 'dog' ]",
			args{[]string{"cat", "dog"}},
			[]byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'},
			false,
		},
		{
			"This crazy struct",
			args{struct {
				Name      string
				Age       uint
				Pointless []string
			}{
				"Bob",
				150,
				[]string{"all", "your", "base"}}},
			[]byte{0xd5, 0x83, 0x42, 0x6f, 0x62, 0x81, 0x96, 0xce, 0x83, 0x61, 0x6c, 0x6c, 0x84, 0x79, 0x6f, 0x75, 0x72, 0x84, 0x62, 0x61, 0x73, 0x65},
			false,
		},
		{
			"This custom struct",
			args{&CustomSerializer{
				"Bob",
				150,
				"all your base"}},
			[]byte{0xd4, 0x83, 0x42, 0x6f, 0x62, 0x81, 0x96, 0x8d, 0x61, 0x6c, 0x6c, 0x20, 0x79, 0x6f, 0x75, 0x72, 0x20, 0x62, 0x61, 0x73, 0x65},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = % x, want % x", got, tt.want)
			}
		})
	}
}
