package crypto

import (
	"crypto/sha256"
	"reflect"
	"testing"
)

func TestSha256(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want [sha256.Size]byte
	}{
		{"test", args{[]byte("my name is inigo montoya")}, [32]byte{0xad, 0xfe, 0xb4, 0xd7, 0x15, 0x5a, 0x3f, 0xd8, 0x87, 0x0b, 0x6c, 0x6a, 0xef, 0xab, 0x83, 0x57, 0x32, 0x79, 0xb3, 0x29, 0x1a, 0x17, 0x3b, 0x88, 0x74, 0x79, 0xf8, 0xad, 0x1d, 0x9f, 0x46, 0x59}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha256(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sha256() = %x, want %xs", got, tt.want)
			}
		})
	}
}

func TestSha3(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want [32]byte
	}{
		{"test", args{[]byte("my name is inigo montoya")}, [32]byte{0x60, 0x73, 0x85, 0x1e, 0x23, 0xfa, 0xbc, 0x64, 0x64, 0x9a, 0x28, 0xb3, 0x8b, 0x08, 0x62, 0xd3, 0x7b, 0x08, 0xd7, 0x5e, 0xc0, 0x63, 0xa0, 0xe5, 0xc1, 0xea, 0xec, 0x7c, 0x80, 0x50, 0x6a, 0x71}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha3(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sha3() = %x, want %x", got, tt.want)
			}
		})
	}
}

func BenchmarkSha256(b *testing.B) {
	t := []byte("my name is inigo montoya")
	for i := 0; i < b.N; i++ {
		Sha256(t)
	}
}

func BenchmarkSha3(b *testing.B) {
	t := []byte("my name is inigo montoya")
	for i := 0; i < b.N; i++ {
		Sha3(t)
	}
}
