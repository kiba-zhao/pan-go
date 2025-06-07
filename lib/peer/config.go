package peer

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var ErrPeerConfigSettingsUnchanged = errors.New("peer.PeerConfig Error: Settings Unchanged")

type PeerConfigListener interface {
	OnPeerConfigUpdated(settings *PeerSettings)
}

type PeerConfig interface {
	ConfigPath() string
	Load() (*PeerSettings, error)
	Save(settings *PeerSettings) error
	EnsureConfig(configPath string) error
	Subscribe(listener PeerConfigListener)
	Unsubscribe(listener PeerConfigListener)
	ConfigListeners() []PeerConfigListener
}

type stdPeerConfig struct {
	configPath   string
	locker       sync.Mutex
	configPathRW sync.RWMutex

	rw       sync.RWMutex
	settings *PeerSettings
	already  bool

	listeners   []PeerConfigListener
	listenersRW sync.RWMutex
}

var _ = (PeerConfig)((*stdPeerConfig)(nil))

func (cfg *stdPeerConfig) ConfigPath() string {
	cfg.configPathRW.RLock()
	defer cfg.configPathRW.RUnlock()
	return cfg.configPath
}

func (cfg *stdPeerConfig) Load() (*PeerSettings, error) {
	cfg.locker.Lock()
	defer cfg.locker.Unlock()

	cfgPath := cfg.ConfigPath()

	privKeyPath := generatePrivateKeyPath(cfgPath)
	privKeyPemBytes, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}

	certPath := generateCertificatePath(cfgPath)
	certPemBytes, err := os.ReadFile(certPath)
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

	settings := &PeerSettings{}
	settings.certificate = certificate
	settings.privKey = certificate.PrivateKey
	settings.pubKey = x509Cert.PublicKey
	settings.peerId = PeerID(pubKeyBytes)

	err = VerifyPeerSettings(settings)
	if err == nil {
		onPeerSettingsChanged(cfg, settings)
	}
	return settings, err
}

func (cfg *stdPeerConfig) Save(settings *PeerSettings) error {
	cfg.locker.Lock()
	defer cfg.locker.Unlock()

	err := VerifyPeerSettings(settings)
	if err != nil {
		return err
	}

	privKeyPemBytes, err := EncodePrivateKeyToPemBytes(settings.PrivateKey())
	if err != nil {
		return err
	}
	certPemBytes := EncodeCertificateToPemBytes(settings.Certificate().Certificate...)

	err = savePrivKeyAndCert(cfg.ConfigPath(), privKeyPemBytes, certPemBytes)
	if err != nil {
		return err
	}

	onPeerSettingsChanged(cfg, settings)
	return err
}

func (cfg *stdPeerConfig) EnsureConfig(configPath string) error {
	cfg.configPathRW.Lock()
	cfg.configPath = configPath
	cfg.configPathRW.Unlock()

	_, err := cfg.Load()
	if err == nil {
		return nil
	}

	settings, err := GeneratePeerSettings()
	if err != nil {
		return err
	}

	return cfg.Save(settings)
}

func (cfg *stdPeerConfig) ConfigListeners() []PeerConfigListener {
	cfg.listenersRW.RLock()
	defer cfg.listenersRW.RUnlock()

	return cfg.listeners
}

func (cfg *stdPeerConfig) Subscribe(listener PeerConfigListener) {
	cfg.listenersRW.Lock()
	defer cfg.listenersRW.Unlock()

	settings, already := cfg.PeerSettingsAndAlready()
	if already {
		listener.OnPeerConfigUpdated(settings)
	}
	cfg.listeners = append(cfg.listeners, listener)
}

func (cfg *stdPeerConfig) Unsubscribe(listener PeerConfigListener) {
	cfg.listenersRW.Lock()
	defer cfg.listenersRW.Unlock()
	for i, l := range cfg.listeners {
		if l == listener {
			cfg.listeners = append(cfg.listeners[:i], cfg.listeners[i+1:]...)
			break
		}
	}
}

func (cfg *stdPeerConfig) PeerSettingsAndAlready() (*PeerSettings, bool) {
	cfg.rw.RLock()
	defer cfg.rw.RUnlock()
	return cfg.settings, cfg.already
}

func onPeerSettingsChanged(cfg *stdPeerConfig, settings *PeerSettings) {
	cfg.rw.Lock()
	cfg.already = true
	cfg.settings = settings
	cfg.rw.Unlock()

	listeners := cfg.ConfigListeners()
	if len(listeners) > 0 {
		for _, l := range listeners {
			l.OnPeerConfigUpdated(settings)
		}
	}
}

func savePrivKeyAndCert(configPath string, privKeyPemBytes, certPemBytes []byte) error {

	err := os.MkdirAll(configPath, 0750)
	if err != nil {
		return err
	}

	privKeyPath := generatePrivateKeyPath(configPath)
	privKeyFile, err := os.Create(privKeyPath)
	if err != nil {
		return err
	}
	defer privKeyFile.Close()

	_, err = privKeyFile.Write(privKeyPemBytes)
	if err != nil {
		return err
	}

	certPath := generateCertificatePath(configPath)
	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()

	_, err = certFile.Write(certPemBytes)
	return err
}

func generatePrivateKeyPath(configPath string) string {
	return filepath.Join(configPath, "key.pem")
}

func generateCertificatePath(configPath string) string {
	return filepath.Join(configPath, "cert.pem")
}
