// Package govsig provides cryptographic functions. It implements an implementation of
// of recoverable ECDSA signatures using Curve P-256.
package govsig

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"io"
	"math/big"

	"github.com/figaro-tech/go-fig-crypto/hasher"
	"github.com/figaro-tech/go-fig-crypto/signature/common"
)

// Constants
const (
	AddressSize               = common.AddressSize
	SignatureSize             = 65
	PublicKeySize             = 33
	DecompressedPublicKeySize = 65
	PrivateKeySize            = 32
	cofactor                  = 1
)

// ToHumanAddress converts a binary address to a Base58 encoded "human readable" address
func ToHumanAddress(binaddr []byte) string {
	return common.ToHumanAddress(binaddr)
}

// ToBinaryAddress converts a Base58 encoded "human readable" address to a binary address
func ToBinaryAddress(humaddr string) []byte {
	return common.ToBinaryAddress(humaddr)
}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkey(pubkey []byte) ([]byte, error) {
	if len(pubkey) != PublicKeySize {
		return nil, common.ErrInvalidKey
	}
	x := new(big.Int).SetBytes(pubkey[1:])
	y := decompressPoint(elliptic.P256().Params(), x, uint(pubkey[0]&1))
	if x == nil {
		return nil, common.ErrInvalidKey
	}
	return marshal(x, y), nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkey(pubkey []byte) ([]byte, error) {
	if len(pubkey) != DecompressedPublicKeySize {
		return nil, common.ErrInvalidKey
	}
	x, y := unmarshal(pubkey)
	cpk := make([]byte, 33)
	if y.Bit(0)&1 == 0 {
		cpk[0] = 02
	} else {
		cpk[0] = 03
	}
	copy(cpk[1:], x.Bytes())
	return cpk, nil
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
	pubkey, _ = CompressPubkey(marshal(key.X, key.Y))
	address = common.AddressFromPublicKey(pubkey)
	return
}

// GenerateKeyFromSeed creates an public/private key pair, along with an address,
// that can be used to verify messages, from a seed string or strings.
func GenerateKeyFromSeed(seed ...string) (privkey, publickey, address []byte, err error) {
	var h []byte
	for _, s := range seed {
		h = hasher.Hash512(h, []byte(s))
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
	pubkey, _ = CompressPubkey(marshal(x, y))
	address = common.AddressFromPublicKey(pubkey)
	return
}

// RecoverFromSignature gets an address which can verify transactions from a signature and message.
func RecoverFromSignature(signature, message []byte) (pubkey, address []byte, err error) {
	if len(signature) != SignatureSize {
		return nil, nil, common.ErrInvalidSignature
	}
	h := hasher.Hash256(message)
	r := signature[:(SignatureSize-1)/2]
	s := signature[(SignatureSize-1)/2 : SignatureSize-1]
	v := signature[SignatureSize-1]
	key, err := recover(new(big.Int).SetBytes(r), new(big.Int).SetBytes(s), h, uint(v), false)
	if err != nil {
		return nil, nil, err
	}
	pubkey, _ = CompressPubkey(marshal(key.X, key.Y))
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
	if len(privkey) != PrivateKeySize {
		return nil, common.ErrInvalidKey
	}
	key, err := SignerFromPrivateKey(privkey)
	if err != nil {
		return nil, err
	}
	return sign(key.(*ecdsa.PrivateKey), message)
}

// Verify verifies that a message was signed by the owner of the compact encoded public key.
func Verify(pubkey, signature, message []byte) bool {
	if len(signature) != SignatureSize {
		return false
	}
	if len(pubkey) != PublicKeySize && len(pubkey) != DecompressedPublicKeySize {
		return false
	}
	var x, y *big.Int
	if len(pubkey) == PublicKeySize {
		key, err := DecompressPubkey(pubkey)
		if err != nil {
			return false
		}
		x, y = elliptic.Unmarshal(elliptic.P256(), key)
	} else {
		x, y = elliptic.Unmarshal(elliptic.P256(), pubkey)
	}
	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return verify(key, signature, message)
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
	_, rec, err := RecoverFromSignature(signature, message)
	if err != nil {
		return false
	}
	return bytes.Equal(address, rec)
}

func sign(privkey *ecdsa.PrivateKey, message []byte) (signature []byte, err error) {
	h := hasher.Hash256(message)
	var r, s *big.Int
	r, s, err = ecdsa.Sign(rand.Reader, privkey, h)
	if err != nil {
		return
	}
	var pub *ecdsa.PublicKey
	for i := 0; i < (cofactor+1)*2; i++ {
		pub, err = recover(r, s, h, uint(i), true)
		if err == nil && pub.X.Cmp(privkey.X) == 0 && pub.Y.Cmp(privkey.Y) == 0 {
			signature = make([]byte, 0, 65)
			signature = append(signature, r.Bytes()...)
			signature = append(signature, s.Bytes()...)
			signature = append(signature, byte(i))
			return
		}
	}
	err = errors.New("no valid solution for pubkey found")
	return
}

func verify(pub *ecdsa.PublicKey, signature, message []byte) bool {
	if len(signature) != SignatureSize {
		return false
	}
	sig := signature[:SignatureSize-1]
	rbytes, sbytes := sig[:len(sig)/2], sig[len(sig)/2:]
	h := hasher.Hash256(message)
	return ecdsa.Verify(pub, h, new(big.Int).SetBytes(rbytes), new(big.Int).SetBytes(sbytes))
}

// Marshal converts a point into the form specified in section 4.3.6 of ANSI
// X9.62.
func marshal(x, y *big.Int) []byte {
	byteLen := (elliptic.P256().Params().BitSize + 7) >> 3
	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // uncompressed point flag
	readBits(x, ret[1:1+byteLen])
	readBits(y, ret[1+byteLen:])
	return ret
}

// Unmarshal converts a point, serialised by Marshal, into an x, y pair. On
// error, x = nil.
func unmarshal(data []byte) (x, y *big.Int) {
	byteLen := (elliptic.P256().Params().BitSize + 7) >> 3
	if len(data) != 1+2*byteLen {
		return
	}
	if data[0] != 4 { // uncompressed form
		return
	}
	x = new(big.Int).SetBytes(data[1 : 1+byteLen])
	y = new(big.Int).SetBytes(data[1+byteLen:])
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

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

func readBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}
