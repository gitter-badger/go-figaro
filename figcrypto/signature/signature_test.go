package signature_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto/signature"
)

func ExampleVerify() {
	priv, pub, addr, _ := signature.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := signature.Sign(priv, msg)
	fullpub := &signature.KeyFromBytes(priv).PublicKey
	log.Printf("Private key: %#x\n", priv)
	log.Printf("Public key: %#x\n", pub)
	log.Printf("Address: %s\n", addr)
	fmt.Printf("Valid identity: %t\n", signature.Identify(sig, msg) == addr)
	fmt.Printf("Valid with pubkey: %t\n", signature.Verify(pub, sig, msg))
	fmt.Printf("Valid with full pubkey: %t\n", signature.VerifyWithKey(fullpub, sig, msg))
	fmt.Printf("Valid with addr: %t\n", signature.VerifyWithAddr(addr, sig, msg))

	wpriv, wpub, waddr, _ := signature.GenerateKey(nil)
	wsig, _ := signature.Sign(wpriv, msg)
	wfullpub := &signature.KeyFromBytes(wpriv).PublicKey
	fmt.Printf("Invalid with wrong sig: %t\n", signature.VerifyWithKey(fullpub, wsig, msg))
	fmt.Printf("Invalid with wrong addr: %t\n", signature.VerifyWithAddr(waddr, sig, msg))
	fmt.Printf("Invalid with wrong pubkey: %t\n", signature.Verify(wpub, sig, msg))
	fmt.Printf("Invalid with alid with full pubkey: %t\n", signature.VerifyWithKey(wfullpub, sig, msg))

	// Output:
	// Valid identity: true
	// Valid with pubkey: true
	// Valid with full pubkey: true
	// Valid with addr: true
	// Invalid with wrong sig: false
	// Invalid with wrong addr: false
	// Invalid with wrong pubkey: false
	// Invalid with alid with full pubkey: false
}

func ExampleRecoverPublicFromPrivate() {
	priv, pub, addr, _ := signature.GenerateKey(nil)

	rpub, raddr, _ := signature.RecoverPublicFromPrivate(priv)
	fmt.Printf("Pub is same: %t\n", bytes.Equal(pub, rpub))
	fmt.Printf("Addr is same: %t\n", addr == raddr)

	// Output:
	// Pub is same: true
	// Addr is same: true
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

func BenchmarkVerifyWithFullKey(b *testing.B) {
	priv, _, _, _ := signature.GenerateKey(nil)
	pub := &signature.KeyFromBytes(priv).PublicKey
	msg := []byte("hello world")
	sig, _ := signature.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		signature.VerifyWithKey(pub, sig, msg)
	}
}
