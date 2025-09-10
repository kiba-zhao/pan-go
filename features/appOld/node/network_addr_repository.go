package node

import (
	"iter"
	"pan/lib/feature"
	"pan/lib/repository"
)

type NetworkAddrRepository interface {
	repository.Repository

	SeqByPeerId(peerId string) (iter.Seq[NetworkAddr], error)
}

type networkAddrRepository struct {
	feature.Repository
}

func NewNetworkAddrRepository() NetworkAddrRepository {
	return &networkAddrRepository{}
}

var _ = (NetworkAddrRepository)((*networkAddrRepository)(nil))

func (repo *networkAddrRepository) SeqByPeerId(peerId string) (iter.Seq[NetworkAddr], error) {
	db := repo.DB()
	if db == nil {
		return nil, feature.ErrFeatureRepositoryDBUnavailable
	}

	rows, err := db.Model(&NetworkAddr{}).InnerJoins("AppNode", db.Where(&AppNode{PeerID: peerId})).Rows()
	if err != nil {
		return nil, err
	}
	return func(yield func(NetworkAddr) bool) {
		defer rows.Close()

		for rows.Next() {
			var model NetworkAddr
			err = db.ScanRows(rows, &model)
			if err == nil {
				if !yield(model) {
					break
				}
			}
			if err != nil {
				break
			}
		}

	}, nil
}
