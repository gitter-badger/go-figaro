package trie_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto/trie"
)

func ExampleTrie() {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	root := trie.Trie(archive)
	fmt.Printf("%#x\n", root)
	// Output: 0xb4cb24d832073729b9626963ea55bc4d621b108348335b452e75eabf2a66aeee
}

func ExampleProof() {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	index := 6
	proof, err := trie.Proof(archive, index)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#x\n", proof)
	// Output: [0xdd6aa78d426ce9647c45069cadc865a1bb9432823fa8b5e0d26f7de4f1c4d6eb 0xb415e5473e716d62ba57995ea4e8da2fddade4ed0bbe7623e1b631c2906392bd 0x913f9bd357c60f2b8999901fe978cec299c77fa078ded8d2c4272c869493e09f 0x770b6dfda8b9dab0b25bf828a50cb7e420c34b6dfda458bcc93b5b2e0e4f0e60 0x34ab56f0be4c128b0aa768a01ca0cd89effccfc7c365902e264ecde8879d3267]
}

func ExampleValidate() {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	index := 6
	root := trie.Trie(archive)
	proof, err := trie.Proof(archive, index)
	if err != nil {
		log.Fatal(err)
	}
	valid := trie.Validate(root, index, archive[index], proof)
	fmt.Println(valid)
	// Output: true
}

func BenchmarkTrie(b *testing.B) {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	for i := 0; i < b.N; i++ {
		trie.Trie(archive)
	}
}

func BenchmarkProof(b *testing.B) {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	for i := 0; i < b.N; i++ {
		trie.Proof(archive, 0)
	}
}

func BenchmarkValidate(b *testing.B) {
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	root := trie.Trie(archive)
	proof, err := trie.Proof(archive, 0)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		trie.Validate(root, 0, archive[0], proof)
	}
}

func TestTrie(t *testing.T) {
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
			if got := trie.Trie(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProof(t *testing.T) {
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
			got, err := trie.Proof(tt.args.data, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Proof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Proof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
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
			if got := trie.Validate(tt.args.root, tt.args.index, tt.args.data, tt.args.proof); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
