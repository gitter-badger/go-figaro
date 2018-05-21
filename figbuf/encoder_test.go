package figbuf_test

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"
	"testing"

	"github.com/figaro-tech/go-figaro/figbuf"
)

type SelfMarshaler struct {
	Name string
	Age  uint
}

type ProtoBeater struct {
	Name     string
	Age      uint
	Height   uint
	Weight   uint
	Alive    uint
	Desc     []byte
	Nickname string
	Num      uint
	Flt      uint
	Data     []byte
}

func ExampleSelfMarshaler() {
	s := &SelfMarshaler{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	b, _ := enc.Encode(s.Name, s.Age)
	fmt.Printf("% #x\n", b)
	// Output: 0xc5 0x83 0x42 0x6f 0x62 0x25
}

func ExampleSelfMarshaler_next() {
	s := &SelfMarshaler{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	b := enc.EncodeNextList(nil, func(buf []byte) []byte {
		next := enc.EncodeNextString(buf, s.Name)
		next = enc.EncodeNextUint(next, s.Age)
		return next
	})
	fmt.Printf("% #x\n", b)
	// Output: 0xc5 0x83 0x42 0x6f 0x62 0x25
}

func ExampleSelfMarshaler_protocompare() {
	enc := &figbuf.Encoder{}
	bench := &ProtoBeater{
		Name:     "Tester",
		Age:      20,
		Height:   58,
		Weight:   180,
		Alive:    1,
		Desc:     []byte("Lets benchmark some json and protobuf"),
		Nickname: "Another name",
		Num:      2314,
		Flt:      123451231,
		Data: []byte(`If you’ve ever heard of ProtoBuf you may be thinking that
		the results of this benchmarking experiment will be obvious, JSON < ProtoBuf.
		My interest was in how much they actually differ in practice.
		How do they compare on a couple of different metrics, specifically serialization
		and de-serialization speeds, and the memory footprint of encoding the data.
		I was also curious about how the different serialization methods would
		behave under small, medium, and large chunks of data.`),
	}
	b := enc.EncodeNextList(nil, func(buf []byte) []byte {
		next := enc.EncodeNextString(buf, bench.Name)
		next = enc.EncodeNextUint(next, bench.Age)
		next = enc.EncodeNextUint(next, bench.Height)
		next = enc.EncodeNextUint(next, bench.Weight)
		next = enc.EncodeNextUint(next, bench.Alive)
		next = enc.EncodeNextBytes(next, bench.Desc)
		next = enc.EncodeNextString(next, bench.Nickname)
		next = enc.EncodeNextUint(next, bench.Num)
		next = enc.EncodeNextUint(next, bench.Flt)
		next = enc.EncodeNextBytes(next, bench.Data)
		return next
	})
	fmt.Printf("% #x\n", b)
	// Output: 0xf9 0x02 0x37 0x86 0x54 0x65 0x73 0x74 0x65 0x72 0x14 0x3a 0x81 0xb4 0x01 0xa5 0x4c 0x65 0x74 0x73 0x20 0x62 0x65 0x6e 0x63 0x68 0x6d 0x61 0x72 0x6b 0x20 0x73 0x6f 0x6d 0x65 0x20 0x6a 0x73 0x6f 0x6e 0x20 0x61 0x6e 0x64 0x20 0x70 0x72 0x6f 0x74 0x6f 0x62 0x75 0x66 0x8c 0x41 0x6e 0x6f 0x74 0x68 0x65 0x72 0x20 0x6e 0x61 0x6d 0x65 0x82 0x09 0x0a 0x84 0x07 0x5b 0xb7 0x5f 0xb9 0x01 0xed 0x49 0x66 0x20 0x79 0x6f 0x75 0xe2 0x80 0x99 0x76 0x65 0x20 0x65 0x76 0x65 0x72 0x20 0x68 0x65 0x61 0x72 0x64 0x20 0x6f 0x66 0x20 0x50 0x72 0x6f 0x74 0x6f 0x42 0x75 0x66 0x20 0x79 0x6f 0x75 0x20 0x6d 0x61 0x79 0x20 0x62 0x65 0x20 0x74 0x68 0x69 0x6e 0x6b 0x69 0x6e 0x67 0x20 0x74 0x68 0x61 0x74 0x0a 0x09 0x09 0x74 0x68 0x65 0x20 0x72 0x65 0x73 0x75 0x6c 0x74 0x73 0x20 0x6f 0x66 0x20 0x74 0x68 0x69 0x73 0x20 0x62 0x65 0x6e 0x63 0x68 0x6d 0x61 0x72 0x6b 0x69 0x6e 0x67 0x20 0x65 0x78 0x70 0x65 0x72 0x69 0x6d 0x65 0x6e 0x74 0x20 0x77 0x69 0x6c 0x6c 0x20 0x62 0x65 0x20 0x6f 0x62 0x76 0x69 0x6f 0x75 0x73 0x2c 0x20 0x4a 0x53 0x4f 0x4e 0x20 0x3c 0x20 0x50 0x72 0x6f 0x74 0x6f 0x42 0x75 0x66 0x2e 0x0a 0x09 0x09 0x4d 0x79 0x20 0x69 0x6e 0x74 0x65 0x72 0x65 0x73 0x74 0x20 0x77 0x61 0x73 0x20 0x69 0x6e 0x20 0x68 0x6f 0x77 0x20 0x6d 0x75 0x63 0x68 0x20 0x74 0x68 0x65 0x79 0x20 0x61 0x63 0x74 0x75 0x61 0x6c 0x6c 0x79 0x20 0x64 0x69 0x66 0x66 0x65 0x72 0x20 0x69 0x6e 0x20 0x70 0x72 0x61 0x63 0x74 0x69 0x63 0x65 0x2e 0x0a 0x09 0x09 0x48 0x6f 0x77 0x20 0x64 0x6f 0x20 0x74 0x68 0x65 0x79 0x20 0x63 0x6f 0x6d 0x70 0x61 0x72 0x65 0x20 0x6f 0x6e 0x20 0x61 0x20 0x63 0x6f 0x75 0x70 0x6c 0x65 0x20 0x6f 0x66 0x20 0x64 0x69 0x66 0x66 0x65 0x72 0x65 0x6e 0x74 0x20 0x6d 0x65 0x74 0x72 0x69 0x63 0x73 0x2c 0x20 0x73 0x70 0x65 0x63 0x69 0x66 0x69 0x63 0x61 0x6c 0x6c 0x79 0x20 0x73 0x65 0x72 0x69 0x61 0x6c 0x69 0x7a 0x61 0x74 0x69 0x6f 0x6e 0x0a 0x09 0x09 0x61 0x6e 0x64 0x20 0x64 0x65 0x2d 0x73 0x65 0x72 0x69 0x61 0x6c 0x69 0x7a 0x61 0x74 0x69 0x6f 0x6e 0x20 0x73 0x70 0x65 0x65 0x64 0x73 0x2c 0x20 0x61 0x6e 0x64 0x20 0x74 0x68 0x65 0x20 0x6d 0x65 0x6d 0x6f 0x72 0x79 0x20 0x66 0x6f 0x6f 0x74 0x70 0x72 0x69 0x6e 0x74 0x20 0x6f 0x66 0x20 0x65 0x6e 0x63 0x6f 0x64 0x69 0x6e 0x67 0x20 0x74 0x68 0x65 0x20 0x64 0x61 0x74 0x61 0x2e 0x0a 0x09 0x09 0x49 0x20 0x77 0x61 0x73 0x20 0x61 0x6c 0x73 0x6f 0x20 0x63 0x75 0x72 0x69 0x6f 0x75 0x73 0x20 0x61 0x62 0x6f 0x75 0x74 0x20 0x68 0x6f 0x77 0x20 0x74 0x68 0x65 0x20 0x64 0x69 0x66 0x66 0x65 0x72 0x65 0x6e 0x74 0x20 0x73 0x65 0x72 0x69 0x61 0x6c 0x69 0x7a 0x61 0x74 0x69 0x6f 0x6e 0x20 0x6d 0x65 0x74 0x68 0x6f 0x64 0x73 0x20 0x77 0x6f 0x75 0x6c 0x64 0x0a 0x09 0x09 0x62 0x65 0x68 0x61 0x76 0x65 0x20 0x75 0x6e 0x64 0x65 0x72 0x20 0x73 0x6d 0x61 0x6c 0x6c 0x2c 0x20 0x6d 0x65 0x64 0x69 0x75 0x6d 0x2c 0x20 0x61 0x6e 0x64 0x20 0x6c 0x61 0x72 0x67 0x65 0x20 0x63 0x68 0x75 0x6e 0x6b 0x73 0x20 0x6f 0x66 0x20 0x64 0x61 0x74 0x61 0x2e
}

func ExampleEncoder_EncodeBytesSlice() {
	enc := &figbuf.Encoder{}
	t := [][]byte{{0xff, 0xfe}, {0xcd, 0x03}}
	b := enc.EncodeBytesSlice(t)
	fmt.Printf("% #x\n", b)
	// Output: 0xc6 0x82 0xff 0xfe 0x82 0xcd 0x03
}

func ExampleEncoder_EncodeBytesSlice_node() {
	enc := &figbuf.Encoder{}
	t := make([][]byte, 0, 17)
	for i := 0; i < 17; i++ {
		t = append(t, bytes.Repeat([]byte{0xff}, 32))
	}
	b := enc.EncodeBytesSlice(t)
	fmt.Printf("% #x\n", b)
	// Output: 0xf9 0x02 0x31 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xa0 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff 0xff
}

func BenchmarkSelfMarshaler(b *testing.B) {
	s := &SelfMarshaler{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	for i := 0; i < b.N; i++ {
		enc.Encode(s.Name, s.Age)
	}
}

func BenchmarkSelfMarshalerNext(b *testing.B) {
	s := &SelfMarshaler{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	var next []byte
	for i := 0; i < b.N; i++ {
		enc.EncodeNextList(nil, func(buf []byte) []byte {
			next = enc.EncodeNextString(buf, s.Name)
			next = enc.EncodeNextUint(next, s.Age)
			return next
		})
	}
}

func BenchmarkEncoder_EncodeString_protocompare(b *testing.B) {
	enc := &figbuf.Encoder{}
	bench := &ProtoBeater{
		Name:     "Tester",
		Age:      20,
		Height:   58,
		Weight:   180,
		Alive:    1,
		Desc:     []byte("Lets benchmark some json and protobuf"),
		Nickname: "Another name",
		Num:      2314,
		Flt:      123451231,
		Data: []byte(`If you’ve ever heard of ProtoBuf you may be thinking that
		the results of this benchmarking experiment will be obvious, JSON < ProtoBuf.
		My interest was in how much they actually differ in practice.
		How do they compare on a couple of different metrics, specifically serialization
		and de-serialization speeds, and the memory footprint of encoding the data.
		I was also curious about how the different serialization methods would
		behave under small, medium, and large chunks of data.`),
	}
	var next []byte
	for i := 0; i < b.N; i++ {
		enc.EncodeNextList(nil, func(buf []byte) []byte {
			next = enc.EncodeNextString(buf, bench.Name)
			next = enc.EncodeNextUint(next, bench.Age)
			next = enc.EncodeNextUint(next, bench.Height)
			next = enc.EncodeNextUint(next, bench.Weight)
			next = enc.EncodeNextUint(next, bench.Alive)
			next = enc.EncodeNextBytes(next, bench.Desc)
			next = enc.EncodeNextString(next, bench.Nickname)
			next = enc.EncodeNextUint(next, bench.Num)
			next = enc.EncodeNextUint(next, bench.Flt)
			next = enc.EncodeNextBytes(next, bench.Data)
			return next
		})
	}
}

func BenchmarkEncoder_Encode_bytes(b *testing.B) {
	enc := &figbuf.Encoder{}
	t := []byte{0xff, 0xee}
	for i := 0; i < b.N; i++ {
		enc.Encode(t)
	}
}

func BenchmarkEncoder_Encode_bytesSlice(b *testing.B) {
	enc := &figbuf.Encoder{}
	t := [][]byte{{0xff, 0xfe}, {0xcd, 0x03}}
	for i := 0; i < b.N; i++ {
		enc.Encode(t)
	}
}

func BenchmarkEncoder_EncodeBytes(b *testing.B) {
	enc := &figbuf.Encoder{}
	t := []byte{0xff, 0xee}
	for i := 0; i < b.N; i++ {
		enc.EncodeBytes(t)
	}
}

func BenchmarkEncoder_EncodeBytesSlice(b *testing.B) {
	enc := &figbuf.Encoder{}
	t := [][]byte{{0xff, 0xfe}, {0xcd, 0x03}}
	for i := 0; i < b.N; i++ {
		enc.EncodeBytesSlice(t)
	}
}

func BenchmarkEncoder_EncodeBytesSlice_node(b *testing.B) {
	enc := &figbuf.Encoder{}
	t := make([][]byte, 0, 17)
	for i := 0; i < 17; i++ {
		t = append(t, bytes.Repeat([]byte{0xff}, 32))
	}
	for i := 0; i < b.N; i++ {
		enc.EncodeBytesSlice(t)
	}
}

func TestEncoder_Encode(t *testing.T) {
	type args struct {
		d interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantB   []byte
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
			enc := &figbuf.Encoder{}
			gotB, err := enc.Encode(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("Encoder.Encode() = % #x, want % #x", gotB, tt.wantB)
			}
		})
	}
}

func TestEncoder_EncodeBytes(t *testing.T) {
	type args struct {
		d []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"A few bytes", args{[]byte{0xff, 0xfe}}, []byte{0x82, 0xff, 0xfe}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeBytes() = % #x, want % #x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeBytesSlice(t *testing.T) {
	type args struct {
		dd [][]byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"A few byte slices", args{[][]byte{{0xff, 0xfe}, {0xcd, 0x03}}}, []byte{0xc6, 0x82, 0xff, 0xfe, 0x82, 0xcd, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeBytesSlice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeBytesSlice() = % #x, want % #x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeString(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encoding a string", &figbuf.Encoder{}, args{"hello world"}, append([]byte{0x8b}, []byte("hello world")...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeString(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeString() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeBool(t *testing.T) {
	type args struct {
		d bool
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeBool(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeBool() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeInt(t *testing.T) {
	type args struct {
		d int
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeInt(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeInt8(t *testing.T) {
	type args struct {
		d int8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeInt8(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt8() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeInt16(t *testing.T) {
	type args struct {
		d int16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeInt16(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt16() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeInt32(t *testing.T) {
	type args struct {
		d int32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeInt32(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt32() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeInt64(t *testing.T) {
	type args struct {
		d int64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeInt64(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt64() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUint(t *testing.T) {
	type args struct {
		d uint
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUint(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUint8(t *testing.T) {
	type args struct {
		d uint8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUint8(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint8() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUint16(t *testing.T) {
	type args struct {
		d uint16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUint16(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint16() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUint32(t *testing.T) {
	type args struct {
		d uint32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUint32(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint32() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUint64(t *testing.T) {
	type args struct {
		d uint64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUint64(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint64() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeStringSlice(t *testing.T) {
	type args struct {
		dd []string
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeStringSlice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeStringSlice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeIntSlice(t *testing.T) {
	type args struct {
		dd []int
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeIntSlice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeIntSlice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeInt8Slice(t *testing.T) {
	type args struct {
		dd []int8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeInt8Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt8Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeInt16Slice(t *testing.T) {
	type args struct {
		dd []int16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeInt16Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt16Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeInt32Slice(t *testing.T) {
	type args struct {
		dd []int32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeInt32Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt32Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeInt64Slice(t *testing.T) {
	type args struct {
		dd []int64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeInt64Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt64Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUintSlice(t *testing.T) {
	type args struct {
		dd []uint
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUintSlice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUintSlice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUint8Slice(t *testing.T) {
	type args struct {
		dd []uint8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUint8Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint8Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUint16Slice(t *testing.T) {
	type args struct {
		dd []uint16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUint16Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint16Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUint32Slice(t *testing.T) {
	type args struct {
		dd []uint32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUint32Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint32Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeUint64Slice(t *testing.T) {
	type args struct {
		dd []uint64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encodce []uint64", &figbuf.Encoder{}, args{[]uint64{999999999999999, 66666666666, 3333333333333}}, []byte{0xd5, 0x87, 0x03, 0x8d, 0x7e, 0xa4, 0xc6, 0x7f, 0xff, 0x85, 0x0f, 0x85, 0xa4, 0x9a, 0xaa, 0x86, 0x03, 0x08, 0x1a, 0x26, 0x35, 0x55}},
		{"Encodce []uint64", &figbuf.Encoder{}, args{[]uint64{9, 6, 3}}, []byte{0xc3, 0x09, 0x06, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeUint64Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint64Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeBinaryMarshaler(t *testing.T) {
	type args struct {
		d encoding.BinaryMarshaler
	}
	tests := []struct {
		name    string
		enc     *figbuf.Encoder
		args    args
		wantBuf []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuf, err := tt.enc.EncodeBinaryMarshaler(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoder.EncodeBinaryMarshaler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBuf, tt.wantBuf) {
				t.Errorf("Encoder.EncodeBinaryMarshaler() = %#x, want %#x", gotBuf, tt.wantBuf)
			}
		})
	}
}

func TestEncoder_EncodeTextMarshaler(t *testing.T) {
	type args struct {
		d encoding.TextMarshaler
	}
	tests := []struct {
		name    string
		enc     *figbuf.Encoder
		args    args
		wantBuf []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuf, err := tt.enc.EncodeTextMarshaler(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoder.EncodeTextMarshaler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBuf, tt.wantBuf) {
				t.Errorf("Encoder.EncodeTextMarshaler() = %#x, want %#x", gotBuf, tt.wantBuf)
			}
		})
	}
}

func TestEncoder_EncodeNextList(t *testing.T) {
	type args struct {
		buf     []byte
		builder func([]byte) []byte
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextList(tt.args.buf, tt.args.builder); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextList() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextBytes(t *testing.T) {
	type args struct {
		buf []byte
		d   []byte
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextBytes(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextBytesSlice(t *testing.T) {
	type args struct {
		buf []byte
		dd  [][]byte
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextBytesSlice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextBytesSlice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextString(t *testing.T) {
	type args struct {
		buf []byte
		d   string
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next string", &figbuf.Encoder{}, args{nil, "hello world"}, []byte{0x8b, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextString(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextString() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextBool(t *testing.T) {
	type args struct {
		buf []byte
		d   bool
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next bool", &figbuf.Encoder{}, args{nil, true}, []byte{0x01}},
		{"Encode next bool", &figbuf.Encoder{}, args{nil, false}, []byte{0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextBool(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextBool() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextInt(t *testing.T) {
	type args struct {
		buf []byte
		d   int
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next int", &figbuf.Encoder{}, args{nil, int(-999999999)}, []byte{0x85, 0xfd, 0xa7, 0xd6, 0xb9, 0x07}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextInt(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextInt() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextInt8(t *testing.T) {
	type args struct {
		buf []byte
		d   int8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next int16", &figbuf.Encoder{}, args{nil, int8(-99)}, []byte{0x82, 0xc5, 0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextInt8(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextInt8() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextInt16(t *testing.T) {
	type args struct {
		buf []byte
		d   int16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next int16", &figbuf.Encoder{}, args{nil, int16(-9999)}, []byte{0x83, 0x9d, 0x9c, 0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextInt16(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextInt16() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextInt32(t *testing.T) {
	type args struct {
		buf []byte
		d   int32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next int32", &figbuf.Encoder{}, args{nil, int32(-999999999)}, []byte{0x85, 0xfd, 0xa7, 0xd6, 0xb9, 0x07}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextInt32(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextInt32() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextInt64(t *testing.T) {
	type args struct {
		buf []byte
		d   int64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next int64", &figbuf.Encoder{}, args{nil, int64(-9999999999999999)}, []byte{0x88, 0xfd, 0xff, 0x87, 0xfc, 0xcd, 0xbc, 0xc3, 0x23}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextInt64(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextInt64() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUint(t *testing.T) {
	type args struct {
		buf []byte
		d   uint
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next uint", &figbuf.Encoder{}, args{nil, uint(999999999)}, []byte{0x84, 0x3b, 0x9a, 0xc9, 0xff}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUint(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUint() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUint8(t *testing.T) {
	type args struct {
		buf []byte
		d   uint8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next uint8", &figbuf.Encoder{}, args{nil, uint8(99)}, []byte{0x63}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUint8(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUint8() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUint16(t *testing.T) {
	type args struct {
		buf []byte
		d   uint16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next uint16", &figbuf.Encoder{}, args{nil, uint16(9999)}, []byte{0x82, 0x27, 0x0f}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUint16(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUint16() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUint32(t *testing.T) {
	type args struct {
		buf []byte
		d   uint32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next uint64", &figbuf.Encoder{}, args{nil, uint32(999999999)}, []byte{0x84, 0x3b, 0x9a, 0xc9, 0xff}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUint32(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUint32() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUint64(t *testing.T) {
	type args struct {
		buf []byte
		d   uint64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next uint64", &figbuf.Encoder{}, args{nil, uint64(9999999999999999)}, []byte{0x87, 0x23, 0x86, 0xf2, 0x6f, 0xc0, 0xff, 0xff}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUint64(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUint64() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextStringSlice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []string
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []string", &figbuf.Encoder{}, args{nil, []string{"hello", "world", "zombie"}}, []byte{0xd3, 0x85, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x85, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x86, 0x7a, 0x6f, 0x6d, 0x62, 0x69, 0x65}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextStringSlice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextStringSlice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextBoolSlice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []bool
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []bool", &figbuf.Encoder{}, args{nil, []bool{true, true, false}}, []byte{0xc3, 0x01, 0x01, 0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextBoolSlice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextBoolSlice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextIntSlice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []int
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []int", &figbuf.Encoder{}, args{nil, []int{-999999999, 6666666, -3333333}}, []byte{0xd0, 0x85, 0xfd, 0xa7, 0xd6, 0xb9, 0x07, 0x84, 0xd4, 0xe6, 0xad, 0x06, 0x84, 0xa9, 0xf3, 0x96, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextIntSlice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextIntSlice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextInt8Slice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []int8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []int8", &figbuf.Encoder{}, args{nil, []int8{-99, 66, -33}}, []byte{0xc7, 0x82, 0xc5, 0x01, 0x82, 0x84, 0x01, 0x41}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextInt8Slice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextInt8Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextInt16Slice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []int16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []int16", &figbuf.Encoder{}, args{nil, []int16{-9999, 6666, -3333}}, []byte{0xca, 0x83, 0x9d, 0x9c, 0x01, 0x82, 0x94, 0x68, 0x82, 0x89, 0x34}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextInt16Slice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextInt16Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextInt32Slice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []int32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []int32", &figbuf.Encoder{}, args{nil, []int32{-999999999, 6666666, -3333333}}, []byte{0xd0, 0x85, 0xfd, 0xa7, 0xd6, 0xb9, 0x07, 0x84, 0xd4, 0xe6, 0xad, 0x06, 0x84, 0xa9, 0xf3, 0x96, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextInt32Slice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextInt32Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextInt64Slice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []int64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []int64", &figbuf.Encoder{}, args{nil, []int64{-999999999999999, 66666666666, -33333333333}}, []byte{0xd7, 0x88, 0xfd, 0xff, 0xb3, 0xcc, 0xd4, 0xdf, 0xc6, 0x03, 0x86, 0xd4, 0xea, 0xa4, 0xda, 0xf0, 0x03, 0x86, 0xa9, 0xb5, 0x92, 0xad, 0xf8, 0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextInt64Slice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextInt64Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUintSlice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []uint
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []uint", &figbuf.Encoder{}, args{nil, []uint{9999999, 666666, 333333}}, []byte{0xcc, 0x83, 0x98, 0x96, 0x7f, 0x83, 0x0a, 0x2c, 0x2a, 0x83, 0x05, 0x16, 0x15}},
		{"Encode next []uint", &figbuf.Encoder{}, args{nil, []uint{9, 6, 3}}, []byte{0xc3, 0x09, 0x06, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUintSlice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUintSlice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUint8Slice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []uint8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []uint8", &figbuf.Encoder{}, args{nil, []uint8{99, 66, 33}}, []byte{0x83, 0x63, 0x42, 0x21}},
		{"Encode next []uint8", &figbuf.Encoder{}, args{nil, []uint8{9, 6, 3}}, []byte{0x83, 0x09, 0x06, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUint8Slice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUint8Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUint16Slice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []uint16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []uint16", &figbuf.Encoder{}, args{nil, []uint16{9999, 666, 333}}, []byte{0xc9, 0x82, 0x27, 0x0f, 0x82, 0x02, 0x9a, 0x82, 0x01, 0x4d}},
		{"Encode next []uint16", &figbuf.Encoder{}, args{nil, []uint16{9, 6, 3}}, []byte{0xc3, 0x09, 0x06, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUint16Slice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUint16Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUint32Slice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []uint32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Encode next []uint32", &figbuf.Encoder{}, args{nil, []uint32{9999999, 666666, 333333}}, []byte{0xcc, 0x83, 0x98, 0x96, 0x7f, 0x83, 0x0a, 0x2c, 0x2a, 0x83, 0x05, 0x16, 0x15}},
		{"Encode next []uint32", &figbuf.Encoder{}, args{nil, []uint32{9, 6, 3}}, []byte{0xc3, 0x09, 0x06, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUint32Slice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUint32Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextUint64Slice(t *testing.T) {
	type args struct {
		buf []byte
		dd  []uint64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{"Encode next []uint64", &figbuf.Encoder{}, args{nil, []uint64{999999999999999, 66666666666, 3333333333333}}, []byte{0xd5, 0x87, 0x03, 0x8d, 0x7e, 0xa4, 0xc6, 0x7f, 0xff, 0x85, 0x0f, 0x85, 0xa4, 0x9a, 0xaa, 0x86, 0x03, 0x08, 0x1a, 0x26, 0x35, 0x55}},
		{"Encode next []uint64", &figbuf.Encoder{}, args{nil, []uint64{9, 6, 3}}, []byte{0xc3, 0x09, 0x06, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextUint64Slice(tt.args.buf, tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextUint64Slice() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextBinaryMarshaler(t *testing.T) {
	type args struct {
		buf []byte
		d   encoding.BinaryMarshaler
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextBinaryMarshaler(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextBinaryMarshaler() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_EncodeNextTextMarshaler(t *testing.T) {
	type args struct {
		buf []byte
		d   encoding.TextMarshaler
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodeNextTextMarshaler(tt.args.buf, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeNextTextMarshaler() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_StringToBytes(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert string to bytes", &figbuf.Encoder{}, args{"Hello World"}, []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.StringToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.StringToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_BoolToBytes(t *testing.T) {
	type args struct {
		d bool
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert true to bytes", &figbuf.Encoder{}, args{true}, []byte{0x01}},
		{"Convert false to bytes", &figbuf.Encoder{}, args{false}, []byte{0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.BoolToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.BoolToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_IntToBytes(t *testing.T) {
	type args struct {
		d int
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert int to bytes", &figbuf.Encoder{}, args{-81009841}, []byte{0xe1, 0xf2, 0xa0, 0x4d}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.IntToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.IntToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_Int8ToBytes(t *testing.T) {
	type args struct {
		d int8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert int16 to bytes", &figbuf.Encoder{}, args{-81}, []byte{0xa1, 0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.Int8ToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.Int8ToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_Int16ToBytes(t *testing.T) {
	type args struct {
		d int16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert int16 to bytes", &figbuf.Encoder{}, args{-8100}, []byte{0xc7, 0x7e}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.Int16ToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.Int16ToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_Int32ToBytes(t *testing.T) {
	type args struct {
		d int32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert int32 to bytes", &figbuf.Encoder{}, args{-81009841}, []byte{0xe1, 0xf2, 0xa0, 0x4d}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.Int32ToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.Int32ToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_Int64ToBytes(t *testing.T) {
	type args struct {
		d int64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert int64 to bytes", &figbuf.Encoder{}, args{-810098410000}, []byte{0x9f, 0x88, 0x86, 0xda, 0x93, 0x2f}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.Int64ToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.Int64ToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_UintToBytes(t *testing.T) {
	type args struct {
		d uint
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert uint to bytes", &figbuf.Encoder{}, args{810098410}, []byte{0x30, 0x49, 0x1e, 0xea}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.UintToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.UintToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_Uint8ToBytes(t *testing.T) {
	type args struct {
		d uint8
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert uint8 to bytes", &figbuf.Encoder{}, args{10}, []byte{0x0a}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.Uint8ToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.Uint8ToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_Uint16ToBytes(t *testing.T) {
	type args struct {
		d uint16
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert uint16 to bytes", &figbuf.Encoder{}, args{8410}, []byte{0x20, 0xda}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.Uint16ToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.Uint16ToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_Uint32ToBytes(t *testing.T) {
	type args struct {
		d uint32
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert uint32 to bytes", &figbuf.Encoder{}, args{810098410}, []byte{0x30, 0x49, 0x1e, 0xea}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.Uint32ToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.Uint32ToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestEncoder_Uint64ToBytes(t *testing.T) {
	type args struct {
		d uint64
	}
	tests := []struct {
		name string
		enc  *figbuf.Encoder
		args args
		want []byte
	}{
		{"Convert uint64 to bytes", &figbuf.Encoder{}, args{69309280810098410}, []byte{0xf6, 0x3c, 0x6c, 0x76, 0x52, 0xde, 0xea}},
		{"Convert uint64 to bytes less then max len", &figbuf.Encoder{}, args{66666666666}, []byte{0x0f, 0x85, 0xa4, 0x9a, 0xaa}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.Uint64ToBytes(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.Uint64ToBytes() = %#x, want %#x", got, tt.want)
			}
		})
	}
}
