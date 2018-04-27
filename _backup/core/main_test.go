package core

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

const LOREM = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func generateRandomBlock() *Block {
	return &Block{
		Index:        uint64(rand.Intn(100)),
		TimeStamp:    time.Now(),
		Data:         generateRandomData(),
		PreviousHash: generateRandomHash(),
	}
}

func generateRandomData() string {
	return LOREM[rand.Intn(100):]
}

func generateRandomHash() string {
	return fmt.Sprintf("0x%x", sha256.New().Sum([]byte(LOREM[rand.Intn(100):])))
}

func copyBlock(block *Block) *Block {
	clone := *block
	return &clone
}
