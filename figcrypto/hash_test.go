package figcrypto_test

import (
	"fmt"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto"
)

func ExampleHash() {
	archive := []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}
	h := figcrypto.Hash(archive)
	fmt.Printf("%#x\n", h)
	// Output: 0xcfae1696d66549c090e4d6ba1266c7b0aa7bb80747c642d8607ecb5d2dec80a5
}

func BenchmarkHash(b *testing.B) {
	archive := []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}
	for i := 0; i < b.N; i++ {
		figcrypto.Hash(archive)
	}
}

func BenchmarkHasher_Hash(b *testing.B) {
	archive := []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}
	hasher := figcrypto.NewHasher()
	for i := 0; i < b.N; i++ {
		hasher.Hash(archive)
	}
}
