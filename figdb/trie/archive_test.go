package trie_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/figaro-tech/go-figaro/figdb/mock"
	"github.com/figaro-tech/go-figaro/figdb/trie"
)

func ExampleArchive_Save() {
	tr := &trie.Archive{
		KeyStore: mock.NewKeyStore(),
	}
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	root, err := tr.Save(archive)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("% #x\n", root)
	// Output: 0x45 0x0a 0x13 0x0d 0x18 0x53 0x03 0xd5 0x3e 0x63 0xd7 0xd8 0x23 0x21 0x1f 0x11 0x7b 0xea 0x61 0x61 0x4b 0xc5 0xe9 0x72 0x6a 0x81 0xda 0xff 0x93 0xa9 0xe5 0x88
}

func ExampleArchive_proof() {
	tr := &trie.Archive{
		KeyStore: mock.NewKeyStore(),
	}
	trv := &trie.ArchiveValidator{}
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	index := 6
	root, err := tr.Save(archive)
	value, proof, err := tr.GetAndProve(root, index)
	valid := trv.Validate(root, index, value, proof)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(valid)
	// Output: true
}

func BenchmarkArchive_Save(b *testing.B) {
	tr := &trie.Archive{
		KeyStore: mock.NewKeyStore(),
	}
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	var err error
	for i := 0; i < b.N; i++ {
		_, err = tr.Save(archive)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkArchive_Retrieve(b *testing.B) {
	tr := &trie.Archive{
		KeyStore: mock.NewKeyStore(),
	}
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	root, _ := tr.Save(archive)
	var err error
	for i := 0; i < b.N; i++ {
		_, err = tr.Retrieve(root)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkArchive_Get(b *testing.B) {
	tr := &trie.Archive{
		KeyStore: mock.NewKeyStore(),
	}
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	root, _ := tr.Save(archive)
	var err error
	for i := 0; i < b.N; i++ {
		_, err = tr.Get(root, 0)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkArchive_GetAndProve(b *testing.B) {
	tr := &trie.Archive{
		KeyStore: mock.NewKeyStore(),
	}
	archive := [][]byte{{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}}
	root, _ := tr.Save(archive)
	var err error
	for i := 0; i < b.N; i++ {
		_, _, err = tr.GetAndProve(root, 0)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkArchiveValidator_Validate(b *testing.B) {
	tr := &trie.Archive{
		KeyStore: mock.NewKeyStore(),
	}
	trv := &trie.ArchiveValidator{}
	archive := [][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}
	root, _ := tr.Save(archive)
	value, proof, _ := tr.GetAndProve(root, 0)
	var r bool
	for i := 0; i < b.N; i++ {
		r = trv.Validate(root, 0, value, proof)
		if r != true {
			log.Fatal("failed validation")
		}
	}
}

func TestArchive_All(t *testing.T) {
	type args struct {
		batch [][]byte
		index int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"A basic test of everything", args{[][]byte{{0xff, 0xee}, {0xdd, 0xcc}}, 1}, false},
		{"A test case that failed once", args{[][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &trie.Archive{
				KeyStore: mock.NewKeyStore(),
			}
			root, err := tr.Save(tt.args.batch)
			if (err != nil) != tt.wantErr {
				t.Errorf("Archive.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			archive, err := tr.Retrieve(root)
			if (err != nil) != tt.wantErr {
				t.Errorf("Archive.Retrieve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(archive, tt.args.batch) {
				t.Errorf("Archive.Retrieve() = %v, want %v", archive, tt.args.batch)
				return
			}
			got, err := tr.Get(root, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Archive.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.args.batch[tt.args.index]) {
				t.Errorf("Archive.Get() = %v, want %v", got, tt.args.batch[tt.args.index])
			}
		})
	}
}

func TestArchive_GetAndProve_TestArchiveValidator_Validate(t *testing.T) {
	type args struct {
		batch [][]byte
		index int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"A basic test of everything", args{[][]byte{{0xff, 0xee}, {0xdd, 0xcc}}, 1}, true, false},
		{"A test case that failed once", args{[][]byte{{0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}, {0xff, 0xee}, {0xdd, 0xcc}, {0xbb, 0xaa}}, 1}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &trie.Archive{
				KeyStore: mock.NewKeyStore(),
			}
			root, err := tr.Save(tt.args.batch)
			if (err != nil) != tt.wantErr {
				t.Errorf("Archive.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, got1, err := tr.GetAndProve(root, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Archive.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.args.batch[tt.args.index]) {
				t.Errorf("Archive.Get() = %v, want %v", got, tt.args.batch[tt.args.index])
			}
			trv := &trie.ArchiveValidator{}
			if valid := trv.Validate(root, tt.args.index, got, got1); valid != tt.want {
				t.Errorf("ArchiveValidator.Validate() = %v, want %v", valid, tt.want)
			}
		})
	}
}
