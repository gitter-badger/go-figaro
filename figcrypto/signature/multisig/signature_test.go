package multisig_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto/signature/common"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/multisig"
)

func ExampleVerify() {
	seed := "hello darkness my old friend"
	priv, pub, addr, err := multisig.GenerateKeyFromSeed(seed)
	if err != nil {
		log.Panic(err)
	}
	msg := []byte("hello world")
	sig, _ := multisig.Sign(priv, msg)
	fmt.Printf("Private key: %#x\n", priv)
	fmt.Printf("Public key: %#x\n", pub)
	fmt.Printf("Address: %#x\n", addr)
	fmt.Printf("Human Address: %s\n", common.ToHumanAddress(addr))
	fmt.Printf("Valid addr: %t\n", multisig.VerifyAddress(pub, addr))
	fmt.Printf("Valid with pub: %t\n", multisig.Verify(pub, sig, msg))

	wseed := "within the sound of silence"
	wpriv, wpub, _, _ := multisig.GenerateKeyFromSeed(wseed)
	wmsg := []byte("foobar")
	wsig, _ := multisig.Sign(wpriv, wmsg)
	wsig[0], wsig[1] = wsig[1], wsig[0]
	fmt.Printf("Valid with wrong pub: %t\n", multisig.Verify(wpub, sig, msg))
	fmt.Printf("Valid with wrong sig: %t\n", multisig.Verify(addr, wsig, msg))
	fmt.Printf("Valid with wrong msg: %t\n", multisig.Verify(addr, sig, wmsg))

	// Output:
	// Private key: 0x2b43891c22187fc9b56df07a2a7a38fe8b44904db30800f0eb44d1e351b57b3dd0ff51ca5130b7dc3b795067e3169c3e2090e75284e32e8a050ae91b4fbb7611
	// Public key: 0xd0ff51ca5130b7dc3b795067e3169c3e2090e75284e32e8a050ae91b4fbb7611
	// Address: 0x005985122bf4dd299d95ccb347ca8f87ecb5c1741090e1cc6b
	// Human Address: 19ALZ1QFxsuhbUFcuoQg9TSeGpY9fu4aX4
	// Valid addr: true
	// Valid with pub: true
	// Valid with wrong pub: false
	// Valid with wrong sig: false
	// Valid with wrong msg: false
}

func BenchmarkVerifyWithAddress(b *testing.B) {
	priv, pubkey, _, _ := multisig.GenerateKey(nil)
	msg := []byte("hello world")
	sig, _ := multisig.Sign(priv, msg)
	for i := 0; i < b.N; i++ {
		multisig.Verify(pubkey, sig, msg)
	}
}
