package signature_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/figaro-tech/go-figaro/figcrypto/signature"
)

func ExampleVerify() {
	pubKey, privKey, err := signature.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}
	msg := []byte("hello world")
	sig, err := signature.Sign(privKey, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PubKey: %#x\n", pubKey)
	log.Printf("PrivKey: %#x\n", privKey)
	log.Printf("Msg: %#x\n", msg)
	log.Printf("Sig: %#x\n", sig)
	valid := signature.Verify(pubKey, msg, sig)
	fmt.Printf("Valid: %t\n", valid)

	wPub, wPriv, err := signature.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}

	sig, err = signature.Sign(wPriv, msg)
	if err != nil {
		log.Fatal(err)
	}
	valid = signature.Verify(pubKey, msg, sig)
	fmt.Printf("Valid wrong priv: %t\n", valid)

	sig, err = signature.Sign(privKey, msg)
	if err != nil {
		log.Fatal(err)
	}
	valid = signature.Verify(wPub, msg, sig)
	fmt.Printf("Valid wrong pub: %t\n", valid)

	// Output:
	// Valid: true
	// Valid wrong priv: false
	// Valid wrong pub: false
}

func BenchmarkVerify(b *testing.B) {
	pubKey, privKey, err := signature.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}
	msg := []byte("hello world")
	sig, err := signature.Sign(privKey, msg)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		signature.Verify(pubKey, msg, sig)
	}
}
