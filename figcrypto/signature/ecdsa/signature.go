// Package ecdsa provides cryptographic functions
package ecdsa

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"io"
	"math/big"

	"github.com/figaro-tech/go-figaro/figcrypto/hash"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/common"
)

// Constants
const (
	AddressSize    = common.AddressSize
	SignatureSize  = 65
	PublicKeySize  = 33
	PrivateKeySize = 32
	lPublicKeySize = 64
	cofactor       = 1
)

// CompactEncodePublicKey64 compacts encodes a 64-byte public key into
// a 33-byte public key that is recoverable. All functions in this package
// use the 33-byte public key. These methods are provided as a convenience
// for functionality that requires the full 64-byte public key.
func CompactEncodePublicKey64(pubkey64 []byte) (pubkey33 []byte, err error) {
	if len(pubkey64) != lPublicKeySize {
		err = common.ErrInvalidKey
		return
	}
	key := bytesToPubKey(pubkey64)
	b := makeLen(key.X.Bytes(), 32)
	pubkey33 = make([]byte, 1, len(b)+1)
	pubkey33[0] = byte(key.Y.Bit(0))
	pubkey33 = append(pubkey33, b...)
	return
}

// CompactDecodePublicKey33 compact decodes a 33-byte public key into
// a 64-byte public key. All functions in this package
// use the 33-byte public key. These methods are provided as a convenience
// for functionality that requires the full 64-byte public key.
func CompactDecodePublicKey33(pubkey33 []byte) (pubkey64 []byte, err error) {
	if len(pubkey33) != PublicKeySize {
		err = common.ErrInvalidKey
		return
	}
	sign, pubkey32 := uint(pubkey33[0]), pubkey33[1:]
	x := new(big.Int).SetBytes(pubkey32)
	y := decompressPoint(elliptic.P256().Params(), x, sign)
	pubkey64 = append(makeLen(x.Bytes(), 32), makeLen(y.Bytes(), 32)...)
	return
}

// GenerateKey generates an public/private key pair, along with an address,
// that can be used to sign and verify messages.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rando io.Reader) (privkey, pubkey, address []byte, err error) {
	if rando == nil {
		rando = rand.Reader
	}
	var key *ecdsa.PrivateKey
	key, err = ecdsa.GenerateKey(elliptic.P256(), rando)
	if err != nil {
		return
	}
	privkey = key.D.Bytes()
	pubkey, err = CompactEncodePublicKey64(pubKeyToBytes(&key.PublicKey))
	address = common.AddressFromPublicKey(pubkey)
	return
}

// GenerateKeyFromSeed creates an public/private key pair, along with an address,
// that can be used to verify messages, from a seed string or strings.
func GenerateKeyFromSeed(seed ...string) (privkey, publickey, address []byte, err error) {
	var h []byte
	for _, s := range seed {
		h = hash.Hash512(h, []byte(s))
	}
	rando := bytes.NewReader(h)
	return GenerateKey(rando)
}

// RecoverFromPrivateKey gets an address which can verify transactions from a private key.
func RecoverFromPrivateKey(privkey []byte) (pubkey, address []byte, err error) {
	if len(privkey) != PrivateKeySize {
		err = common.ErrInvalidKey
		return
	}
	c := elliptic.P256()
	x, y := c.ScalarBaseMult(privkey)
	b := append(makeLen(x.Bytes(), 32), makeLen(y.Bytes(), 32)...)
	pubkey, err = CompactEncodePublicKey64(b)
	address = common.AddressFromPublicKey(pubkey)
	return
}

// SignerFromPrivateKey returns a `crypto.Signer` from private key bytes.
func SignerFromPrivateKey(privkey []byte) (crypto.Signer, error) {
	if len(privkey) != PrivateKeySize {
		return nil, common.ErrInvalidKey
	}
	c := elliptic.P256()
	d := new(big.Int).SetBytes(privkey)
	x, y := c.ScalarBaseMult(d.Bytes())
	return &ecdsa.PrivateKey{
		D:         d,
		PublicKey: ecdsa.PublicKey{Curve: c, X: x, Y: y},
	}, nil
}

// Sign signs the message using the private key
func Sign(privkey, message []byte) (signature []byte, err error) {
	var key crypto.Signer
	key, err = SignerFromPrivateKey(privkey)
	if err != nil {
		return
	}
	return signWithKey(key.(*ecdsa.PrivateKey), message)
}

// Verify verifies that a message was signed by the owner of the compact encoded public key.
func Verify(pubkey, signature, message []byte) bool {
	if len(pubkey) != PublicKeySize {
		return false
	}
	b, err := CompactDecodePublicKey33(pubkey)
	if err != nil {
		return false
	}
	key := bytesToPubKey(b)
	return verifyWithKey(key, signature, message)
}

// VerifyAddress verifies that the public key is valid for the address.
func VerifyAddress(publickey, address []byte) bool {
	if len(publickey) != PublicKeySize || len(address) != AddressSize {
		return false
	}
	derived := common.AddressFromPublicKey(publickey)
	return bytes.Equal(address, derived)
}

// VerifyWithAddress verifies that a message was signed by the owner of the address.
func VerifyWithAddress(address, signature, message []byte) bool {
	if len(address) != AddressSize {
		return false
	}
	if len(signature) != SignatureSize {
		return false
	}
	v := uint(signature[0] - 27)
	sig := signature[1:]
	rbytes, sbytes := sig[:len(sig)/2], sig[len(sig)/2:]
	r, s := new(big.Int).SetBytes(rbytes), new(big.Int).SetBytes(sbytes)
	a := recoverAddress(r, s, v, message)
	if len(a) == 0 {
		return false
	}
	return bytes.Equal(a, address)
}

