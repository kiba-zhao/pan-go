package user

import (
	"iter"
	"pan/pkg/repository"
	"strconv"
)

const (
	UserExtraTableName = "user_extras"
)

func generateUserExtraTableName(userId uint) string {
	return UserExtraTableName + "_" + strconv.FormatUint(uint64(userId), 10)
}

type UserExtraRepository interface {
	Select(userId uint, serialNumber uint) (UserExtra, error)
	ScanWithUserID(userId uint) (iter.Seq2[UserExtra, error], error)
}

type stdUserExtraRepository struct {
	repository.RepositoryBase
}

var _ = (UserExtraRepository)((*stdUserExtraRepository)(nil))

func (repo *stdUserExtraRepository) Select(userId uint, serialNumber uint) (UserExtra, error) {
	db := repo.DB()
	if db == nil {
		return UserExtra{}, repository.ErrRepositoryDBUnavailable
	}
	userExtraTableName := generateUserExtraTableName(userId)
	if !db.Migrator().HasTable(userExtraTableName) {
		return UserExtra{}, nil
	}

	db = db.Scopes(repository.TableName(userExtraTableName))

	var userExtra UserExtra
	results := db.Where("user_id = ? and serial_number = ?", userId, serialNumber).First(&userExtra)
	return userExtra, results.Error
}

func (repo *stdUserExtraRepository) ScanWithUserID(userId uint) (iter.Seq2[UserExtra, error], error) {
	db := repo.DB()
	if db == nil {
		return nil, repository.ErrRepositoryDBUnavailable
	}
	userExtraTableName := generateUserExtraTableName(userId)
	if !db.Migrator().HasTable(userExtraTableName) {
		return nil, nil
	}

	db = db.Scopes(repository.TableName(userExtraTableName))

	rows, err := db.Where("user_id = ?", userId).Rows()
	if err != nil {
		return nil, err
	}

	return func(yield func(UserExtra, error) bool) {
		defer rows.Close()
		for rows.Next() {
			var userExtra UserExtra
			err := db.ScanRows(rows, &userExtra)
			if !yield(userExtra, err) {
				break
			}
		}
	}, nil
}
