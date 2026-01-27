// Define Peer Settings
package peer

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

var ErrPeerSettingsUnavailable = errors.New("peer.PeerSettings Error: Unavailable")
var ErrPeerSettingsInvalid = errors.New("peer.PeerSettings Error: Invalid")

type PeerSettings struct {
	certificate tls.Certificate
	pubKey      any
	peerId      PeerID
	privKey     crypto.PrivateKey
}

func (ps *PeerSettings) Certificate() tls.Certificate {
	return ps.certificate
}

func (ps *PeerSettings) PublicKey() any {
	return ps.pubKey
}

func (ps *PeerSettings) PeerID() PeerID {
	return ps.peerId
}

func (ps *PeerSettings) PrivateKey() crypto.PrivateKey {
	return ps.privKey
}

func VerifyPeerSettings(settings *PeerSettings) error {
	certificate := settings.Certificate()

	// Verify the public key
	x509Cert, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		return err
	}
	certPublicKeyBytes, err := x509.MarshalPKIXPublicKey(x509Cert.PublicKey)
	if err != nil {
		return err
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(settings.PublicKey())
	if err != nil {
		return err
	}
	if !bytes.Equal(certPublicKeyBytes, publicKeyBytes) || !bytes.Equal(certPublicKeyBytes, settings.PeerID()) {
		return ErrPeerSettingsInvalid
	}
	//

	// Verify the private key
	certPrivateKeyBytes, err := x509.MarshalPKCS8PrivateKey(certificate.PrivateKey)
	if err != nil {
		return err
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(settings.PrivateKey())
	if err != nil {
		return err
	}
	if !bytes.Equal(certPrivateKeyBytes, privateKeyBytes) {
		return ErrPeerSettingsInvalid
	}
	//

	return VerifyPairKey(settings.PrivateKey(), settings.PublicKey())
}

func GeneratePeerSettings() (*PeerSettings, error) {
	// generate private key
	caPrivkey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}
	//

	// generate public key
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&caPrivkey.PublicKey)
	if err != nil {
		return nil, err
	}
	//

	// generate certificate
	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值
	template := &x509.Certificate{
		SerialNumber:          serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, // 典型用法是指定叶子证书中的公钥的使用目的。它包括一系列的OID，每一个都指定一种用途。例如{id pkix 31}表示用于服务器端的TLS/SSL连接；{id pkix 34}表示密钥可以用于保护电子邮件。
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,                      // 指定了这份证书包含的公钥可以执行的密码操作，例如只能用于签名，但不能用来加密
		IsCA:                  true,                                                                       // 指示证书是不是ca证书
		BasicConstraintsValid: true,                                                                       // 指示证书是不是ca证书
	}

	certDer, err := x509.CreateCertificate(rand.Reader, template, template, &caPrivkey.PublicKey, caPrivkey)
	if err != nil {
		return nil, err
	}

	privKeyPKCS8, err := x509.MarshalPKCS8PrivateKey(caPrivkey)
	if err != nil {
		return nil, err
	}

	keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privKeyPKCS8})
	certPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})

	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}
	//

	settings := &PeerSettings{}
	settings.peerId = PeerID(pubKeyBytes)
	settings.pubKey = &caPrivkey.PublicKey
	settings.privKey = caPrivkey
	settings.certificate = cert
	return settings, nil
}
