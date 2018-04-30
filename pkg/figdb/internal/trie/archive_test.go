package trie_test

import (
	"reflect"
	"testing"

	"github.com/figaro-tech/figaro/pkg/figdb/internal"
	"github.com/figaro-tech/figaro/pkg/figdb/internal/trie"
)

func TestArchive_Save(t *testing.T) {
	type fields struct {
		KeyStore internal.KeyStore
		Hasher   internal.Hasher
		Encdec   internal.EncoderDecoder
	}
	type args struct {
		batch [][]byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &trie.Archive{
				KeyStore: tt.fields.KeyStore,
				Hasher:   tt.fields.Hasher,
				Encdec:   tt.fields.Encdec,
			}
			if got := tr.Save(tt.args.batch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Archive.Save() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchive_Retrieve(t *testing.T) {
	type fields struct {
		KeyStore internal.KeyStore
		Hasher   internal.Hasher
		Encdec   internal.EncoderDecoder
	}
	type args struct {
		root []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &trie.Archive{
				KeyStore: tt.fields.KeyStore,
				Hasher:   tt.fields.Hasher,
				Encdec:   tt.fields.Encdec,
			}
			if got := tr.Retrieve(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Archive.Retrieve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchive_Get(t *testing.T) {
	type fields struct {
		KeyStore internal.KeyStore
		Hasher   internal.Hasher
		Encdec   internal.EncoderDecoder
	}
	type args struct {
		root  []byte
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &trie.Archive{
				KeyStore: tt.fields.KeyStore,
				Hasher:   tt.fields.Hasher,
				Encdec:   tt.fields.Encdec,
			}
			if got := tr.Get(tt.args.root, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Archive.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchive_Prove(t *testing.T) {
	type fields struct {
		KeyStore internal.KeyStore
		Hasher   internal.Hasher
		Encdec   internal.EncoderDecoder
	}
	type args struct {
		root  []byte
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
		want1  [][]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &trie.Archive{
				KeyStore: tt.fields.KeyStore,
				Hasher:   tt.fields.Hasher,
				Encdec:   tt.fields.Encdec,
			}
			got, got1 := tr.Prove(tt.args.root, tt.args.index)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Archive.Prove() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Archive.Prove() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestArchiveValidator_Validate(t *testing.T) {
	type fields struct {
		Hasher internal.Hasher
	}
	type args struct {
		root  []byte
		index int
		data  []byte
		proof [][]byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &trie.ArchiveValidator{
				Hasher: tt.fields.Hasher,
			}
			if got := tr.Validate(tt.args.root, tt.args.index, tt.args.data, tt.args.proof); got != tt.want {
				t.Errorf("ArchiveValidator.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
