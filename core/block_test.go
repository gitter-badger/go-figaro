package core

import (
	"reflect"
	"testing"
	"time"
)

const YYYY_MM_DD = "2006-01-02"

func TestNewBlock(t *testing.T) {
	sealedBlock := generateRandomBlock().Seal()

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
		{
			"ReturnsASealedBlock",
			args{
				sealedBlock.Index,
				sealedBlock.TimeStamp,
				sealedBlock.Data,
				sealedBlock.PreviousHash,
			},
			sealedBlock,
		},
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
	unsealedBlock := generateRandomBlock()

	sealedBlock := copyBlock(unsealedBlock)
	sealedBlock.Hash = blockHash(sealedBlock)

	tests := []struct {
		name string
		b    *Block
		want *Block
	}{
		{"unsealedBlock", unsealedBlock, sealedBlock},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Seal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Block.Seal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_Valid(t *testing.T) {
	unalteredBlock := generateRandomBlock().Seal()

	alteredIndex := generateRandomBlock().Seal()
	alteredIndex.Index = 123

	alteredTimeStamp := generateRandomBlock().Seal()
	alteredTimeStamp.TimeStamp = time.Now()

	alteredData := generateRandomBlock().Seal()
	alteredData.Data = generateRandomData()

	alteredPrevHash := generateRandomBlock().Seal()
	alteredPrevHash.PreviousHash = generateRandomHash()

	tests := []struct {
		name string
		b    *Block
		want bool
	}{
		{"NoAlterations", unalteredBlock, true},
		{"AlteredIndex", alteredIndex, false},
		{"AlteredTimeStamp", alteredTimeStamp, false},
		{"AlteredData", alteredData, false},
		{"AlteredPreviousHash", alteredPrevHash, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Valid(); got != tt.want {
				t.Errorf("Block.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blockHash(t *testing.T) {
	blockOrig := generateRandomBlock()

	blockSame := copyBlock(blockOrig)

	blockDiffIndex := copyBlock(blockOrig)
	blockDiffIndex.Index = 123

	diffTime, _ := time.Parse(YYYY_MM_DD, "2001-01-01")
	blockDiffTime := copyBlock(blockOrig)
	blockDiffTime.TimeStamp = diffTime

	blockDiffData := copyBlock(blockOrig)
	blockDiffData.Data = "Hal"

	blockDiffPrev := copyBlock(blockOrig)
	blockDiffPrev.PreviousHash = "3000"

	type args struct {
		b *Block
	}
	tests := []struct {
		name   string
		args   args
		test   string
		equals bool
	}{
		{"WithSame", args{blockOrig}, blockHash(blockOrig), true},
		{"WithSameValues", args{blockOrig}, blockHash(blockSame), true},
		{"WithDiffIndex", args{blockDiffIndex}, blockHash(blockOrig), false},
		{"WithDiffTimes", args{blockDiffTime}, blockHash(blockOrig), false},
		{"WithDiffData", args{blockDiffData}, blockHash(blockOrig), false},
		{"WithDiffPrev", args{blockDiffPrev}, blockHash(blockOrig), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := blockHash(tt.args.b); (got == tt.test) != tt.equals {
				t.Errorf("blockHash(), (%v == %v) == %v", got, tt.test, tt.equals)
			}
		})
	}
}
