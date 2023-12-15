package core

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"
)

// ParseCertWithPem ...
func ParseCertWithPem(cert []byte) (x509Cert *x509.Certificate, err error) {
	block, _ := pem.Decode(cert)
	x509Cert, err = x509.ParseCertificate(block.Bytes)
	return
}

// ExtractPublicKeyFromCert ...
func ExtractPublicKeyFromCert(x509Cert *x509.Certificate) (pubKey []byte, err error) {
	pubKey, err = x509.MarshalPKIXPublicKey(x509Cert.PublicKey)
	return
}

// VerifyWithPublicKey ...
func VerifyWithPublicKey(key, data, sig []byte, hash crypto.Hash) (err error) {

	pubKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return
	}

	switch pubKey.(type) {
	case *rsa.PublicKey:
		err = verifyWithRSA(pubKey.(*rsa.PublicKey), hash, data, sig)
	case *ecdsa.PublicKey:
		err = verifyWithECDSA(pubKey.(*ecdsa.PublicKey), data, sig)

	default:
		err = errors.New("Unsupported Verify Algorithm")
	}

	return
}

// verifyWithRSA ...
func verifyWithRSA(publicKey *rsa.PublicKey, hash crypto.Hash, data, sig []byte) (err error) {

	hashed, err := hashWithRSA(hash, data)
	if err != nil {
		return
	}

	err = rsa.VerifyPKCS1v15(publicKey, hash, hashed, sig)
	return
}

// verifyWithECDSA ...
func verifyWithECDSA(publicKey *ecdsa.PublicKey, data, sig []byte) (err error) {

	curveName := publicKey.Params().Name
	hash, err := hashWithECDSA(curveName, data)
	if err != nil {
		return
	}

	valid := ecdsa.VerifyASN1(publicKey, hash, sig)
	if valid == false {
		err = errors.New("Verify Failed")
	}
	return
}

// SignWithPrivateKey ...
func SignWithPrivateKey(key, data []byte, hash crypto.Hash) (sig []byte, err error) {

	block, _ := pem.Decode(key)
	x509Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	switch x509Key.(type) {
	case *rsa.PrivateKey:
		sig, err = signWithRSA(x509Key.(*rsa.PrivateKey), hash, data)
	case *ecdsa.PrivateKey:
		sig, err = signWithECDSA(x509Key.(*ecdsa.PrivateKey), data)
	default:
		err = errors.New("Unsupported Sign Algorithm")
	}
	return
}

// signWithRSA ...
func signWithRSA(key *rsa.PrivateKey, hash crypto.Hash, data []byte) (sig []byte, err error) {

	hashed, err := hashWithRSA(hash, data)
	if err == nil {
		sig, err = rsa.SignPKCS1v15(nil, key, hash, hashed)
	}

	return
}

// signWithECDSA ...
func signWithECDSA(key *ecdsa.PrivateKey, data []byte) (sig []byte, err error) {
	curveName := key.Params().Name

	hash, err := hashWithECDSA(curveName, data)
	if err != nil {
		return
	}

	sig, err = ecdsa.SignASN1(rand.Reader, key, hash)
	return
}

// hashWith ...
func hashWithRSA(hash crypto.Hash, data []byte) (hashed []byte, err error) {

	hashed, err = shaWithCryptoHash(hash, data)
	return
}

// hashWithECDSA ...
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
		err = fmt.Errorf("Unsupported Curve Name: %s", curveName)
	}
	return
}

// shaWith ...
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
		err = fmt.Errorf("Unsupported Hash: %d", hash)
	}
	return
}

// SizeOfSignature ...
func SizeOfSignature(pubkey []byte, hash crypto.Hash) (size int, err error) {
	pubKey, err := x509.ParsePKIXPublicKey(pubkey)
	if err != nil {
		return
	}

	switch pubKey.(type) {
	case *rsa.PublicKey:
		size, err = sizeOfSignatureWithRSA(hash)
	case *ecdsa.PublicKey:
		size, err = sizeOfSignatureWithECDSA(pubKey.(*ecdsa.PublicKey))

	default:
		err = errors.New("Unsupported Verify Algorithm")
	}

	return

}

