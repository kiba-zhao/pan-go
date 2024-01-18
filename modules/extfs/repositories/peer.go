package repositories

import (
	"pan/modules/extfs/models"

	"gorm.io/gorm"
)

type PeerRepository interface {
	FindOne(id string) (models.Peer, error)
	Save(peer models.Peer) error
}

type peerRepositoryImpl struct {
	DB *gorm.DB
}

// NewPeerRepository ...
func NewPeerRepository(db *gorm.DB) PeerRepository {
	repo := new(peerRepositoryImpl)
	repo.DB = db
	return repo
}

// FineOne ...
func (repo *peerRepositoryImpl) FindOne(id string) (models.Peer, error) {

	var peer models.Peer
	peer.ID = id
	results := repo.DB.Take(&peer)
	if results.Error != nil {
		return peer, results.Error
	}

	return peer, nil
}

// Save ...
func (repo *peerRepositoryImpl) Save(peer models.Peer) error {

	results := repo.DB.Save(&peer)
	return results.Error
}
