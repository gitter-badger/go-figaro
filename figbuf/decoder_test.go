package figbuf_test

import (
	"encoding"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/figaro-tech/go-figaro/figbuf"
)

type SelfUnmarshaler struct {
	Name string
	Age  uint
}

type SelfUnmarshalerNext struct {
	Name string
	Age  uint
}

func ExampleSelfUnmarshaler() {
	s := &SelfUnmarshaler{}
	dec := &figbuf.Decoder{}
	b := []byte{0xc5, 0x83, 0x42, 0x6f, 0x62, 0x25}
	_, err := dec.Decode(b, &s.Name, &s.Age)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", *s)
	// Output: {Name:Bob Age:37}
}

func ExampleSelfUnmarshaler_next() {
	s := &SelfUnmarshalerNext{}
	dec := &figbuf.Decoder{}
	bs := []byte{0xc5, 0x83, 0x42, 0x6f, 0x62, 0x25}
	var r []byte
	r = dec.DecodeNextList(bs, func(b []byte) {
		s.Name, r = dec.DecodeNextString(b)
		s.Age, r = dec.DecodeNextUint(r)
	})
	if len(r) > 0 {
		log.Fatal("invalid encoding")
	}
	fmt.Printf("%+v\n", *s)
	// Output: {Name:Bob Age:37}
}

