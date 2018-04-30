package trie_test

import (
	"reflect"
	"testing"

	"github.com/figaro-tech/figaro/pkg/figdb/internal"
	"github.com/figaro-tech/figaro/pkg/figdb/internal/trie"
)

func TestState_Get(t *testing.T) {
	type fields struct {
		KeyStore internal.KeyStore
		Hasher   internal.Hasher
		Encdec   internal.EncoderDecoder
	}
	type args struct {
		root []byte
		key  []byte
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
			tr := &trie.State{
				KeyStore: tt.fields.KeyStore,
				Hasher:   tt.fields.Hasher,
				Encdec:   tt.fields.Encdec,
			}
			if got := tr.Get(tt.args.root, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("State.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_Set(t *testing.T) {
	type fields struct {
		KeyStore internal.KeyStore
		Hasher   internal.Hasher
		Encdec   internal.EncoderDecoder
	}
	type args struct {
		root  []byte
		key   []byte
		value []byte
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
			tr := &trie.State{
				KeyStore: tt.fields.KeyStore,
				Hasher:   tt.fields.Hasher,
				Encdec:   tt.fields.Encdec,
			}
			if got := tr.Set(tt.args.root, tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("State.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_Prove(t *testing.T) {
	type fields struct {
		KeyStore internal.KeyStore
		Hasher   internal.Hasher
		Encdec   internal.EncoderDecoder
	}
	type args struct {
		root []byte
		key  []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
		want1  [][][]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &trie.State{
				KeyStore: tt.fields.KeyStore,
				Hasher:   tt.fields.Hasher,
				Encdec:   tt.fields.Encdec,
			}
			got, got1 := tr.Prove(tt.args.root, tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("State.Prove() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("State.Prove() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestStateValidator_Validate(t *testing.T) {
	type fields struct {
		Hasher internal.Hasher
	}
	type args struct {
		root  []byte
		key   []byte
		data  []byte
		proof [][][]byte
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
			tr := &trie.StateValidator{
				Hasher: tt.fields.Hasher,
			}
			if got := tr.Validate(tt.args.root, tt.args.key, tt.args.data, tt.args.proof); got != tt.want {
				t.Errorf("StateValidator.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
