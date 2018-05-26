package bloom_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/figaro-tech/go-figaro/figdb/bloom"
)

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456000789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
func ExampleBloom_Add() {
	bloom := bloom.NewWithEstimates(10, 0.01)
	data := [][]byte{[]byte("dog"), []byte("doge"), []byte("coin")}
	for _, datum := range data {
		bloom.Add(datum)
	}
	fmt.Printf("Dog is in set: %t\n", bloom.Has([]byte("dog")))
	fmt.Printf("Horse is in set: %t\n", bloom.Has([]byte("horse")))

	// Output:
	// Dog is in set: true
	// Horse is in set: false
}

func BenchmarkBloom_Add(b *testing.B) {
	bloom := bloom.NewWithEstimates(56000, 0.03)
	for i := 0; i < 56000; i++ {
		bloom.Add([]byte(RandomString(5)))
	}
	test := []byte("happy")
	for i := 0; i < b.N; i++ {
		bloom.Add(test)
	}
}

func BenchmarkBloom_Marshal(b *testing.B) {
	bloom := bloom.NewWithEstimates(56000, 0.03)
	for i := 0; i < 56000; i++ {
		bloom.Add([]byte(RandomString(5)))
	}
	for i := 0; i < b.N; i++ {
		bloom.Marshal()
	}
}

func BenchmarkBloom_Unmarshal(b *testing.B) {
	bl := bloom.NewWithEstimates(5600, 0.03)
	for i := 0; i < 56000; i++ {
		bl.Add([]byte(RandomString(5)))
	}
	m, _ := bl.Marshal()

	for i := 0; i < b.N; i++ {
		bloom.Unmarshal(m)
	}
}
