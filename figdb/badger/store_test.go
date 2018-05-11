package badger_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/figaro-tech/go-figaro/figdb/badger"
)

func TestKeyStore_Set_Get_Delete(t *testing.T) {
	defer helperCleanup()
	type args struct {
		key   []byte
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test basic functionality", args{[]byte("key"), []byte("value")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks := badger.NewKeyStore("testdata")
			defer ks.Close()
			err := ks.Set(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyStore.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := ks.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyStore.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.args.value) {
				t.Errorf("KeyStore.Get() = %v, want %v", got, tt.args.value)
			}
			err = ks.Delete(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err = ks.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyStore.Get() after Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				t.Errorf("KeyStore.Get() after Delete() = %v, want %v", got, nil)
			}
		})
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
