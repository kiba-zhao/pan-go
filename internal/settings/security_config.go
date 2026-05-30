package settings

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
	"os"
	"pan/internal/config"
	"pan/internal/net"
	"path/filepath"
	"time"
)

var ErrSettingsInvalidSecurityConfig = errors.New("settings.SecurityConfig Error: Invalid SecurityConfig")

var (
	PrivateKeyFileName  = "key.pem"
	CertificateFileName = "cert.pem"
)

type SecurityConfig interface {
	Certificate() tls.Certificate
	PublicKey() any
	PrivateKey() crypto.PrivateKey
	PeerID() net.PeerID
}

type SecurityConfigurerListener = config.ConfigurerListener[SecurityConfig]
type SecurityConfigurer = config.Configurer[SecurityConfig]

func newSecurityConfig(homePath string) (SecurityConfig, error) {
	security, err := loadSecurityConfig(homePath)
	if err == nil {
		return security, err
	}

	security, err = generateSecurityConfig()
	if err != nil {
		return nil, err
	}

	err = saveSecurityConfig(homePath, security)
	return security, err
}

type stdSecurityConfig struct {
	certificate tls.Certificate
	pubKey      any
	privKey     crypto.PrivateKey
	peerId      net.PeerID
}

var _ = ((SecurityConfig)((*stdSecurityConfig)(nil)))

func (s *stdSecurityConfig) Certificate() tls.Certificate {
	return s.certificate
}

func (s *stdSecurityConfig) PublicKey() any {
	return s.pubKey
}

func (s *stdSecurityConfig) PeerID() net.PeerID {
	return s.peerId
}

func (s *stdSecurityConfig) PrivateKey() crypto.PrivateKey {
	return s.privKey
}

func VerifySecurityConfig(security SecurityConfig) error {
	certificate := security.Certificate()

	// Verify the public key
	x509Cert, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		return err
	}
	certPublicKeyBytes, err := x509.MarshalPKIXPublicKey(x509Cert.PublicKey)
	if err != nil {
		return err
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(security.PublicKey())
	if err != nil {
		return err
	}
	if !bytes.Equal(certPublicKeyBytes, publicKeyBytes) || !bytes.Equal(certPublicKeyBytes, security.PeerID()) {
		return ErrSettingsInvalidSecurityConfig
	}
	//

	// Verify the private key
	certPrivateKeyBytes, err := x509.MarshalPKCS8PrivateKey(certificate.PrivateKey)
	if err != nil {
		return err
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(security.PrivateKey())
	if err != nil {
		return err
	}
	if !bytes.Equal(certPrivateKeyBytes, privateKeyBytes) {
		return ErrSettingsInvalidSecurityConfig
	}
	//

	return net.VerifyPairKey(security.PrivateKey(), security.PublicKey())
}

func loadSecurityConfig(homePath string) (SecurityConfig, error) {

	privKeyPemBytes, err := os.ReadFile(filepath.Join(homePath, PrivateKeyFileName))
	if err != nil {
		return nil, err
	}

	certPemBytes, err := os.ReadFile(filepath.Join(homePath, CertificateFileName))
	if err != nil {
		return nil, err
	}

	certificate, err := tls.X509KeyPair(certPemBytes, privKeyPemBytes)
	if err != nil {
		return nil, err
	}

	x509Cert, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(x509Cert.PublicKey)
	if err != nil {
		return nil, err
	}

	security := &stdSecurityConfig{}
	security.certificate = certificate
	security.privKey = certificate.PrivateKey
	security.pubKey = pubKeyBytes
	security.peerId = net.PeerID(pubKeyBytes)

	err = VerifySecurityConfig(security)
	if err != nil {
		return nil, err
	}

	return security, nil
}

func saveSecurityConfig(homePath string, security SecurityConfig) error {
	err := VerifySecurityConfig(security)
	if err != nil {
		return err
	}

	privKeyPemBytes, err := net.EncodePrivateKeyToPemBytes(security.PrivateKey())
	if err != nil {
		return err
	}
	certPemBytes := net.EncodeCertificateToPemBytes(security.Certificate().Certificate...)

	err = os.MkdirAll(homePath, 0750)
	if err != nil {
		return err
	}

	privKeyFile, err := os.Create(filepath.Join(homePath, PrivateKeyFileName))
	if err != nil {
		return err
	}
	defer privKeyFile.Close()

	_, err = privKeyFile.Write(privKeyPemBytes)
	if err != nil {
		return err
	}

	certFile, err := os.Create(filepath.Join(homePath, CertificateFileName))
	if err != nil {
		return err
	}
	defer certFile.Close()

	_, err = certFile.Write(certPemBytes)
	return err
}

func generateSecurityConfig() (SecurityConfig, error) {
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

	security := &stdSecurityConfig{}
	security.peerId = net.PeerID(pubKeyBytes)
	security.pubKey = &caPrivkey.PublicKey
	security.privKey = caPrivkey
	security.certificate = cert
	return security, nil
}
