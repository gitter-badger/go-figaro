package fastsig_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto/signature/common"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/fastsig"
)

func ExampleVerify() {
	seed := "hello darkness my old friend"
	priv, pub, addr, err := fastsig.GenerateKeyFromSeed(seed)
	if err != nil {
		log.Panic(err)
	}
	msg := []byte("hello world")
	sig, err := fastsig.Sign(priv, msg)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Private key: %#x\n", priv)
	fmt.Printf("Public key: %#x\n", pub)
	fmt.Printf("Address: %#x\n", addr)
	fmt.Printf("Human Address: %s\n", common.ToHumanAddress(addr))
	fmt.Printf("Valid addr: %t\n", fastsig.VerifyAddress(pub, addr))
	fmt.Printf("Valid with pub: %t\n", fastsig.Verify(pub, sig, msg))
	fmt.Printf("Valid with addr: %t\n", fastsig.VerifyWithAddress(addr, sig, msg))

	wseed := "within the sound of silence"
	wpriv, wpub, waddr, _ := fastsig.GenerateKeyFromSeed(wseed)
	wmsg := []byte("foobar")
	wsig, _ := fastsig.Sign(wpriv, wmsg)
	wsig[0], wsig[1] = wsig[1], wsig[0]
	fmt.Printf("Valid with wrong pub: %t\n", fastsig.Verify(wpub, sig, msg))
	fmt.Printf("Valid with wrong sig: %t\n", fastsig.Verify(pub, wsig, msg))
	fmt.Printf("Valid with wrong msg: %t\n", fastsig.Verify(pub, sig, wmsg))
	fmt.Printf("Valid with wrong addr: %t\n", fastsig.VerifyWithAddress(waddr, sig, wmsg))

	// Output:
	// Private key: 0xb56df07a2a7a38fec23f03ab4f3339a33cedd196bfaecde67d47ae4dbfe6335f
	// Public key: 0x029294ee1a518a34c6552a596af899778e677456de0362df8a9762e3142f30a563
	// Address: 0x00af5c7673dbb07b04e0af246fa5ed6e5d0ea5f0dc0db72766
	// Human Address: 1GzE1zX9KGRGnijvXEzxKvj9AoVphYnpqs
	// Valid addr: true
	// Valid with pub: true
	// Valid with addr: true
	// Valid with wrong pub: false
	// Valid with wrong sig: false
	// Valid with wrong msg: false
	// Valid with wrong addr: false

}

func BenchmarkAddressFromPublicKey(b *testing.B) {
	_, pub, _, _ := fastsig.GenerateKey(nil)
	for i := 0; i < b.N; i++ {
		common.AddressFromPublicKey(pub)
	}
}

func BenchmarkAddressFromPublicKey_with_verify(b *testing.B) {
	_, pub, addr, _ := fastsig.GenerateKey(nil)
	for i := 0; i < b.N; i++ {
		test := common.AddressFromPublicKey(pub)
		bytes.Equal(addr, test)
	}
}

func BenchmarkVerify(b *testing.B) {
	priv, pub, _, _ := fastsig.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := fastsig.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		fastsig.Verify(pub, sig, msg)
	}
}

func BenchmarkVerifyWithAddress(b *testing.B) {
	priv, _, addr, _ := fastsig.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := fastsig.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		fastsig.VerifyWithAddress(addr, sig, msg)
	}
}
