package figbuf_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/figaro-tech/figaro/pkg/figbuf"
)

type SelfMarshaler struct {
	Name string
	Age  uint
}

func (s *SelfMarshaler) MarshalDeterministicBinary(enc *figbuf.Encoder) ([]byte, error) {
	return enc.Encode(s.Name, s.Age)
}

func ExampleDeterministicBinaryMarshaler() {
	s := &SelfMarshaler{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	b, _ := enc.EncodeDeterministicBinaryMarshaler(s)
	fmt.Printf("% #x\n", b)
	// Output: 0xc5 0x83 0x42 0x6f 0x62 0x25
}

func BenchmarkDeterministicBinaryMarshaler(b *testing.B) {
	s := &SelfMarshaler{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	for i := 0; i < b.N; i++ {
		enc.EncodeDeterministicBinaryMarshaler(s)
	}
}

type SelfMarshalerCompose struct {
	Name string
	Age  uint
}

func (s *SelfMarshalerCompose) MarshalDeterministicBinary(enc *figbuf.Encoder) ([]byte, error) {
	// This is the equivalent of
	// List(s.Name, s.Age, List(s.Name, s.Age))
	/*
		return enc.EncodeList(
			enc.EncodeString(s.Name), enc.EncodeUint(s.Age), enc.EncodeList(
				enc.EncodeString(s.Name), enc.EncodeUint(s.Age)
			)
		), nil
	*/
	return enc.EncodeList(enc.Copy(enc.EncodeString(s.Name)), enc.Copy(enc.EncodeUint(s.Age))), nil
}

func ExampleDeterministicBinaryMarshaler_compose() {
	s := &SelfMarshalerCompose{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	b, _ := enc.EncodeDeterministicBinaryMarshaler(s)
	fmt.Printf("% #x\n", b)
	// Output: 0xc5 0x83 0x42 0x6f 0x62 0x25
}

func BenchmarkDeterministicBinaryMarshalerCompose(b *testing.B) {
	s := &SelfMarshalerCompose{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	for i := 0; i < b.N; i++ {
		enc.EncodeDeterministicBinaryMarshaler(s)
	}
}

type SelfMarshalerNext struct {
	Name string
	Age  uint
}

func (s *SelfMarshalerNext) MarshalDeterministicBinary(enc *figbuf.Encoder) ([]byte, error) {
	// This is the equivalent of
	// List(s.Name, s.Age, List(s.Name, s.Age), List(s.Name, s.Age), List(s.Name, List(s.Name, s.Age)))
	/*
		next := enc.EncodeNextString(nil, s.Name)
		next = enc.EncodeNextUint(next, s.Age)
		idx := uint(len(next))
		next = enc.EncodeNextString(next, s.Name)
		next = enc.EncodeNextUint(next, s.Age)
		next = enc.EncodeNextList(next, idx)
		idx = uint(len(next))
		next = enc.EncodeNextString(next, s.Name)
		next = enc.EncodeNextUint(next, s.Age)
		next = enc.EncodeNextList(next, idx)
		idx = uint(len(next))
		next = enc.EncodeNextString(next, s.Name)
		idxi := uint(len(next))
		next = enc.EncodeNextString(next, s.Name)
		next = enc.EncodeNextUint(next, s.Age)
		next = enc.EncodeNextList(next, idxi)
		next = enc.EncodeNextList(next, idx)
		next = enc.EncodeNextList(next, 0)
		return next, nil
	*/
	next := enc.EncodeNextString(nil, s.Name)
	next = enc.EncodeNextUint(next, s.Age)
	next = enc.EncodeNextList(next, 0)
	return next, nil
}

func ExampleDeterministicBinaryMarshaler_next() {
	s := &SelfMarshalerNext{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	b, _ := enc.EncodeDeterministicBinaryMarshaler(s)
	fmt.Printf("% #x\n", b)
	// Output: 0xc5 0x83 0x42 0x6f 0x62 0x25
}

func BenchmarkDeterministicBinaryMarshalerNext(b *testing.B) {
	s := &SelfMarshalerNext{Name: "Bob", Age: 37}
	enc := &figbuf.Encoder{}
	for i := 0; i < b.N; i++ {
		enc.EncodeDeterministicBinaryMarshaler(s)
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

func BenchmarkEncoder_EncodeBytes(b *testing.B) {
	enc := &figbuf.Encoder{}
	t := []byte{0xff, 0xee}
	for i := 0; i < b.N; i++ {
		enc.EncodeBytes(t)
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

func BenchmarkEncoder_EncodeBytesSlice(b *testing.B) {
	enc := &figbuf.Encoder{}
	t := [][]byte{{0xff, 0xfe}, {0xcd, 0x03}}
	for i := 0; i < b.N; i++ {
		enc.EncodeBytesSlice(t)
	}
}

func BenchmarkEncoder_EncodeNode(b *testing.B) {
	enc := &figbuf.Encoder{}
	t := make([][]byte, 0, 17)
	for i := 0; i < 17; i++ {
		t = append(t, bytes.Repeat([]byte{0xff}, 32))
	}
	for i := 0; i < b.N; i++ {
		enc.EncodeBytesSlice(t)
	}
}

func TestEncoder_EncodeList(t *testing.T) {
	type args struct {
		d [][]byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeList(tt.args.d...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeList() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeString(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeString() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeInt(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeInt8(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt8() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeInt16(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt16() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeInt32(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt32() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeInt64(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt64() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUint(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUint8(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint8() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUint16(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint16() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUint32(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint32() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUint64(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint64() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeStringSlice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeStringSlice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeIntSlice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeIntSlice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeInt8Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt8Slice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeInt16Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt16Slice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeInt32Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt32Slice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeInt64Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeInt64Slice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUintSlice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUintSlice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUint8Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint8Slice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUint16Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint16Slice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUint32Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint32Slice() = % #x, want % #x", got, tt.want)
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
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &figbuf.Encoder{}
			if got := enc.EncodeUint64Slice(tt.args.dd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.EncodeUint64Slice() = % #x, want % #x", got, tt.want)
			}
		})
	}
}
