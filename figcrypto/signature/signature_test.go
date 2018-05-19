package signature_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto/signature"
)

func ExampleVerify() {
	seed := "hello darkness my old friend"
	priv, pub, addr, err := signature.GenerateFromSeed(seed)
	// priv, pub, addr, err := signature.GenerateKey(nil)
	if err != nil {
		log.Panic(err)
	}
	msg := []byte("hello world")
	sig, _ := signature.Sign(priv, msg)
	fmt.Printf("Private key: %#x\n", priv)
	fmt.Printf("Public key: %#x\n", pub)
	fmt.Printf("Address: %#x\n", addr)
	fmt.Printf("Human Address: %s\n", signature.ToHumanAddress(addr))
	fmt.Printf("Valid identity: %t\n", bytes.Equal(signature.Identify(sig, msg), addr))
	fmt.Printf("Valid with pubkey: %t\n", signature.Verify(pub, sig, msg))
	fmt.Printf("Valid with addr: %t\n", signature.VerifyWithAddress(addr, sig, msg))

	_, wpub, waddr, _ := signature.GenerateKey(nil)
	fmt.Printf("Valid with wrong pubkey: %t\n", signature.Verify(wpub, sig, msg))
	fmt.Printf("Valid with wrong addr: %t\n", signature.VerifyWithAddress(waddr, sig, msg))

	// Output:
	// Private key: 0xfa98dc283ee4e866809227266e23c63048af7103d43fd65576a67269f8299f21
	// Public key: 0x01eace460e3fe0e07725f2e01ce2ee7e6fcb21078d8933841a5df5d980d7126bd7
	// Address: 0x0066423b3664eb7d98de3f77ceb32e49e25b4317787535351b
	// Human Address: 1AKhGQ5vTbMdSPDLYaueVAtgaNwS6mYkK4
	// Valid identity: true
	// Valid with pubkey: true
	// Valid with addr: true
	// Valid with wrong pubkey: false
	// Valid with wrong addr: false
}

func BenchmarkIdentify(b *testing.B) {
	priv, _, _, _ := signature.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := signature.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		signature.Identify(sig, msg)
	}
}

func BenchmarkVerifyWithKey(b *testing.B) {
	priv, pub, _, _ := signature.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := signature.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		signature.Verify(pub, sig, msg)
	}
}
