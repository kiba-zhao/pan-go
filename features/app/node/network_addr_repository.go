package node

import (
	"iter"
	"pan/lib/sample"

	"gorm.io/gorm"
)

type NetworkAddrRepository interface {
	SeqByPeerId(peerId string) (iter.Seq[NetworkAddr], error)
}

type networkAddrRepository struct {
	db *gorm.DB
}

func NewNetworkAddrRepository(db *gorm.DB) NetworkAddrRepository {
	return &networkAddrRepository{db: db}
}

func (repo *networkAddrRepository) SeqByPeerId(peerId string) (iter.Seq[NetworkAddr], error) {
	db := repo.db
	if db == nil {
		return nil, sample.ErrSampleDBUnavailable
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
