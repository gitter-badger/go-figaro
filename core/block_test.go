package core

import (
	"reflect"
	"testing"
	"time"
)

func TestNewBlock(t *testing.T) {
	type args struct {
		index        uint64
		timestamp    time.Time
		data         string
		previousHash string
	}
	tests := []struct {
		name string
		args args
		want *Block
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlock(tt.args.index, tt.args.timestamp, tt.args.data, tt.args.previousHash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_Seal(t *testing.T) {
	tests := []struct {
		name    string
		b       *Block
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Seal(); (err != nil) != tt.wantErr {
				t.Errorf("Block.Seal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlock_Validate(t *testing.T) {
	tests := []struct {
		name string
		b    *Block
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Validate(); got != tt.want {
				t.Errorf("Block.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blockHash(t *testing.T) {
	type args struct {
		b *Block
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := blockHash(tt.args.b); got != tt.want {
				t.Errorf("blockHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
