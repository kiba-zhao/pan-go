package peer

import (
	"crypto/x509"

	"github.com/google/uuid"
)

type Settings interface {
	BaseId() uuid.UUID
	Cert() *x509.Certificate
	RouteDiscardThreshold() uint8
	PeerIdGeneratorDefaultDeny() bool
}

type settingsSt struct {
}

// NewSettings ...
func NewSettings() Settings {
	settings := new(settingsSt)
	return settings
}

func (s *settingsSt) BaseId() uuid.UUID {
	// TODO: Implement reading from toml file
	return uuid.New()
}
func (s *settingsSt) Cert() *x509.Certificate {
	// TODO: Implement reading from toml file
	return nil
}
func (s *settingsSt) RouteDiscardThreshold() uint8 {
	// TODO: Implement reading from toml file
	return 3
}
func (s *settingsSt) PeerIdGeneratorDefaultDeny() bool {
	// TODO: Implement reading from toml file
	return true
}
