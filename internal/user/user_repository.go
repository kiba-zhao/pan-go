package user

import (
	"pan/pkg/repository"
)

const (
	DevicesUserAssociation = "Devices"
)

type UserRepository interface {
	SelectWithGenesis(genesis string, code string) (User, error)
}

type stdUserRepository struct {
	repository.RepositoryBase
}

var _ = (UserRepository)((*stdUserRepository)(nil))

func (repo *stdUserRepository) SelectWithGenesis(genesis string, code string) (User, error) {
	db := repo.DB()
	if db == nil {
		return User{}, repository.ErrRepositoryDBUnavailable
	}

	var user User
	results := db.Where("genesis_signature = ? AND code = ?", genesis, code).First(&user)
	return user, results.Error
}
