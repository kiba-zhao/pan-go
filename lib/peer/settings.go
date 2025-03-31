// Define Peer Settings
package peer

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/fs"
	"math/big"
	"os"
	"pan/lib/config"
	"pan/lib/injection"
	"pan/lib/runtime"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

var ErrPeerSettingsUnavailable = errors.New("peer.PeerSettings Error: Unavailable")

type PeerSettings interface {

	// id of peer
	PeerID() PeerID
	// public key of peer
	PubKey() any
	// private key of peer
	PrivKey() crypto.PrivateKey
	// certificate of peer
	Certificate() tls.Certificate
	// if peer settings is available
	Available() bool
}

// PeerSettingsListener is a listener for peer settings updates
type PeerSettingsListener interface {
	// OnPeerSettingsUpdated is called when the peer settings are updated.
	//
	// The function is called with the new peer settings.
	OnPeerSettingsUpdated(PeerSettings)
}

type peerSettings struct {
	Config     config.AppConfig
	registry   runtime.Registry
	registryRW sync.RWMutex
	locker     sync.RWMutex
	peerId     PeerID
	pubKey     any
	privKey    crypto.PrivateKey
	cert       tls.Certificate
	hashCode   []byte
}

// Init initializes the peer settings with the given registry.
//
// It sets the registry and does not return an error.
func (ns *peerSettings) Init(registry runtime.Registry) error {

	ns.registryRW.Lock()
	ns.registry = registry
	ns.registryRW.Unlock()

	return ns.Generate()
}

// Components returns a slice of injection.Component representing the components
// provided by the peer settings. Currently, the only component provided is the
// PeerSettings itself, which is scoped internally.
func (ns *peerSettings) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(ns, injection.ComponentNoneScope),
	}
}

// Generate initializes and updates the peer settings by parsing fields from the
// configuration path. If parsing fails due to a path error, it generates the fields.
// It notifies all registered PeerSettingsListeners about the updates. The function
// returns an error if the registry is unavailable or if any error occurs during the
// parsing or generation of fields.

func (ns *peerSettings) Generate() error {
	ns.registryRW.RLock()
	registry := ns.registry
	ns.registryRW.RUnlock()
	if registry == nil {
		return ErrPeerSettingsUnavailable
	}

	ns.locker.Lock()
	configPath := filepath.Dir(ns.Config.ConfigFilePath())
	err := ns.ParseFields(configPath)
	if _, ok := err.(*fs.PathError); ok {
		err = ns.GenerateFields(configPath)
	}
	ns.locker.Unlock()

	if err == nil {
		listeners := runtime.ModulesForType[PeerSettingsListener](registry)
		for _, listener := range listeners {
			listener.OnPeerSettingsUpdated(ns)
		}
	}

	return err
}

// GenerateFields generates a pair of EC private key and certificate and writes them to the paths
// derived from the configPath. It then updates the peer settings with the generated key pair and
// notifies all registered PeerSettingsListeners about the updates. The function returns an error if
// any error occurs during the generation or writing of the key pair.
func (ns *peerSettings) GenerateFields(configPath string) error {

	caPrivkey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&caPrivkey.PublicKey)
	if err != nil {
		return err
	}

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
		return err
	}

	privKeyPKCS8, err := x509.MarshalPKCS8PrivateKey(caPrivkey)
	if err != nil {
		return err
	}

	keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privKeyPKCS8})
	certPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})

	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return err
	}

	privKeyPath, certificatePath := generatePrivKeyPathAndCertificatePath(configPath)
	err = os.MkdirAll(filepath.Dir(privKeyPath), 0750)
	if err != nil {
		return err
	}
	keyFile, err := os.Create(privKeyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	_, err = keyFile.Write(keyPEMBlock)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(certificatePath), 0750)
	if err != nil {
		return err
	}
	certFile, err := os.Create(certificatePath)
	if err != nil {
		return err
	}
	defer certFile.Close()
	_, err = certFile.Write(certPEMBlock)
	if err != nil {
		return err
	}

	hash := sha512.New()
	hash.Write(certPEMBlock)
	hash.Write(keyPEMBlock)
	hashCode := hash.Sum(nil)

	ns.peerId = pubKeyBytes
	ns.pubKey = &caPrivkey.PublicKey
	ns.privKey = caPrivkey
	ns.cert = cert
	ns.hashCode = hashCode
	ns.locker.Unlock()

	return err
}

// ParseFields reads the private key and certificate from the given configPath and
// loads the peer settings from them. If the peer settings are already loaded and
// the hash of the private key and certificate matches the one stored in the peer
// settings, the function returns nil without doing anything. Otherwise, it will
// load the peer settings from the given files and store the hash of the private key
// and certificate in the peer settings. If any error occurs, it is returned.
func (ns *peerSettings) ParseFields(configPath string) error {

	privKeyPath, certificatePath := generatePrivKeyPathAndCertificatePath(configPath)

	var err error
	keyPEMBlock, err := os.ReadFile(privKeyPath)
	var certPEMBlock []byte
	var cert tls.Certificate
	var privKey crypto.PrivateKey
	var hashCode []byte
	var x509Cert *x509.Certificate
	if err == nil {
		certPEMBlock, err = os.ReadFile(certificatePath)
	}

	if err == nil {
		hash := sha512.New()
		hash.Write(certPEMBlock)
		hash.Write(keyPEMBlock)
		hashCode = hash.Sum(nil)
		if ns.hashCode != nil && bytes.Equal(ns.hashCode, hashCode) {
			return nil
		}

		cert, err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	}

	if err == nil {
		block, _ := pem.Decode(keyPEMBlock)
		privKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	if err == nil {
		x509Cert, err = x509.ParseCertificate(cert.Certificate[0])
	}

	if err != nil {
		return err
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(x509Cert.PublicKey)
	if err == nil {
		ns.peerId = pubKeyBytes
		ns.pubKey = x509Cert.PublicKey
		ns.privKey = privKey
		ns.cert = cert
		ns.hashCode = hashCode
	}

	return err
}

// PeerID returns the peer ID of the peer settings, which is the marshaled bytes of
// the public key of the peer's certificate.
func (ns *peerSettings) PeerID() PeerID {
	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.peerId
}

// PubKey returns the public key of the peer's certificate, which is used to
// identify the peer.
func (ns *peerSettings) PubKey() any {
	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.pubKey
}

// PrivKey returns the private key of the peer's certificate, which is used to
// decrypt encrypted messages and sign messages.
func (ns *peerSettings) PrivKey() crypto.PrivateKey {
	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.privKey
}

// Certificate returns the TLS certificate of the peer's settings,
// which is used for establishing secure connections and authenticating
// the peer within the p2p network.

func (ns *peerSettings) Certificate() tls.Certificate {

	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.cert
}

// Available returns true if the peer settings have been loaded and the hash of the
// private key and certificate has been stored, and false otherwise.
func (ns *peerSettings) Available() bool {
	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.hashCode != nil
}

// EngineTypes returns a slice of reflect.Type representing the various engine types
// associated with the peer settings module. These types include:
//
//   - PeerSettingsListener: Represents a listener for peer settings updates.
func (ns *peerSettings) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[PeerSettingsListener](),
	}
}

// generatePrivKeyPathAndCertificatePath generates the paths for the private key and
// certificate of the peer settings based on the given root path. The private key is
// stored as "key.pem" and the certificate is stored as "cert.pem" in the root path.
func generatePrivKeyPathAndCertificatePath(rootPath string) (string, string) {

	return filepath.Join(rootPath, "key.pem"), filepath.Join(rootPath, "cert.pem")
}
