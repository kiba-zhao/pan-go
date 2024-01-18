package peer

import (
	"bytes"
	"crypto/x509"
	"errors"
	"pan/core"
	"pan/memory"

	"github.com/google/uuid"
)

// compareUUIDBytes ...
func compareUUID(prev, next uuid.UUID) int {
	return bytes.Compare(prev[:], next[:])
}

type PeerIdGenerator interface {
	LocalPeerId() PeerId
	Generate(baseId []byte, node Node) (PeerId, error)
}

type PeerPassport struct {
	*memory.BucketItem[uuid.UUID]
	VerifyBaseId bool
	VerifyPeerId bool
}

type SimplePeerIdGenerator struct {
	*memory.Bucket[uuid.UUID, *PeerPassport]
	provider Provider
	peerId   PeerId
}

// NewPeerIdGenerator ...
func NewPeerIdGenerator(provider Provider) (*SimplePeerIdGenerator, error) {

	baseId := getBaseId(provider)
	cert := getCert(provider)
	peerId, err := generatePeerId(baseId, cert)
	if err != nil {
		return nil, err
	}

	bucket := memory.NewBucket[uuid.UUID, *PeerPassport](compareUUID)

	generator := new(SimplePeerIdGenerator)
	generator.Bucket = bucket
	generator.provider = provider
	generator.peerId = peerId

	return generator, nil
}

// LocalPeerId ...
func (pg *SimplePeerIdGenerator) LocalPeerId() PeerId {
	return pg.peerId
}

// Generate ...
func (pg *SimplePeerIdGenerator) Generate(baseId []byte, node Node) (peerId PeerId, err error) {
	space, err := uuid.FromBytes(baseId)
	if err != nil {
		return
	}

	cert := node.Certificate()
	id, err := generatePeerId(space, cert)
	if err != nil {
		return
	}

	pass := getPeerIdGeneratorPass(pg.provider)
	passport := pg.GetItem(id)
	if passport.VerifyPeerId == true {
		pass = !pass
	} else {
		passport = pg.GetItem(space)
		if passport.VerifyBaseId == true {
			pass = !pass
		}
	}

	if !pass {
		err = errors.New("Deny Peer Id")
		return
	}

	peerId = PeerId(id)
	return
}

// generatePeerId ...
func generatePeerId(space uuid.UUID, cert *x509.Certificate) (peerId PeerId, err error) {
	pubKey, err := core.ExtractPublicKeyFromCert(cert)
	if err != nil {
		return
	}

	peerId = uuid.NewSHA1(space, pubKey)
	return
}

// getPeerIdGeneratorPass ...
func getPeerIdGeneratorPass(provider Provider) bool {
	settings := provider.Settings()
	defaultDeny := settings.PeerIdGeneratorDefaultDeny()
	return !defaultDeny
}
