package core

import (
	"reflect"
	"testing"
)

func TestNewBlockChain(t *testing.T) {
	tests := []struct {
		name string
		want *BlockChain
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlockChain(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlockChain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockChain_Head(t *testing.T) {
	tests := []struct {
		name string
		c    *BlockChain
		want *Block
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Head(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BlockChain.Head() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockChain_CreateBlock(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		c    *BlockChain
		args args
		want *Block
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.CreateBlock(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BlockChain.CreateBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockChain_ReceiveBlock(t *testing.T) {
	type args struct {
		block *Block
	}
	tests := []struct {
		name    string
		c       *BlockChain
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.ReceiveBlock(tt.args.block); (err != nil) != tt.wantErr {
				t.Errorf("BlockChain.ReceiveBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
