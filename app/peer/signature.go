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
