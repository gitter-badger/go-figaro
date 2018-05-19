// Package signature provides cryptographic functions
package signature

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"github.com/figaro-tech/go-figaro/figcrypto/hash"
)

const (
	cofactor = 1
)

// Errors used in package
var (
	version       = []byte{0x00} // Using different versions for different nets means addresses won't cross-validate
	ErrInvalidKey = errors.New("figcrypto signature: invalid public or private key")
)

func init() {
	if v, ok := os.LookupEnv("ADDRESS_VERSION_CODE"); ok {
		b, _ := hex.DecodeString(v)
		if len(b) == 1 {
			version = b
		}
	}
}

// ToHumanAddress converts a binary address to a Base58 encoded "human readable" address
func ToHumanAddress(binaddr []byte) string {
	return base58.Encode(binaddr)
}

// ToBinaryAddress converts a Base58 encoded "human readable" address to a binary address
func ToBinaryAddress(humaddr string) []byte {
	return base58.Decode(humaddr)
}

// PublicKeyFromPrivateKey recovers the public key from a private key.
func PublicKeyFromPrivateKey(privkey32 []byte) (pubkey33 []byte, err error) {
	if len(privkey32) != 32 {
		err = ErrInvalidKey
		return
	}
	c := elliptic.P256()
	x, y := c.ScalarBaseMult(privkey32)
	b := append(makeLen(x.Bytes(), 32), makeLen(y.Bytes(), 32)...)
	pubkey33, err = CompactEncodePublicKey64(b)
	return
}

// AddressFromPublicKey creates an human-readable address from a public key, following the Bitcoin address protocol.
func AddressFromPublicKey(pubkey33 []byte) (address []byte, err error) {
	if len(pubkey33) != 33 {
		err = ErrInvalidKey
		return
	}
	h := hash.Hash256(pubkey33)
	h = hash.Hash160(h)
	address = append(version, h...)
	h2 := hash.Hash256(address)
	h2 = hash.Hash256(h2)
	checksum := h2[:4]
	address = append(address, checksum...)
	return
}

// SignerFromPrivateKey hydrates a *ecsda.PrivateKey from private key bytes for use
// when a Go `crypto.Signer` is needed.
func SignerFromPrivateKey(privkey32 []byte) (*ecdsa.PrivateKey, error) {
	if len(privkey32) != 32 {
		return nil, ErrInvalidKey
	}
	c := elliptic.P256()
	d := new(big.Int).SetBytes(privkey32)
	x, y := c.ScalarBaseMult(d.Bytes())
	return &ecdsa.PrivateKey{
		D:         d,
		PublicKey: ecdsa.PublicKey{Curve: c, X: x, Y: y},
	}, nil
}

// GenerateKey generates a public/private key pair that can be used to sign and verify messages,
// along with a corresponding address. if rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rando io.Reader) (privkey32, pubkey33, address []byte, err error) {
	if rando == nil {
		rando = rand.Reader
	}
	var key *ecdsa.PrivateKey
	key, err = ecdsa.GenerateKey(elliptic.P256(), rando)
	if err != nil {
		return
	}
	privkey32 = key.D.Bytes()
	pubkey33, err = CompactEncodePublicKey64(pubKeyToBytes(&key.PublicKey))
	address, _ = AddressFromPublicKey(pubkey33)
	return
}

// GenerateFromSeed creates a private/public key pair that can be used to verify messages,
// along with a corresponding address, from a seed hash.
func GenerateFromSeed(seed ...string) (privkey32, pubkey33, address []byte, err error) {
	c := elliptic.P256()
	var h []byte
	for _, s := range seed {
		h = hash.Hash256(h, []byte(s))
	}
	privkey32 = fieldElement(c, h).Bytes()
	x, y := c.ScalarBaseMult(privkey32)
	b := append(makeLen(x.Bytes(), 32), makeLen(y.Bytes(), 32)...)
	pubkey33, err = CompactEncodePublicKey64(b)
	if err != nil {
		return
	}
	address, err = AddressFromPublicKey(pubkey33)
	return
}

// Sign signs the message using the private key
func Sign(privkey32, message []byte) (signature []byte, err error) {
	var key *ecdsa.PrivateKey
	key, err = SignerFromPrivateKey(privkey32)
	if err != nil {
		return
	}
	return SignWithKey(key, message)
}