func signWithKey(privkey *ecdsa.PrivateKey, message []byte) (signature []byte, err error) {
	h := hash.Hash256(message)
	var r, s *big.Int
	r, s, err = ecdsa.Sign(rand.Reader, privkey, h)
	if err != nil {
		return
	}
	curve := privkey.Curve.Params()
	curvelen := (curve.BitSize + 7) / 8
	var pub *ecdsa.PublicKey
	for i := 0; i < (cofactor+1)*2; i++ {
		pub, err = recover(r, s, h, uint(i), true)
		if err == nil && pub.X.Cmp(privkey.X) == 0 && pub.Y.Cmp(privkey.Y) == 0 {
			signature = make([]byte, 1, 2*curve.BitSize/8+1)
			signature[0] = byte(i + 27)
			signature = append(signature, makeLen(r.Bytes(), curvelen)...)
			signature = append(signature, makeLen(s.Bytes(), curvelen)...)
			return
		}
	}
	err = errors.New("no valid solution for pubkey found")
	return
}

func verifyWithKey(pub *ecdsa.PublicKey, signature, message []byte) bool {
	if len(signature) != SignatureSize {
		return false
	}
	sig := signature[1:]
	rbytes, sbytes := sig[:len(sig)/2], sig[len(sig)/2:]
	h := hash.Hash256(message)
	return ecdsa.Verify(pub, h, new(big.Int).SetBytes(rbytes), new(big.Int).SetBytes(sbytes))
}

func bytesToPubKey(pubkey64 []byte) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(pubkey64[:32]),
		Y:     new(big.Int).SetBytes(pubkey64[32:]),
	}
}

func pubKeyToBytes(pubkey *ecdsa.PublicKey) (pubkey64 []byte) {
	x := makeLen(pubkey.X.Bytes(), 32)
	y := makeLen(pubkey.Y.Bytes(), 32)
	return append(x, y...)
}

func recoverAddress(r, s *big.Int, v uint, message []byte) (address []byte) {
	h := hash.Hash256(message)
	key, err := recover(r, s, h, v, true)
	if err != nil {
		return
	}
	b, _ := CompactEncodePublicKey64(pubKeyToBytes(key))
	address = common.AddressFromPublicKey(b)
	return
}

// Based on SEC 1 Ver 2.0 Sec. 4.1.6
// See http://www.secg.org/sec1-v2.pdf
func recover(r, s *big.Int, msg []byte, iter uint, doChecks bool) (pub *ecdsa.PublicKey, err error) {
	c := elliptic.P256()
	// Step 1.1
	rx := new(big.Int).Mul(c.Params().N, new(big.Int).SetInt64(int64(iter/2)))
	rx.Add(rx, r)
	if rx.Cmp(c.Params().P) != -1 {
		err = errors.New("calculated X is larger than curve P")
		return
	}
	// Steps 1.2 and 1.3
	ry := decompressPoint(c.Params(), rx, iter&1)
	// Step 1.4
	if doChecks {
		nrx, nry := c.ScalarMult(rx, ry, c.Params().N.Bytes())
		if nrx.Sign() != 0 || nry.Sign() != 0 {
			err = errors.New("n*R does not equal the point at infinity")
			return
		}
	}
	// Step 1.5
	e := hashToInt(msg, c)
	// Step 1.6.1:
	invr := new(big.Int).ModInverse(r, c.Params().N)
	invrS := new(big.Int).Mul(invr, s)
	invrS.Mod(invrS, c.Params().N)
	srx, sry := c.ScalarMult(rx, ry, invrS.Bytes())

	e.Neg(e)
	e.Mod(e, c.Params().N)
	e.Mul(e, invr)
	e.Mod(e, c.Params().N)
	minuseGx, minuseGy := c.ScalarBaseMult(e.Bytes())

	qx, qy := c.Add(srx, sry, minuseGx, minuseGy)

	return &ecdsa.PublicKey{
		Curve: c,
		X:     qx,
		Y:     qy,
	}, nil
}

func decompressPoint(curve *elliptic.CurveParams, x *big.Int, sign uint) (y *big.Int) {
	// y² = x³ - 3x + b
	y = new(big.Int).Mul(x, x)
	y.Mul(y, x)
	threeX := new(big.Int).Lsh(x, 1)
	threeX.Add(threeX, x)
	y.Sub(y, threeX)
	y.Add(y, curve.B)
	y.Mod(y, curve.P)
	y.ModSqrt(y, curve.P)
	if y == nil {
		panic(common.ErrInvalidKey)
	}
	if y.Bit(0) != sign&1 {
		y.Neg(y)
		y.Mod(y, curve.P)
	}
	return
}

// This is borrowed from crypto/ecdsa.
func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

var one = new(big.Int).SetInt64(1)

// This is largely borrowed from crypto/ecdsa.
// fieldElement returns an element of the field underlying the given
// curve using the procedure given in [NSA] A.2.1.
func fieldElement(c elliptic.Curve, seed []byte) (k *big.Int) {
	params := c.Params()
	b := make([]byte, params.BitSize/8+8)
	copy(b, seed)
	k = new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)
	return
}

func makeLen(b []byte, l int) (result []byte) {
	if len(b) > l {
		panic("bytes exceed len")
	}
	if len(b) < l {
		result = append(result, make([]byte, l-len(b))...)
	}
	result = append(result, b...)
	return
}
