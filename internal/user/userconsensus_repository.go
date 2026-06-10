package user

import (
	"iter"
	"pan/internal/repository"
	"strconv"
)

const (
	UserConsensusTableName = "user_consensuses"
)

func generateUserConsensusTableName(userId uint) string {
	return UserConsensusTableName + "_" + strconv.FormatUint(uint64(userId), 10)
}

type UserConsensusRepository interface {
	Select(userId uint, height uint64) (UserConsensus, error)
	SelectLatestWithUserID(userId uint) (UserConsensus, error)
	ScanWithUserIDLessThanHeight(userId uint, height uint64) (iter.Seq2[UserConsensus, error], error)
}

type stdUserConsensusRepository struct {
	repository.RepositoryBase
}

var _ = (UserConsensusRepository)((*stdUserConsensusRepository)(nil))

func (repo *stdUserConsensusRepository) Select(userId uint, height uint64) (UserConsensus, error) {
	db := repo.DB()
	if db == nil {
		return UserConsensus{}, repository.ErrRepositoryDBUnavailable
	}
	userConsensusTableName := generateUserConsensusTableName(userId)
	if !db.Migrator().HasTable(userConsensusTableName) {
		return UserConsensus{}, nil
	}

	db = db.Scopes(repository.TableName(userConsensusTableName))

	var userConsensus UserConsensus
	results := db.Where("user_id = ? and height = ?", userId, height).First(&userConsensus)
	return userConsensus, results.Error
}

func (repo *stdUserConsensusRepository) SelectLatestWithUserID(userId uint) (UserConsensus, error) {
	db := repo.DB()
	if db == nil {
		return UserConsensus{}, repository.ErrRepositoryDBUnavailable
	}
	userConsensusTableName := generateUserConsensusTableName(userId)
	if !db.Migrator().HasTable(userConsensusTableName) {
		return UserConsensus{}, nil
	}

	db = db.Scopes(repository.TableName(userConsensusTableName))

	var userConsensus UserConsensus
	results := db.Where("user_id = ?", userId).Last(&userConsensus)
	return userConsensus, results.Error
}

func (repo *stdUserConsensusRepository) ScanWithUserIDLessThanHeight(userId uint, height uint64) (iter.Seq2[UserConsensus, error], error) {
	db := repo.DB()
	if db == nil {
		return nil, repository.ErrRepositoryDBUnavailable
	}
	userConsensusTableName := generateUserConsensusTableName(userId)
	if !db.Migrator().HasTable(userConsensusTableName) {
		return nil, nil
	}
	db = db.Scopes(repository.TableName(userConsensusTableName))

	rows, err := db.Where("user_id = ? And height <= ?", userId, height).Order("id asc").Rows()
	if err != nil {
		return nil, err
	}
	return func(yield func(UserConsensus, error) bool) {
		defer rows.Close()
		for rows.Next() {
			var userConsensus UserConsensus
			err := db.ScanRows(rows, &userConsensus)
			if !yield(userConsensus, err) {
				break
			}
		}
	}, nil
}
