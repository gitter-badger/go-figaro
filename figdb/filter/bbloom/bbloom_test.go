package bbloom_test

import (
	"fmt"

	"github.com/figaro-tech/go-figaro/figdb/filter/bbloom"
)

func ExampleBloom_Add() {
	bloom := bbloom.NewWithEstimates(10, 0.01)
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