// sizeOfSignatureWithRSA ...
func sizeOfSignatureWithRSA(hash crypto.Hash) (size int, err error) {

	// TODO: Implement
	return
}

// sizeOfSignatureWithECDSA ...
func sizeOfSignatureWithECDSA(pubKey *ecdsa.PublicKey) (size int, err error) {
	curveName := pubKey.Params().Name

	// TODO: Implement
	switch curveName {
	case "P-224":
		size = 1
	case "P-256":
		size = 2
	case "P-384":
		size = 3
	case "P-521":
		size = 4
	default:
		err = fmt.Errorf("Unsupported Curve Name: %s", curveName)
	}
	return
}

type generateKeyAndCertConfig struct {
	algorithm     x509.PublicKeyAlgorithm
	hash          crypto.Hash
	template      *x509.Certificate
	parentCert    *x509.Certificate
	parentPrivKey any
	bits          int
}

// default ...
func defaultGenerateKeyAndCertConfig(hostname string) *generateKeyAndCertConfig {
	cfg := new(generateKeyAndCertConfig)
	cfg.algorithm = x509.ECDSA
	cfg.hash = crypto.SHA256

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	cfg.template = &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{CommonName: hostname},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	return cfg

}

type generateKeyAndCertConfigFunc func(cfg *generateKeyAndCertConfig)

// // GenerateKeyAndCertWithSettings ...
// func GenerateKeyAndCertWithSettings() generateKeyAndCertConfigFunc {
// 	return func (cfg *generateKeyAndCertConfig){}
// }

// GenerateKey ...
func GenerateKeyAndCert(cfgFns ...generateKeyAndCertConfigFunc) (key []byte, cert []byte, err error) {

	hostname, err := os.Hostname()
	if err != nil {
		return
	}
	cfg := defaultGenerateKeyAndCertConfig(hostname)
	for _, fn := range cfgFns {
		fn(cfg)
	}

	var pubKey any
	var privKey any
	var pemType string
	switch cfg.algorithm {
	case x509.RSA:
		pemType = "RSA PRIVATE KEY"
		pubKey, privKey, err = generateKeyWithRSA(cfg.bits)
	case x509.ECDSA:
		pemType = "EC PRIVATE KEY"
		pubKey, privKey, err = generateKeyWithECDSA(cfg.hash)
	default:
		err = fmt.Errorf("Unsupported Signature Algorithm: %s", cfg.algorithm.String())
	}

	if err != nil {
		return
	}

	var parentCert *x509.Certificate
	if cfg.parentCert == nil {
		parentCert = cfg.template
	} else {
		parentCert = cfg.parentCert
	}

	var parentPrivKey any
	if cfg.parentPrivKey == nil {
		parentPrivKey = privKey
	} else {
		parentPrivKey = cfg.parentPrivKey
	}

	certDer, err := x509.CreateCertificate(rand.Reader, cfg.template, parentCert, pubKey, parentPrivKey)
	if err != nil {
		return
	}

	privKeyPKCS8, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return
	}

	key = pem.EncodeToMemory(&pem.Block{Type: pemType, Bytes: privKeyPKCS8})
	cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})
	return
}

// generateKeyWithRSA ...
func generateKeyWithRSA(bits int) (pubKey any, privKey any, err error) {
	caPrivkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err == nil {
		pubKey = &caPrivkey.PublicKey
		privKey = caPrivkey
	}
	return
}

// generateKeyWithECDSA ...
func generateKeyWithECDSA(hash crypto.Hash) (pubKey any, privKey any, err error) {
	var curve elliptic.Curve

	switch hash {
	case crypto.SHA224:
		curve = elliptic.P224()
	case crypto.SHA256:
		curve = elliptic.P256()
	case crypto.SHA384:
		curve = elliptic.P384()
	case crypto.SHA512:
		curve = elliptic.P521()
	default:
		err = fmt.Errorf("Unsupported Hash: %d", hash)
	}

	if err != nil {
		return
	}

	caPrivkey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err == nil {
		pubKey = &caPrivkey.PublicKey
		privKey = caPrivkey
	}
	return
}
