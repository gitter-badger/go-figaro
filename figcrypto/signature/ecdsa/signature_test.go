package ecdsa_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto/signature/common"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/ecdsa"
)

func ExampleVerify() {
	seed := "hello darkness my old friend"
	priv, pub, addr, err := ecdsa.GenerateKeyFromSeed(seed)
	if err != nil {
		log.Panic(err)
	}
	msg := []byte("hello world")
	sig, _ := ecdsa.Sign(priv, msg)
	fmt.Printf("Private key: %#x\n", priv)
	fmt.Printf("Public key: %#x\n", pub)
	fmt.Printf("Address: %#x\n", addr)
	fmt.Printf("Human Address: %s\n", common.ToHumanAddress(addr))
	fmt.Printf("Valid addr: %t\n", ecdsa.VerifyAddress(pub, addr))
	fmt.Printf("Valid with pub: %t\n", ecdsa.Verify(pub, sig, msg))
	fmt.Printf("Valid with addr: %t\n", ecdsa.VerifyWithAddress(addr, sig, msg))

	wseed := "within the sound of silence"
	wpriv, wpub, waddr, _ := ecdsa.GenerateKeyFromSeed(wseed)
	wmsg := []byte("foobar")
	wsig, _ := ecdsa.Sign(wpriv, wmsg)
	wsig[0], wsig[1] = wsig[1], wsig[0]
	fmt.Printf("Valid with wrong pub: %t\n", ecdsa.Verify(wpub, sig, msg))
	fmt.Printf("Valid with wrong sig: %t\n", ecdsa.Verify(pub, wsig, msg))
	fmt.Printf("Valid with wrong msg: %t\n", ecdsa.Verify(pub, sig, wmsg))
	fmt.Printf("Valid with wrong addr: %t\n", ecdsa.VerifyWithAddress(waddr, sig, wmsg))

	// Output:
	// Private key: 0xd7867042dd1e3019969b77aebb609aa4ab0084814048d254748bffff6c434e0f
	// Public key: 0x006260119dab61eb2695160ddcc25da40e822e124bebce976d132d557c4ce7a0d7
	// Address: 0x0050843897153e22418c21155ff268151ab650d5276edd5e72
	// Human Address: 18LjTE8fQngBDpNyHjdGgbWgoLciQZqKKw
	// Valid addr: true
	// Valid with pub: true
	// Valid with addr: true
	// Valid with wrong pub: false
	// Valid with wrong sig: false
	// Valid with wrong msg: false
	// Valid with wrong addr: false
}

func BenchmarkAddressFromPublicKey(b *testing.B) {
	_, pub, _, _ := ecdsa.GenerateKey(nil)
	for i := 0; i < b.N; i++ {
		common.AddressFromPublicKey(pub)
	}
}

func BenchmarkAddressFromPublicKey_with_verify(b *testing.B) {
	_, pub, addr, _ := ecdsa.GenerateKey(nil)
	for i := 0; i < b.N; i++ {
		test := common.AddressFromPublicKey(pub)
		bytes.Equal(addr, test)
	}
}

func BenchmarkVerify(b *testing.B) {
	priv, pub, _, _ := ecdsa.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := ecdsa.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		ecdsa.Verify(pub, sig, msg)
	}
}

func BenchmarkVerifyWithAddress(b *testing.B) {
	priv, _, addr, _ := ecdsa.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := ecdsa.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		ecdsa.VerifyWithAddress(addr, sig, msg)
	}
}
