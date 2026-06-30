package user

import (
	"context"
	"errors"
	"pan/pkg/repository"
)

var ErrUserSecretRepositoryNotFound = errors.New("user.UserSecretRepository Error: Not Found")
var ErrUserSecretRepositoryConflict = errors.New("user.UserSecretRepository Error: Conflict")

type UserSecretRepository interface {
	SelectWithUserID(ctx context.Context, userId uint) (UserSecret, error)
	SaveWithUserID(ctx context.Context, secret UserSecret) (UserSecret, error)
	DeleteWithUserID(ctx context.Context, userId uint) error
}

type stdUserSecretRepository struct {
	repository.RepositoryBase
}

var _ = (UserSecretRepository)((*stdUserSecretRepository)(nil))

func (repo *stdUserSecretRepository) SelectWithUserID(ctx context.Context, userId uint) (UserSecret, error) {
	db := repo.WithContext(ctx)
	if db == nil {
		return UserSecret{}, repository.ErrRepositoryDBUnavailable
	}

	var userSecret UserSecret
	results := db.Where("user_id = ?", userId).First(&userSecret)
	return userSecret, results.Error
}

func (repo *stdUserSecretRepository) SaveWithUserID(ctx context.Context, secret UserSecret) (UserSecret, error) {
	db := repo.WithContext(ctx)
	if db == nil {
		return UserSecret{}, repository.ErrRepositoryDBUnavailable
	}

	isUpdated := secret.ID > 0
	results := db.Save(&secret)
	err := results.Error
	if err == nil && results.RowsAffected < 1 {
		if isUpdated {
			err = ErrUserSecretRepositoryNotFound
		} else {
			err = ErrUserSecretRepositoryConflict
		}
	}
	if err != nil {
		return UserSecret{}, results.Error
	}
	return secret, nil
}

func (repo *stdUserSecretRepository) DeleteWithUserID(ctx context.Context, userId uint) error {
	db := repo.WithContext(ctx)
	if db == nil {
		return repository.ErrRepositoryDBUnavailable
	}
	results := db.Where("user_id = ?", userId).Delete(&UserSecret{})
	return results.Error
}
