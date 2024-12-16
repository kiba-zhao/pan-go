package appbroadcast

import (
	"errors"
	"pan/app/sample"

	"gorm.io/gorm"
)

var ErrAppBroadcastInfoNotFound = errors.New("appbroadcast.AppBroadcastInfoRepository Error: Not Found")

type AppBroadcastInfoRepository interface {
	SelectOrCreate(baseInfo AppBroadcastInfo) (AppBroadcastInfo, bool, error)
	UpdateHeightest(info AppBroadcastInfo) (AppBroadcastInfo, error)
	DeleteByPeerID(peerId string) error
}

type appBroadcastInfoRepository struct {
	db *gorm.DB
}

func NewAppBroadcastInfoRepository(db *gorm.DB) AppBroadcastInfoRepository {
	return &appBroadcastInfoRepository{db: db}
}

func (repo *appBroadcastInfoRepository) SelectOrCreate(baseInfo AppBroadcastInfo) (AppBroadcastInfo, bool, error) {

	db := repo.db
	if db == nil {
		return AppBroadcastInfo{}, false, sample.ErrSampleDBUnavailable
	}
	var info AppBroadcastInfo
	results := db.Where(AppBroadcastInfo{PeerID: baseInfo.PeerID}).Attrs(baseInfo).FirstOrCreate(&info)
	if results.Error == gorm.ErrRecordNotFound {
		return info, results.RowsAffected == 0, ErrAppBroadcastInfoNotFound
	}
	return info, results.RowsAffected == 0, results.Error
}

func (repo *appBroadcastInfoRepository) UpdateHeightest(info AppBroadcastInfo) (AppBroadcastInfo, error) {
	db := repo.db
	if db == nil {
		return AppBroadcastInfo{}, sample.ErrSampleDBUnavailable
	}

	results := db.Model(&info).Where("hightest < ?", info.Hightest).Limit(1)
	if results.Error == nil && results.RowsAffected != 1 {
		return info, ErrAppBroadcastInfoNotFound
	}
	return info, results.Error
}

func (repo *appBroadcastInfoRepository) DeleteByPeerID(peerId string) error {
	db := repo.db
	if db == nil {
		return sample.ErrSampleDBUnavailable
	}
	results := db.Delete(&AppBroadcastInfo{PeerID: peerId}).Limit(1)
	if results.Error == nil && results.RowsAffected != 1 {
		return ErrAppBroadcastInfoNotFound
	}
	return results.Error
}
