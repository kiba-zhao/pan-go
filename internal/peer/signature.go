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
	"encoding/pem"
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

func Verify(data, sig []byte, key any) error {
	pubKey, ok := key.(*ecdsa.PublicKey)
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

func VerifyWithPublicKeyBytes(data, sig, key []byte) error {
	x509Key, err := x509.ParsePKIXPublicKey(key)
	if err == nil {
		err = Verify(data, sig, x509Key)
	}
	return err
}

func VerifyPairKey(privKey crypto.PrivateKey, pubKey any) error {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return err
	}

	sig, err := Sign(randomBytes, privKey)
	if err != nil {
		return err
	}

	return Verify(randomBytes, sig, pubKey)

}

func VerifyPairKeyBytes(privKey, pubKey []byte) error {
	x509PrivKey, err := x509.ParseECPrivateKey(privKey)
	if err == nil {
		err = VerifyPairKey(x509PrivKey, pubKey)
	}
	return err
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

func EncodePrivateKeyToPemBytes(key crypto.PrivateKey) ([]byte, error) {
	privKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrPeerSignatureInvalidPrivKey
	}

	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{}
	block.Type = "EC PRIVATE KEY"
	block.Bytes = privKeyBytes

	return pem.EncodeToMemory(block), nil
}

func EncodeCertificateToPemBytes(certificates ...[]byte) []byte {
	certificateBytes := make([]byte, 0)
	for _, certificate := range certificates {
		certPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate})
		certificateBytes = append(certificateBytes, certPEMBlock...)
	}
	return certificateBytes
}
