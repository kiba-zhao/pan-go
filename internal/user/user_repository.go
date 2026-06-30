package user

import (
	"context"
	"iter"
	"pan/pkg/repository"
)

const (
	DevicesUserAssociation = "Devices"
)

type UserRepository interface {
	SelectWithGenesis(ctx context.Context, genesis string, code string) (User, error)
	Scan(ctx context.Context) (iter.Seq2[User, error], error)
}

type stdUserRepository struct {
	repository.RepositoryBase
}

var _ = (UserRepository)((*stdUserRepository)(nil))

func (repo *stdUserRepository) SelectWithGenesis(ctx context.Context, genesis string, code string) (User, error) {
	db := repo.WithContext(ctx)
	if db == nil {
		return User{}, repository.ErrRepositoryDBUnavailable
	}

	var user User
	results := db.Where("genesis_signature = ? AND code = ?", genesis, code).First(&user)
	return user, results.Error
}

func (repo *stdUserRepository) Scan(ctx context.Context) (iter.Seq2[User, error], error) {
	db := repo.WithContext(ctx)
	if db == nil {
		return nil, repository.ErrRepositoryDBUnavailable
	}

	rows, err := db.Model(&User{}).Rows()
	if err != nil {
		return nil, err
	}
	return repository.NewSeq2WithRows[User](rows, db), nil
}
