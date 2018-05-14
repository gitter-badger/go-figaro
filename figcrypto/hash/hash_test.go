package hash_test

import (
	"fmt"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto/hash"
)

func ExampleHash256() {
	archive := []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}
	h := hash.Hash256(archive)
	fmt.Printf("%#x\n", h)
	// Output: 0xcfae1696d66549c090e4d6ba1266c7b0aa7bb80747c642d8607ecb5d2dec80a5
}

func BenchmarkHash256(b *testing.B) {
	archive := []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}
	for i := 0; i < b.N; i++ {
		hash.Hash256(archive)
	}
}
