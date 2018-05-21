package signature_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto/signature"
)

func ExampleVerify() {
	seed := "hello darkness my old friend"
	priv, pub, addr, err := signature.GenerateFromSeed(seed)
	if err != nil {
		log.Panic(err)
	}
	msg := []byte("hello world")
	sig, _ := signature.Sign(priv, msg)
	fmt.Printf("Private key: %#x\n", priv)
	fmt.Printf("Public key: %#x\n", pub)
	fmt.Printf("Address: %#x\n", addr)
	fmt.Printf("Human Address: %s\n", signature.ToHumanAddress(addr))
	fmt.Printf("Valid with pubkey: %t\n", signature.Verify(pub, sig, msg))
	fmt.Printf("Valid with addr: %t\n", signature.VerifyWithAddress(addr, sig, msg))

	wseed := "within the sound of silence"
	wpriv, wpub, waddr, _ := signature.GenerateFromSeed(wseed)
	wmsg := []byte("foobar")
	wsig, _ := signature.Sign(wpriv, wmsg)
	wsig[0], wsig[1] = wsig[1], wsig[0]
	fmt.Printf("Valid with wrong pubkey: %t\n", signature.Verify(wpub, sig, msg))
	fmt.Printf("Valid with wrong addr: %t\n", signature.VerifyWithAddress(waddr, sig, msg))
	fmt.Printf("Valid with wrong sig: %t\n", signature.VerifyWithAddress(addr, wsig, msg))
	fmt.Printf("Valid with wrong msg: %t\n", signature.VerifyWithAddress(addr, sig, wmsg))

	// Output:
	// Private key: 0xfa98dc283ee4e866809227266e23c63048af7103d43fd65576a67269f8299f21
	// Public key: 0x01eace460e3fe0e07725f2e01ce2ee7e6fcb21078d8933841a5df5d980d7126bd7
	// Address: 0x0066423b3664eb7d98de3f77ceb32e49e25b4317787535351b
	// Human Address: 1AKhGQ5vTbMdSPDLYaueVAtgaNwS6mYkK4
	// Valid with pubkey: true
	// Valid with addr: true
	// Valid with wrong pubkey: false
	// Valid with wrong addr: false
	// Valid with wrong sig: false
	// Valid with wrong msg: false
}

func BenchmarkVerifyWithAddress(b *testing.B) {
	priv, _, addr, _ := signature.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := signature.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		signature.VerifyWithAddress(addr, sig, msg)
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
