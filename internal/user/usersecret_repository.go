package user

import (
	"errors"
	"pan/internal/repository"
)

var ErrUserSecretRepositoryNotFound = errors.New("user.UserSecretRepository Error: Not Found")
var ErrUserSecretRepositoryConflict = errors.New("user.UserSecretRepository Error: Conflict")

type UserSecretRepository interface {
	SelectWithUserID(userId uint) (UserSecret, error)
	SaveWithUserID(secret UserSecret) (UserSecret, error)
	DeleteWithUserID(userId uint) error
}

type stdUserSecretRepository struct {
	repository.RepositoryBase
}

var _ = (UserSecretRepository)((*stdUserSecretRepository)(nil))

func (repo *stdUserSecretRepository) SelectWithUserID(userId uint) (UserSecret, error) {
	db := repo.DB()
	if db == nil {
		return UserSecret{}, repository.ErrRepositoryDBUnavailable
	}

	var userSecret UserSecret
	results := db.Where("user_id = ?", userId).First(&userSecret)
	return userSecret, results.Error
}

func (repo *stdUserSecretRepository) SaveWithUserID(secret UserSecret) (UserSecret, error) {
	db := repo.DB()
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

func (repo *stdUserSecretRepository) DeleteWithUserID(userId uint) error {
	db := repo.DB()
	if db == nil {
		return repository.ErrRepositoryDBUnavailable
	}
	results := db.Where("user_id = ?", userId).Delete(&UserSecret{})
	return results.Error
}
