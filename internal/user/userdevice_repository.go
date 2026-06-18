package user

import (
	"iter"
	"pan/pkg/repository"
	"strconv"
)

const (
	UserDeviceTableName = "user_devices"
)

func generateUserDeviceTableName(userId uint) string {
	return UserDeviceTableName + "_" + strconv.FormatUint(uint64(userId), 10)
}

type UserDeviceRepository interface {
	Select(userId uint, peerId string) (UserDevice, error)
	ScanWithUserID(userId uint) (iter.Seq2[UserDevice, error], error)
}

type stdUserDeviceRepository struct {
	repository.RepositoryBase
}

var _ = (UserDeviceRepository)((*stdUserDeviceRepository)(nil))

func (repo *stdUserDeviceRepository) Select(userId uint, peerId string) (UserDevice, error) {
	db := repo.DB()
	if db == nil {
		return UserDevice{}, repository.ErrRepositoryDBUnavailable
	}
	userDeviceTableName := generateUserDeviceTableName(userId)
	if !db.Migrator().HasTable(userDeviceTableName) {
		return UserDevice{}, nil
	}
	db = db.Scopes(repository.TableName(userDeviceTableName))

	var userDevice UserDevice
	results := db.Where("user_id = ? and peer_id = ?", userId, peerId).First(&userDevice)
	return userDevice, results.Error
}

func (repo *stdUserDeviceRepository) ScanWithUserID(userId uint) (iter.Seq2[UserDevice, error], error) {
	db := repo.DB()
	if db == nil {
		return nil, repository.ErrRepositoryDBUnavailable
	}
	userDeviceTableName := generateUserDeviceTableName(userId)
	if !db.Migrator().HasTable(userDeviceTableName) {
		return nil, nil
	}

	db = db.Scopes(repository.TableName(userDeviceTableName))

	rows, err := db.Where("user_id = ?", userId).Rows()
	if err != nil {
		return nil, err
	}
	return func(yield func(UserDevice, error) bool) {
		defer rows.Close()
		for rows.Next() {
			var userDevice UserDevice
			err := db.ScanRows(rows, &userDevice)
			if !yield(userDevice, err) {
				break
			}
		}
	}, nil
}
