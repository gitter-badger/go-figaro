package figcrypto_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto"
)

func ExampleBMTrieRoot() {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	root := figcrypto.BMTrieRoot(archive)
	fmt.Printf("%#x\n", root)
	// Output: 0xb4cb24d832073729b9626963ea55bc4d621b108348335b452e75eabf2a66aeee
}

func ExampleBMTrieProof() {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	index := 6
	proof, err := figcrypto.BMTrieProof(archive, index)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#x\n", proof)
	// Output: [0x772544556b5de9b4ea6e9f083a9066e17c26385785ba94ddf33213c7538783a6 0xdd6aa78d426ce9647c45069cadc865a1bb9432823fa8b5e0d26f7de4f1c4d6eb 0xb415e5473e716d62ba57995ea4e8da2fddade4ed0bbe7623e1b631c2906392bd 0x913f9bd357c60f2b8999901fe978cec299c77fa078ded8d2c4272c869493e09f 0x770b6dfda8b9dab0b25bf828a50cb7e420c34b6dfda458bcc93b5b2e0e4f0e60 0x34ab56f0be4c128b0aa768a01ca0cd89effccfc7c365902e264ecde8879d3267]
}

func ExampleBMTrieValidate() {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	index := 6
	root := figcrypto.BMTrieRoot(archive)
	proof, err := figcrypto.BMTrieProof(archive, index)
	if err != nil {
		log.Fatal(err)
	}
	valid := figcrypto.BMTrieValidate(root, index, archive[index], proof)
	fmt.Println(valid)
	// Output: true
}

func BenchmarkBMTrieRoot(b *testing.B) {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	for i := 0; i < b.N; i++ {
		figcrypto.BMTrieRoot(archive)
	}
}

func BenchmarkBMTrieProof(b *testing.B) {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	for i := 0; i < b.N; i++ {
		figcrypto.BMTrieProof(archive, 0)
	}
}

func BenchmarkBMTrieValidate(b *testing.B) {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	root := figcrypto.BMTrieRoot(archive)
	proof, err := figcrypto.BMTrieProof(archive, 0)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		figcrypto.BMTrieValidate(root, 0, archive[0], proof)
	}
}

func TestBMTrieRoot(t *testing.T) {
	type args struct {
		data [][]byte
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
			if got := figcrypto.BMTrieRoot(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BMTrie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBMTrieProof(t *testing.T) {
	type args struct {
		data  [][]byte
		index int
	}
	tests := []struct {
		name    string
		args    args
		want    [][]byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := figcrypto.BMTrieProof(tt.args.data, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("BMTrieProof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BMTrieProof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBMTrieValidate(t *testing.T) {
	type args struct {
		root  []byte
		index int
		data  []byte
		proof [][]byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := figcrypto.BMTrieValidate(tt.args.root, tt.args.index, tt.args.data, tt.args.proof); got != tt.want {
				t.Errorf("BMTrieValidate() = %v, want %v", got, tt.want)
			}
		})
	}
}
