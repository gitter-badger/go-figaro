package govsig_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/figaro-tech/go-fig-crypto/signature/common"
	"github.com/figaro-tech/go-fig-crypto/signature/govsig"
)

func ExampleVerify() {
	seed := "hello darkness my old friend"
	priv, pub, addr, err := govsig.GenerateKeyFromSeed(seed)
	if err != nil {
		log.Panic(err)
	}
	msg := []byte("hello world")
	sig, _ := govsig.Sign(priv, msg)
	fmt.Printf("Private key: %#x\n", priv)
	fmt.Printf("Public key: %#x\n", pub)
	fmt.Printf("Address: %#x\n", addr)
	fmt.Printf("Human Address: %s\n", common.ToHumanAddress(addr))
	fmt.Printf("Valid addr: %t\n", govsig.VerifyAddress(pub, addr))
	fmt.Printf("Valid with pub: %t\n", govsig.Verify(pub, sig, msg))
	fmt.Printf("Valid with addr: %t\n", govsig.VerifyWithAddress(addr, sig, msg))

	wseed := "within the sound of silence"
	wpriv, wpub, waddr, _ := govsig.GenerateKeyFromSeed(wseed)
	wmsg := []byte("foobar")
	wsig, _ := govsig.Sign(wpriv, wmsg)
	wsig[0], wsig[1] = wsig[1], wsig[0]
	fmt.Printf("Valid with wrong pub: %t\n", govsig.Verify(wpub, sig, msg))
	fmt.Printf("Valid with wrong sig: %t\n", govsig.Verify(pub, wsig, msg))
	fmt.Printf("Valid with wrong msg: %t\n", govsig.Verify(pub, sig, wmsg))
	fmt.Printf("Valid with wrong addr: %t\n", govsig.VerifyWithAddress(waddr, sig, wmsg))

	// Output:
	// Private key: 0xd7867042dd1e3019969b77aebb609aa4ab0084814048d254748bffff6c434e0f
	// Public key: 0x026260119dab61eb2695160ddcc25da40e822e124bebce976d132d557c4ce7a0d7
	// Address: 0x0011c89829b40d7a9c50082db180181307c19dad0ca203d1bb
	// Human Address: 12d2nobiDn7DEHGcEbnc9Rkd2ztG2YN2Lr
	// Valid addr: true
	// Valid with pub: true
	// Valid with addr: true
	// Valid with wrong pub: false
	// Valid with wrong sig: false
	// Valid with wrong msg: false
	// Valid with wrong addr: false
}

func BenchmarkAddressFromPublicKey(b *testing.B) {
	_, pub, _, _ := govsig.GenerateKey(nil)
	for i := 0; i < b.N; i++ {
		common.AddressFromPublicKey(pub)
	}
}

func BenchmarkAddressFromPublicKey_with_verify(b *testing.B) {
	_, pub, addr, _ := govsig.GenerateKey(nil)
	for i := 0; i < b.N; i++ {
		test := common.AddressFromPublicKey(pub)
		bytes.Equal(addr, test)
	}
}

func BenchmarkVerify(b *testing.B) {
	priv, pub, _, _ := govsig.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := govsig.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		govsig.Verify(pub, sig, msg)
	}
}

func BenchmarkVerifyWithAddress(b *testing.B) {
	priv, _, addr, _ := govsig.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := govsig.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		govsig.VerifyWithAddress(addr, sig, msg)
	}
}