func ExampleDecoder_DecodeBytesSlice_node() {
	dec := &figbuf.Decoder{}
	t := []byte{0xf9, 0x02, 0x31, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	d, _, err := dec.DecodeBytesSlice(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("% #v\n", d)
	// Output: [][]uint8{[]uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}, []uint8{ 0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff,  0xff}}
}

func BenchmarkSelfUnmarshaler(b *testing.B) {
	s := &SelfUnmarshaler{}
	dec := &figbuf.Decoder{}
	bs := []byte{0xc5, 0x83, 0x42, 0x6f, 0x62, 0x25}
	var err error
	for i := 0; i < b.N; i++ {
		_, err = dec.Decode(bs, &s.Name, &s.Age)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkSelfUnmarshalerNext(b *testing.B) {
	s := &SelfUnmarshalerNext{}
	dec := &figbuf.Decoder{}
	bs := []byte{0xc5, 0x83, 0x42, 0x6f, 0x62, 0x25}
	for i := 0; i < b.N; i++ {
		dec.DecodeNextList(bs, func(b []byte) {
			s.Name, b = dec.DecodeNextString(b)
			s.Age, b = dec.DecodeNextUint(b)
		})
	}
}

func BenchmarkDecoder_BytesToUint64(b *testing.B) {
	dec := &figbuf.Decoder{}
	t := []byte{0x0f, 0x85, 0xa4, 0x9a, 0xaa}
	var n uint64
	for i := 0; i < b.N; i++ {
		n = dec.BytesToUint64(t)
		if n != 66666666666 || len(t) != 5 {
			log.Fatal("not equal")
		}
	}
}

func BenchmarkDecoder_DecodeBytes(b *testing.B) {
	dec := &figbuf.Decoder{}
	t := []byte{0x82, 0xff, 0xfe}
	for i := 0; i < b.N; i++ {
		_, _, err := dec.DecodeBytes(t)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkDecoder_DecodeBytesSlice(b *testing.B) {
	dec := &figbuf.Decoder{}
	t := []byte{0xc6, 0x82, 0xff, 0xfe, 0x82, 0xcd, 0x03}
	for i := 0; i < b.N; i++ {
		_, _, err := dec.DecodeBytesSlice(t)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkDecoder_DecodeBytesSlice_node(b *testing.B) {
	dec := &figbuf.Decoder{}
	t := []byte{0xf9, 0x02, 0x31, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	for i := 0; i < b.N; i++ {
		_, _, err := dec.DecodeBytesSlice(t)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestDecoder_DecodeBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"A few bytes", args{[]byte{0x82, 0xff, 0xfe}}, []byte{0xff, 0xfe}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &figbuf.Decoder{}
			got, _, err := dec.DecodeBytes(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.DecodeBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_DecodeBytesSlice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		args    args
		want    [][]byte
		wantErr bool
	}{
		{"A few byte slices", args{[]byte{0xc6, 0x82, 0xff, 0xfe, 0x82, 0xcd, 0x03}}, [][]byte{{0xff, 0xfe}, {0xcd, 0x03}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &figbuf.Decoder{}
			got, _, err := dec.DecodeBytesSlice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeBytesSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.DecodeBytesSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	type args struct {
		b    []byte
		dest []interface{}
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := tt.dec.Decode(tt.args.b, tt.args.dest...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.Decode() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   string
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeString(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeString() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeString() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeBool(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   bool
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeBool(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeBool() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeBool() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeInt(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   int
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeInt(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeInt() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeInt() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeInt8(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   int8
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeInt8(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeInt8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeInt8() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeInt8() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeInt16(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   int16
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeInt16(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeInt16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeInt16() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeInt16() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeInt32(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   int32
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeInt32(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeInt32() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeInt32() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeInt64(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   int64
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeInt64(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeInt64() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeInt64() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUint(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   uint
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeUint(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeUint() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUint() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUint8(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   uint8
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeUint8(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUint8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeUint8() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUint8() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUint16(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   uint16
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeUint16(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeUint16() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUint16() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUint32(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   uint32
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeUint32(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeUint32() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUint32() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUint64(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantD   uint64
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR, err := tt.dec.DecodeUint64(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeUint64() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUint64() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeStringSlice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []string
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeStringSlice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeStringSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeStringSlice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeStringSlice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeIntSlice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []int
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeIntSlice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeIntSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeIntSlice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeIntSlice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeInt8Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []int8
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeInt8Slice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeInt8Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeInt8Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeInt8Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeInt16Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []int16
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeInt16Slice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeInt16Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeInt16Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeInt16Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeInt32Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []int32
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeInt32Slice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeInt32Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeInt32Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeInt32Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeInt64Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []int64
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeInt64Slice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeInt64Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeInt64Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeInt64Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUintSlice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []uint
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeUintSlice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUintSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeUintSlice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUintSlice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUint8Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []uint8
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeUint8Slice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUint8Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeUint8Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUint8Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUint16Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []uint16
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeUint16Slice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUint16Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeUint16Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUint16Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUint32Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []uint32
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeUint32Slice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUint32Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeUint32Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUint32Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeUint64Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantDd  []uint64
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR, err := tt.dec.DecodeUint64Slice(tt.args.bb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeUint64Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeUint64Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeUint64Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeBinaryUnmarshaler(t *testing.T) {
	type args struct {
		b    []byte
		dest encoding.BinaryUnmarshaler
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := tt.dec.DecodeBinaryUnmarshaler(tt.args.b, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeBinaryUnmarshaler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeBinaryUnmarshaler() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeTextUnmarshaler(t *testing.T) {
	type args struct {
		b    []byte
		dest encoding.TextUnmarshaler
	}
	tests := []struct {
		name    string
		dec     *figbuf.Decoder
		args    args
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := tt.dec.DecodeTextUnmarshaler(tt.args.b, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.DecodeTextUnmarshaler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeTextUnmarshaler() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextList(t *testing.T) {
	type args struct {
		b       []byte
		builder func([]byte)
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.DecodeNextList(tt.args.b, tt.args.builder); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.DecodeNextList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_DecodeNextBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD []byte
		wantR []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextBytes(tt.args.b)
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("Decoder.DecodeNextBytes() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextBytes() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextBytesSlice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd [][]byte
		wantR  []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextBytesSlice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextBytesSlice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextBytesSlice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD string
		wantR []byte
	}{
		{"Decode next string", &figbuf.Decoder{}, args{[]byte{0x8b, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64}}, "hello world", []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextString(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextString() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextString() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextBool(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD bool
		wantR []byte
	}{
		{"Decode next bool", &figbuf.Decoder{}, args{[]byte{0x01}}, true, []byte{}},
		{"Decode next bool", &figbuf.Decoder{}, args{[]byte{0x00}}, false, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextBool(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextBool() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextBool() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextInt(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD int
		wantR []byte
	}{
		{"Decode next int", &figbuf.Decoder{}, args{[]byte{0x85, 0xfd, 0xa7, 0xd6, 0xb9, 0x07}}, int(-999999999), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextInt(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextInt() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextInt() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextInt8(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD int8
		wantR []byte
	}{
		{"Decode next int16", &figbuf.Decoder{}, args{[]byte{0x82, 0xc5, 0x01}}, int8(-99), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextInt8(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextInt8() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextInt8() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextInt16(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD int16
		wantR []byte
	}{
		{"Decode next int16", &figbuf.Decoder{}, args{[]byte{0x83, 0x9d, 0x9c, 0x01}}, int16(-9999), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextInt16(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextInt16() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextInt16() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextInt32(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD int32
		wantR []byte
	}{
		{"Decode next int32", &figbuf.Decoder{}, args{[]byte{0x85, 0xfd, 0xa7, 0xd6, 0xb9, 0x07}}, int32(-999999999), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextInt32(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextInt32() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextInt32() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextInt64(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD int64
		wantR []byte
	}{
		{"Decode next int64", &figbuf.Decoder{}, args{[]byte{0x88, 0xfd, 0xff, 0x87, 0xfc, 0xcd, 0xbc, 0xc3, 0x23}}, int64(-9999999999999999), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextInt64(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextInt64() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextInt64() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUint(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD uint
		wantR []byte
	}{
		{"Decode next uint", &figbuf.Decoder{}, args{[]byte{0x84, 0x3b, 0x9a, 0xc9, 0xff}}, uint(999999999), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextUint(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextUint() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUint() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUint8(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD uint8
		wantR []byte
	}{
		{"Decode next uint8", &figbuf.Decoder{}, args{[]byte{0x63}}, uint8(99), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextUint8(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextUint8() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUint8() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUint16(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD uint16
		wantR []byte
	}{
		{"Decode next uint16", &figbuf.Decoder{}, args{[]byte{0x82, 0x27, 0x0f}}, uint16(9999), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextUint16(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextUint16() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUint16() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUint32(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD uint32
		wantR []byte
	}{
		{"Decode next uint32", &figbuf.Decoder{}, args{[]byte{0x84, 0x3b, 0x9a, 0xc9, 0xff}}, uint32(999999999), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextUint32(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextUint32() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUint32() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUint64(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantD uint64
		wantR []byte
	}{
		{"Decode next uint64", &figbuf.Decoder{}, args{[]byte{0x87, 0x23, 0x86, 0xf2, 0x6f, 0xc0, 0xff, 0xff}}, uint64(9999999999999999), []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, gotR := tt.dec.DecodeNextUint64(tt.args.b)
			if gotD != tt.wantD {
				t.Errorf("Decoder.DecodeNextUint64() gotD = %v, want %v", gotD, tt.wantD)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUint64() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextStringSlice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []string
		wantR  []byte
	}{
		{"Decode next []string", &figbuf.Decoder{}, args{[]byte{0xd3, 0x85, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x85, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x86, 0x7a, 0x6f, 0x6d, 0x62, 0x69, 0x65}}, []string{"hello", "world", "zombie"}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextStringSlice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextStringSlice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextStringSlice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextBoolSlice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []bool
		wantR  []byte
	}{
		{"Decode next []bool", &figbuf.Decoder{}, args{[]byte{0xc3, 0x01, 0x01, 0x00}}, []bool{true, true, false}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextBoolSlice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextBoolSlice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextBoolSlice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextIntSlice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []int
		wantR  []byte
	}{
		{"Decode next []int", &figbuf.Decoder{}, args{[]byte{0xd0, 0x85, 0xfd, 0xa7, 0xd6, 0xb9, 0x07, 0x84, 0xd4, 0xe6, 0xad, 0x06, 0x84, 0xa9, 0xf3, 0x96, 0x03}}, []int{-999999999, 6666666, -3333333}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextIntSlice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextIntSlice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextIntSlice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextInt8Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []int8
		wantR  []byte
	}{
		{"Decode next []int8", &figbuf.Decoder{}, args{[]byte{0xc7, 0x82, 0xc5, 0x01, 0x82, 0x84, 0x01, 0x41}}, []int8{-99, 66, -33}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextInt8Slice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextInt8Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextInt8Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextInt16Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []int16
		wantR  []byte
	}{
		{"Decode next []int16", &figbuf.Decoder{}, args{[]byte{0xca, 0x83, 0x9d, 0x9c, 0x01, 0x82, 0x94, 0x68, 0x82, 0x89, 0x34}}, []int16{-9999, 6666, -3333}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextInt16Slice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextInt16Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextInt16Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextInt32Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []int32
		wantR  []byte
	}{
		{"Decode next []int32", &figbuf.Decoder{}, args{[]byte{0xd0, 0x85, 0xfd, 0xa7, 0xd6, 0xb9, 0x07, 0x84, 0xd4, 0xe6, 0xad, 0x06, 0x84, 0xa9, 0xf3, 0x96, 0x03}}, []int32{-999999999, 6666666, -3333333}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextInt32Slice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextInt32Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextInt32Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextInt64Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []int64
		wantR  []byte
	}{
		{"Decode next []int64", &figbuf.Decoder{}, args{[]byte{0xd7, 0x88, 0xfd, 0xff, 0xb3, 0xcc, 0xd4, 0xdf, 0xc6, 0x03, 0x86, 0xd4, 0xea, 0xa4, 0xda, 0xf0, 0x03, 0x86, 0xa9, 0xb5, 0x92, 0xad, 0xf8, 0x01}}, []int64{-999999999999999, 66666666666, -33333333333}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextInt64Slice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextInt64Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextInt64Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUintSlice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []uint
		wantR  []byte
	}{
		{"Decode next []uint", &figbuf.Decoder{}, args{[]byte{0xcc, 0x83, 0x98, 0x96, 0x7f, 0x83, 0x0a, 0x2c, 0x2a, 0x83, 0x05, 0x16, 0x15}}, []uint{9999999, 666666, 333333}, []byte{}},
		{"Decode next []uint", &figbuf.Decoder{}, args{[]byte{0xc3, 0x09, 0x06, 0x03}}, []uint{9, 6, 3}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextUintSlice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextUintSlice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUintSlice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUint8Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []uint8
		wantR  []byte
	}{
		{"Decode next []uint8", &figbuf.Decoder{}, args{[]byte{0x83, 0x63, 0x42, 0x21}}, []uint8{99, 66, 33}, []byte{}},
		{"Decode next []uint8", &figbuf.Decoder{}, args{[]byte{0x83, 0x09, 0x06, 0x03}}, []uint8{9, 6, 3}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextUint8Slice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextUint8Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUint8Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUint16Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []uint16
		wantR  []byte
	}{
		{"Decode next []uint16", &figbuf.Decoder{}, args{[]byte{0xc9, 0x82, 0x27, 0x0f, 0x82, 0x02, 0x9a, 0x82, 0x01, 0x4d}}, []uint16{9999, 666, 333}, []byte{}},
		{"Decode next []uint16", &figbuf.Decoder{}, args{[]byte{0xc3, 0x09, 0x06, 0x03}}, []uint16{9, 6, 3}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextUint16Slice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextUint16Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUint16Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUint32Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []uint32
		wantR  []byte
	}{
		{"Decode next []uint32", &figbuf.Decoder{}, args{[]byte{0xcc, 0x83, 0x98, 0x96, 0x7f, 0x83, 0x0a, 0x2c, 0x2a, 0x83, 0x05, 0x16, 0x15}}, []uint32{9999999, 666666, 333333}, []byte{}},
		{"Decode next []uint32", &figbuf.Decoder{}, args{[]byte{0xc3, 0x09, 0x06, 0x03}}, []uint32{9, 6, 3}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextUint32Slice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextUint32Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUint32Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextUint64Slice(t *testing.T) {
	type args struct {
		bb []byte
	}
	tests := []struct {
		name   string
		dec    *figbuf.Decoder
		args   args
		wantDd []uint64
		wantR  []byte
	}{
		{"Decode next []uint64", &figbuf.Decoder{}, args{[]byte{0xd5, 0x87, 0x03, 0x8d, 0x7e, 0xa4, 0xc6, 0x7f, 0xff, 0x85, 0x0f, 0x85, 0xa4, 0x9a, 0xaa, 0x86, 0x03, 0x08, 0x1a, 0x26, 0x35, 0x55}}, []uint64{999999999999999, 66666666666, 3333333333333}, []byte{}},
		{"Decode next []uint64", &figbuf.Decoder{}, args{[]byte{0xc3, 0x09, 0x06, 0x03}}, []uint64{9, 6, 3}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDd, gotR := tt.dec.DecodeNextUint64Slice(tt.args.bb)
			if !reflect.DeepEqual(gotDd, tt.wantDd) {
				t.Errorf("Decoder.DecodeNextUint64Slice() gotDd = %v, want %v", gotDd, tt.wantDd)
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextUint64Slice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextBinaryUnmarshaler(t *testing.T) {
	type args struct {
		b    []byte
		dest encoding.BinaryUnmarshaler
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantR []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := tt.dec.DecodeNextBinaryUnmarshaler(tt.args.b, tt.args.dest); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextBinaryUnmarshaler() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_DecodeNextTextUnmarshaler(t *testing.T) {
	type args struct {
		b    []byte
		dest encoding.TextUnmarshaler
	}
	tests := []struct {
		name  string
		dec   *figbuf.Decoder
		args  args
		wantR []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := tt.dec.DecodeNextTextUnmarshaler(tt.args.b, tt.args.dest); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Decoder.DecodeNextTextUnmarshaler() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestDecoder_BytesToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want string
	}{
		{"Convert bytes to string", &figbuf.Decoder{}, args{[]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64}}, "Hello World"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToString(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToBool(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want bool
	}{
		{"Convert bytes to bool", &figbuf.Decoder{}, args{[]byte{0x01}}, true},
		{"Convert bytes to bool", &figbuf.Decoder{}, args{[]byte{0x00}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToBool(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToInt(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want int
	}{
		{"Convert bytes to int", &figbuf.Decoder{}, args{[]byte{0xe1, 0xf2, 0xa0, 0x4d}}, -81009841},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToInt(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToInt8(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want int8
	}{
		{"Convert bytes to int16", &figbuf.Decoder{}, args{[]byte{0xa1, 0x01}}, -81},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToInt8(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToInt8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToInt16(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want int16
	}{
		{"Convert bytes to int16", &figbuf.Decoder{}, args{[]byte{0xc7, 0x7e}}, -8100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToInt16(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToInt32(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want int32
	}{
		{"Convert bytes to int32", &figbuf.Decoder{}, args{[]byte{0xe1, 0xf2, 0xa0, 0x4d}}, -81009841},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToInt32(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToInt64(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want int64
	}{
		{"Convert bytes to int64", &figbuf.Decoder{}, args{[]byte{0x9f, 0x88, 0x86, 0xda, 0x93, 0x2f}}, -810098410000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToInt64(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToUint(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want uint
	}{
		{"Convert bytes to uint", &figbuf.Decoder{}, args{[]byte{0x30, 0x49, 0x1e, 0xea}}, 810098410},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToUint(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToUint8(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want uint8
	}{
		{"Convert bytes to uint8", &figbuf.Decoder{}, args{[]byte{0x0a}}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToUint8(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToUint16(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want uint16
	}{
		// TODO: Add test cases.
		{"Convert bytes to uint16", &figbuf.Decoder{}, args{[]byte{0x20, 0xda}}, 8410},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToUint16(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToUint32(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want uint32
	}{
		{"Convert bytes to uint32", &figbuf.Decoder{}, args{[]byte{0x30, 0x49, 0x1e, 0xea}}, 810098410},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToUint32(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_BytesToUint64(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		dec  *figbuf.Decoder
		args args
		want uint64
	}{
		{"Convert bytes to uint64", &figbuf.Decoder{}, args{[]byte{0xf6, 0x3c, 0x6c, 0x76, 0x52, 0xde, 0xea}}, 69309280810098410},
		{"Convert bytes to uint64", &figbuf.Decoder{}, args{[]byte{0x0f, 0x85, 0xa4, 0x9a, 0xaa}}, 66666666666},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.BytesToUint64(tt.args.b); got != tt.want {
				t.Errorf("Decoder.BytesToUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}
