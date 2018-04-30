package mem_test

import (
	"reflect"
	"testing"

	"github.com/figaro-tech/figaro/pkg/figdb/internal"
	"github.com/figaro-tech/figaro/pkg/figdb/internal/mem"
)

func TestStore_Set_Get_Delete(t *testing.T) {
	type args struct {
		key   []byte
		value []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test basic functionality", args{[]byte("key"), []byte("value")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mem.NewStore()
			s.Set(tt.args.key, tt.args.value)
			if got := s.Get(tt.args.key); !reflect.DeepEqual(got, tt.args.value) {
				t.Errorf("Store.Get() after Store.Set() = %v, want %v", got, tt.args.value)
			}
			if got := s.Get(tt.args.key); !reflect.DeepEqual(got, tt.args.value) {
				t.Errorf("Store.Get() after Store.Delete() = %v, want %v", got, nil)
			}
		})
	}
}

func TestStore_Batch(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mem.NewStore()
			s.Batch()
		})
	}
}

func TestStore_Write(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mem.NewStore()
			s.Write()
		})
	}
}

func TestStore_BatchUpdate(t *testing.T) {
	type args struct {
		updates internal.KeyStoreUpdateBatch
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mem.NewStore()
			s.BatchUpdate(tt.args.updates)
		})
	}
}
