// Define peer signature
package peer

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"errors"
)

var ErrPeerSignatureInvalidPrivKey = errors.New("peer.signature Error: Invalid Private Key")
var ErrPeerSignatureInvalidPublicKey = errors.New("peer.signature Error: Invalid Public Key")
var ErrPeerSignatureUnsupportedHash = errors.New("peer.signature Error: Unsupported Hash")
var ErrPeerSignatureUnsupportedCurve = errors.New("peer.signature Error: Unsupported Curve")
var ErrPeerSignatureVerifyFailed = errors.New("peer.signature Error: Verify Failed")

// Sign generates a signature for the given data using the given private key.
//
// The key must be an *ecdsa.PrivateKey, and the curve name must be one of:
//   - P-224
//   - P-256
//   - P-384
//   - P-521
//
// The signature algorithm used is ECDSA with the given curve.
//
// The error returned is ErrPeerSignatureInvalidPrivKey if the key is not an *ecdsa.PrivateKey,
// and ErrPeerSignatureUnsupportedCurve if the curve is not supported.
func Sign(data []byte, key crypto.PrivateKey) ([]byte, error) {
	privKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrPeerSignatureInvalidPrivKey
	}

	curveName := privKey.Params().Name
	hash, err := hashWithECDSA(curveName, data)
	if err != nil {
		return nil, err
	}

	return ecdsa.SignASN1(rand.Reader, privKey, hash)
}

// Verify checks the validity of a given signature for the provided data using
// the specified public key.
//
// The `key` must be in the PKIX, ASN.1 DER format and represent an ECDSA public
// key. If the public key cannot be parsed or is not an ECDSA key, the function
// returns ErrPeerSignatureInvalidPublicKey.
//
// The function calculates a hash of the `data` based on the curve used by the
// public key and verifies the signature using ECDSA. If the signature does not
// match, it returns ErrPeerSignatureVerifyFailed.
//
// Returns an error if the public key is invalid, the curve is unsupported, or
// the signature verification fails.

func Verify(data, sig, key []byte) error {
	x509Key, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return err
	}
	pubKey, ok := x509Key.(*ecdsa.PublicKey)
	if !ok {
		return ErrPeerSignatureInvalidPublicKey
	}

	curveName := pubKey.Params().Name
	hash, err := hashWithECDSA(curveName, data)
	if err != nil {
		return err
	}

	valid := ecdsa.VerifyASN1(pubKey, hash, sig)
	if !valid {
		return ErrPeerSignatureVerifyFailed
	}

	return nil
}

// hashWithECDSA generates a hash for the given data using the specified curve.
//
// The `curveName` parameter must be one of the following curve names:
//   - P-224
//   - P-256
//   - P-384
//   - P-521
//
// The function returns the hash as a byte slice, and an error if the curve is
// unsupported.
func hashWithECDSA(curveName string, data []byte) (hash []byte, err error) {
	switch curveName {
	case "P-224":
		hash, err = shaWithCryptoHash(crypto.SHA224, data)
	case "P-256":
		hash, err = shaWithCryptoHash(crypto.SHA256, data)
	case "P-384":
		hash, err = shaWithCryptoHash(crypto.SHA384, data)
	case "P-521":
		hash, err = shaWithCryptoHash(crypto.SHA512, data)
	default:
		err = ErrPeerSignatureUnsupportedCurve
	}
	return
}

// shaWithCryptoHash hashes the input data using the specified hash algorithm.
//
// The `hash` parameter is a crypto.Hash value that specifies the hashing algorithm,
// and it must be one of the following:
//   - crypto.SHA1
//   - crypto.SHA224
//   - crypto.SHA256
//   - crypto.SHA384
//   - crypto.SHA512
//
// The function returns the hashed data as a byte slice. If the specified hash
// algorithm is unsupported, it returns an error indicating the issue.

func shaWithCryptoHash(hash crypto.Hash, data []byte) (hashed []byte, err error) {
	switch hash {
	case crypto.SHA1:
		hash20 := sha1.Sum(data)
		hashed = hash20[:]
	case crypto.SHA224:
		hash32 := sha256.Sum224(data)
		hashed = hash32[:]
	case crypto.SHA256:
		hash32 := sha256.Sum256(data)
		hashed = hash32[:]
	case crypto.SHA384:
		hash48 := sha512.Sum384(data)
		hashed = hash48[:]
	case crypto.SHA512:
		hash48 := sha512.Sum512(data)
		hashed = hash48[:]
	default:
		err = ErrPeerSignatureUnsupportedHash
	}
	return
}
