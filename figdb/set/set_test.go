package set_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/figaro-tech/go-figaro/figdb/mock"

	"github.com/figaro-tech/go-figaro/figdb/set"
)

func ExampleSet_Save() {
	bb := set.Set{
		KeyStore: mock.NewKeyStore(),
	}
	data := [][]byte{[]byte("dog"), []byte("doge"), []byte("coin")}
	key, _ := bb.Save(data, 0.01)
	fmt.Printf("Dog is in set: %t\n", bb.Has(key, []byte("dog")))

	// Output:
	// Dog is in set: true
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456000789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func BenchmarkSet_Save(b *testing.B) {
	bb := set.Set{
		KeyStore: mock.NewKeyStore(),
	}
	data := make([][]byte, 56000)
	for i := range data {
		data[i] = []byte(RandomString(5))
	}
	for i := 0; i < b.N; i++ {
		bb.Save(data, 0.01)
	}
}

func BenchmarkSet_Has(b *testing.B) {
	bb := set.Set{
		KeyStore: mock.NewKeyStore(),
	}
	data := make([][]byte, 56000)
	for i := range data {
		data[i] = []byte(RandomString(5))
	}
	key, _ := bb.Save(data, 0.01)
	for i := 0; i < b.N; i++ {
		bb.Has(key, data[0])
	}
}

func BenchmarkSet_HasBatch(b *testing.B) {
	bb := set.Set{
		KeyStore: mock.NewKeyStore(),
	}
	data := make([][]byte, 56000)
	for i := range data {
		data[i] = []byte(RandomString(5))
	}
	key, _ := bb.Save(data, 0.01)
	for i := 0; i < b.N; i++ {
		bb.HasBatch(key, data)
	}
}

func BenchmarkSet_Get_and_test(b *testing.B) {
	bb := set.Set{
		KeyStore: mock.NewKeyStore(),
	}
	data := make([][]byte, 56000)
	for i := range data {
		data[i] = []byte(RandomString(5))
	}
	key, _ := bb.Save(data, 0.01)
	bloom, _ := bb.Get(key)
	for i := 0; i < b.N; i++ {
		bloom.Has(data[0])
	}
}
