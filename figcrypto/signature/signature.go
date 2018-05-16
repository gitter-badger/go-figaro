// Package signature provides cryptographic functions
package signature

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"io"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"github.com/figaro-tech/go-figaro/figcrypto/hash"
)

// AddrFromBytes converts a public key to a Base58 encoded "human readable" address to reduce human error
func AddrFromBytes(pub []byte) string {
	return base58.Encode(pub)
}

// AddrToBytes converts a Base58 encoded "human readable" address to a public key
func AddrToBytes(s string) []byte {
	return base58.Decode(s)
}

// Creates an "address" from a public key, following the Bitcoin address protocol.
func generateAddress(pub []byte) string {
	h := hash.Hash256(pub)
	h = hash.Hash160(h)
	a := append([]byte{version}, h...)
	h2 := hash.Hash256(a)
	h2 = hash.Hash256(h2)
	checksum := h2[:4]
	a = append(a, checksum...)
	return AddrFromBytes(a)
}

const (
	version  = 0x00 // Using different versions for different nets means addresses won't cross-validate
	cofactor = 1
)

// Errors used in package
var (
	ErrInvalidKey = errors.New("figcrypto signature: invalid public or private key")
)

// GenerateKey generates a public/private key pair that can be used to sign and verify messages,
// along with a corresponding address. if rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rando io.Reader) (priv, pub []byte, addr string, err error) {
	if rando == nil {
		rando = rand.Reader
	}
	key, err := ecdsa.GenerateKey(elliptic.P256(), rando)
	if err != nil {
		return
	}
	priv = key.D.Bytes()
	pub = compactEncodePubKey(&key.PublicKey)
	addr = generateAddress(pub)
	return
}

// RecoverPublicFromPrivate generates a public key that can be used to verify messages from a private key,
// along with a corresponding address. If rand is nil, crypto/rand.Reader will be used.
func RecoverPublicFromPrivate(priv []byte) (pub []byte, addr string, err error) {
	c := elliptic.P256()
	x, y := c.ScalarBaseMult(priv)
	pub = compactEncodePubKey(&ecdsa.PublicKey{Curve: c, X: x, Y: y})
	addr = generateAddress(pub)
	return
}

// KeyFromBytes hydrates a *ecsda.PrivateKey from private key hash for convenience.
func KeyFromBytes(priv []byte) *ecdsa.PrivateKey {
	c := elliptic.P256()
	d := new(big.Int).SetBytes(priv)
	x, y := c.ScalarBaseMult(d.Bytes())
	return &ecdsa.PrivateKey{
		D:         d,
		PublicKey: ecdsa.PublicKey{Curve: c, X: x, Y: y},
	}
}

// Sign signs the message using the private key
func Sign(priv []byte, msg []byte) (sig []byte, err error) {
	key := KeyFromBytes(priv)
	return SignWithKey(key, msg)
}

// SignWithKey signs the message using the private key
func SignWithKey(priv *ecdsa.PrivateKey, msg []byte) (sig []byte, err error) {
	h := hash.Hash256(msg)
	var r, s *big.Int
	r, s, err = ecdsa.Sign(rand.Reader, priv, h)
	if err != nil {
		return
	}
	curve := priv.Curve.Params()
	for i := 0; i < (cofactor+1)*2; i++ {
		pub, err := recover(r, s, h, uint(i), true)
		if err == nil && pub.X.Cmp(priv.X) == 0 && pub.Y.Cmp(priv.Y) == 0 {
			result := make([]byte, 1, 2*curve.BitSize/8+1)
			result[0] = byte(i)
			curvelen := (curve.BitSize + 7) / 8
			bytelen := (r.BitLen() + 7) / 8
			if bytelen < curvelen {
				result = append(result,
					make([]byte, curvelen-bytelen)...)
			}
			result = append(result, r.Bytes()...)
			bytelen = (s.BitLen() + 7) / 8
			if bytelen < curvelen {
				result = append(result,
					make([]byte, curvelen-bytelen)...)
			}
			result = append(result, s.Bytes()...)
			return result, nil
		}
	}
	return nil, errors.New("no valid solution for pubkey found")
}

// Identify determines the address that signed a message
func Identify(sig []byte, msg []byte) (addr string) {
	v := uint(sig[0])
	sig = sig[1:]
	rbytes, sbytes := sig[:len(sig)/2], sig[len(sig)/2:]
	r, s := new(big.Int).SetBytes(rbytes), new(big.Int).SetBytes(sbytes)
	h := hash.Hash256(msg)
	key, err := recover(r, s, h, v, false)
	if err != nil {
		return ""
	}
	pub := compactEncodePubKey(key)
	return generateAddress(pub)
}

// Verify verifies that a message was signed by the owner of the compact encoded public key
func Verify(pub []byte, sig []byte, msg []byte) bool {
	key := compactDecodePubKey(pub)
	return VerifyWithKey(key, sig, msg)
}

// VerifyWithKey verifies that a message was signed by the owner of the public key
func VerifyWithKey(pub *ecdsa.PublicKey, sig []byte, msg []byte) bool {
	sig = sig[1:]
	rbytes, sbytes := sig[:len(sig)/2], sig[len(sig)/2:]
	h := hash.Hash256(msg)
	return ecdsa.Verify(pub, h, new(big.Int).SetBytes(rbytes), new(big.Int).SetBytes(sbytes))
}

// VerifyWithAddr verifies that a message was signed by the owner of the address
func VerifyWithAddr(addr string, sig []byte, msg []byte) bool {
	a := Identify(sig, msg)
	return a == addr
}

func compactEncodePubKey(key *ecdsa.PublicKey) (pub []byte) {
	b := key.X.Bytes()
	pub = make([]byte, 1, len(b)+1)
	pub[0] = byte(key.Y.Bit(0))
	pub = append(pub, b...)
	return
}

func compactDecodePubKey(pub []byte) *ecdsa.PublicKey {
	curve := elliptic.P256()
	sign := uint(pub[0])
	x := new(big.Int).SetBytes(pub[1:])
	y := decompressPoint(curve.Params(), x, sign)
	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
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
