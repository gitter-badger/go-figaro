package level_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/figaro-tech/figaro/pkg/figdb/internal"
	"github.com/figaro-tech/figaro/pkg/figdb/internal/level"
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
			s := level.NewStore("testdata")
			defer s.Close()
			s.Set(tt.args.key, tt.args.value)
			if got := s.Get(tt.args.key); !reflect.DeepEqual(got, tt.args.value) {
				t.Errorf("Store.Get() after Store.Set() = %v, want %v", got, tt.args.value)
			}
			if got := s.Get(tt.args.key); !reflect.DeepEqual(got, tt.args.value) {
				t.Errorf("Store.Get() after Store.Delete() = %v, want %v", got, nil)
			}
		})
		helperCleanup()
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
			s := level.NewStore("testdata")
			defer s.Close()
			s.Batch()
		})
		helperCleanup()
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
			s := level.NewStore("testdata")
			defer s.Close()
			s.Write()
		})
		helperCleanup()
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
			s := level.NewStore("testdata")
			defer s.Close()
			s.BatchUpdate(tt.args.updates)
		})
		helperCleanup()
	}
}

func helperCleanup() {
	cwd, _ := os.Getwd()
	directory := cwd + "/testdata/"
	dirRead, _ := os.Open(directory)
	dirFiles, _ := dirRead.Readdir(0)
	for index := range dirFiles {
		fileHere := dirFiles[index]
		nameHere := fileHere.Name()
		if strings.HasSuffix(nameHere, ".keep") {
			continue
		}
		fullPath := directory + nameHere
		os.Remove(fullPath)
	}
}
