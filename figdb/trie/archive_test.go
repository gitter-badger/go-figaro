package trie_test

import (
	"reflect"
	"testing"

	"github.com/figaro-tech/go-figaro/figdb/mock"
	"github.com/figaro-tech/go-figaro/figdb/trie"
)

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

func TestArchive_GetAndProve(t *testing.T) {
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
			if valid := trie.ValidateBMT(root, tt.args.index, got, got1); valid != tt.want {
				t.Errorf("ArchiveValidator.Validate() = %v, want %v", valid, tt.want)
			}
		})
	}
}
