package user

import (
	"context"
	"iter"
	"pan/pkg/repository"
	"strconv"
)

const (
	UserConsensusTableName = "user_consensuses"
)

func generateUserConsensusTableName(userId uint) string {
	return UserConsensusTableName + "_" + strconv.FormatUint(uint64(userId), 10)
}

type UserConsensusRepository interface {
	Select(ctx context.Context, userId uint, height uint64) (UserConsensus, error)
	SelectLatestWithUserID(ctx context.Context, userId uint) (UserConsensus, error)
	ScanWithUserIDLessThanHeight(ctx context.Context, userId uint, height uint64) (iter.Seq2[UserConsensus, error], error)
}

type stdUserConsensusRepository struct {
	repository.RepositoryBase
}

var _ = (UserConsensusRepository)((*stdUserConsensusRepository)(nil))

func (repo *stdUserConsensusRepository) Select(ctx context.Context, userId uint, height uint64) (UserConsensus, error) {
	db := repo.WithContext(ctx)
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

func (repo *stdUserConsensusRepository) SelectLatestWithUserID(ctx context.Context, userId uint) (UserConsensus, error) {
	db := repo.WithContext(ctx)
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

func (repo *stdUserConsensusRepository) ScanWithUserIDLessThanHeight(ctx context.Context, userId uint, height uint64) (iter.Seq2[UserConsensus, error], error) {
	db := repo.WithContext(ctx)
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
	return repository.NewSeq2WithRows[UserConsensus](rows, db), nil
}