// SignWithKey signs the message using the private key
func SignWithKey(privkey *ecdsa.PrivateKey, message []byte) (signature []byte, err error) {
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
			signature[0] = byte(i)
			signature = append(signature, makeLen(r.Bytes(), curvelen)...)
			signature = append(signature, makeLen(s.Bytes(), curvelen)...)
			return
		}
	}
	err = errors.New("no valid solution for pubkey found")
	return
}

// Identify determines the address that signed a message
func Identify(signature, message []byte) (address []byte) {
	if len(signature) != 65 {
		return
	}
	v := uint(signature[0])
	sig := signature[1:]
	rbytes, sbytes := sig[:len(sig)/2], sig[len(sig)/2:]
	r, s := new(big.Int).SetBytes(rbytes), new(big.Int).SetBytes(sbytes)
	h := hash.Hash256(message)
	key, err := recover(r, s, h, v, false)
	if err != nil {
		return
	}
	b, _ := CompactEncodePublicKey64(pubKeyToBytes(key))
	address, _ = AddressFromPublicKey(b)
	return
}

// Verify verifies that a message was signed by the owner of the compact encoded public key
func Verify(pubkey33, signature, message []byte) bool {
	if len(pubkey33) != 33 {
		return false
	}
	b, err := CompactDecodePublicKey33(pubkey33)
	if err != nil {
		return false
	}
	key := bytesToPubKey(b)
	return VerifyWithKey(key, signature, message)
}

// VerifyWithAddress verifies that a message was signed by the owner of the address
func VerifyWithAddress(address, signature, message []byte) bool {
	a := Identify(signature, message)
	if len(a) == 0 {
		return false
	}
	return bytes.Equal(a, address)
}

// VerifyWithKey verifies that a message was signed by the owner of the public key
func VerifyWithKey(pub *ecdsa.PublicKey, signature, message []byte) bool {
	if len(signature) != 65 {
		return false
	}
	sig := signature[1:]
	rbytes, sbytes := sig[:len(sig)/2], sig[len(sig)/2:]
	h := hash.Hash256(message)
	return ecdsa.Verify(pub, h, new(big.Int).SetBytes(rbytes), new(big.Int).SetBytes(sbytes))
}

// CompactEncodePublicKey64 compacts encodes a 64-byte public key into
// a 33-byte public key that is recoverable. All functions in this package
// use the 33-byte public key. These methods are provided as a convenience
// for functionality that requires the full 64-byte public key.
func CompactEncodePublicKey64(pubkey64 []byte) (pubkey33 []byte, err error) {
	if len(pubkey64) != 64 {
		err = ErrInvalidKey
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
	if len(pubkey33) != 33 {
		err = ErrInvalidKey
		return
	}
	sign, pubkey32 := uint(pubkey33[0]), pubkey33[1:]
	x := new(big.Int).SetBytes(pubkey32)
	y := decompressPoint(elliptic.P256().Params(), x, sign)
	pubkey64 = append(makeLen(x.Bytes(), 32), makeLen(y.Bytes(), 32)...)
	return
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

// Based on SEC 1 Ver 2.0 Sec. 4.1.6
// See http://www.secg.org/sec1-v2.pdf
func recover(r, s *big.Int, msg []byte, iter uint, doChecks bool) (pub *ecdsa.PublicKey, err error) {
	c := elliptic.P256()
	// 1.1 x = (n * i) + r
	rx := new(big.Int).Mul(c.Params().N, new(big.Int).SetInt64(int64(iter/2)))
	rx.Add(rx, r)
	if rx.Cmp(c.Params().P) != -1 {
		err = errors.New("calculated X is larger than curve P")
		return
	}
	// (step 1.2 and 1.3). If we are on an odd
	// iteration then 1.6 will be done with -R, so we calculate the other
	// term when uncompressing the point.
	ry := decompressPoint(c.Params(), rx, iter&1)

	// 1.4 Check n*R is point at infinity
	if doChecks {
		nrx, nry := c.ScalarMult(rx, ry, c.Params().N.Bytes())
		if nrx.Sign() != 0 || nry.Sign() != 0 {
			err = errors.New("n*R does not equal the point at infinity")
			return
		}
	}

	// 1.5 calculate e from message using the same algorithm as ecdsa
	// signature calculation.
	e := hashToInt(msg, c)

	// Step 1.6.1:
	// We calculate the two terms sR and eG separately multiplied by the
	// inverse of r (from the signature). We then add them to calculate
	// Q = r^-1(sR-eG)
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
		panic(ErrInvalidKey)
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
