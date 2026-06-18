package user

import (
	"errors"
	"pan/pkg/repository"
)

var ErrUserDataRepositoryInvalidData = errors.New("user.UserDataRepository Error: Invalid Data")

type UserDataRepository interface {
	Save(user User, secret UserSecret, userConsensuses []UserConsensus, userDevices []UserDevice, userExtras []UserExtra) error
	Clean(user User) error
}

type stdUserDataRepository struct {
	*repository.RepositoryBase
}

var _ = (UserDataRepository)((*stdUserDataRepository)(nil))

func (repo *stdUserDataRepository) Save(user User, secret UserSecret, userConsensuses []UserConsensus, userDevices []UserDevice, userExtras []UserExtra) error {
	db := repo.DB()
	if db == nil {
		return repository.ErrRepositoryDBUnavailable
	}

	if len(userDevices) < 1 || len(userConsensuses) < 1 {
		return ErrUserDataRepositoryInvalidData
	}

	tx := db.Begin()

	isNew := false
	results := tx.Model(&user).Where("code = ? and genesis_signature = ? and height < ?", user.Code, user.GenesisSignature, user.Height).Select("signature", "height").Updates(&user)
	err := results.Error
	if err == nil && results.RowsAffected < 1 {
		results = tx.Create(&user)
		err = results.Error
		isNew = true
	}

	if err == nil && len(secret.UserKey) > 0 && len(secret.UserSecretKey) > 0 {
		results = tx.Model(&secret).Where("user_id = ?", user.ID).Select("user_key", "user_secret_key").Updates(&secret)
		err = results.Error
		if err == nil && results.RowsAffected < 1 {
			results = tx.Create(&secret)
			err = results.Error
		}
	}

	if err == nil {
		userConsensusTableName := generateUserConsensusTableName(user.ID)
		tx_ := tx.Scopes(repository.TableName(userConsensusTableName))
		if isNew {
			err = tx_.Migrator().CreateTable(&UserConsensus{})
		} else {
			results = tx_.Delete(&UserConsensus{})
			err = results.Error
			if err == nil {
				results = tx.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name =?", userConsensusTableName)
				err = results.Error
			}
		}

		if err == nil {
			results = tx_.Create(userConsensuses)
			err = results.Error
		}
	}
	if err == nil {
		userDeviceTableName := generateUserDeviceTableName(user.ID)
		tx_ := tx.Scopes(repository.TableName(userDeviceTableName))
		if isNew {
			err = tx_.Migrator().CreateTable(&UserDevice{})
		} else {
			results = tx_.Delete(&UserDevice{})
			err = results.Error
			if err == nil {
				results = tx.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name =?", userDeviceTableName)
				err = results.Error
			}
		}

		if err == nil {
			results = tx_.Create(userDevices)
			err = results.Error
		}
	}
	if err == nil {
		userExtraTableName := generateUserExtraTableName(user.ID)
		tx_ := tx.Scopes(repository.TableName(userExtraTableName))
		if isNew {
			err = tx_.Migrator().CreateTable(&UserExtra{})
		} else {
			results = tx_.Delete(&UserExtra{})
			err = results.Error
			if err == nil {
				results = tx.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name =?", userExtraTableName)
				err = results.Error
			}
		}

		if err == nil {
			results = tx_.Create(userExtras)
			err = results.Error
		}

	}

	if err != nil {
		db.Rollback()
		return err
	}

	results = tx.Commit()
	return results.Error
}

func (repo *stdUserDataRepository) Clean(user User) error {
	db := repo.DB()
	if db == nil {
		return repository.ErrRepositoryDBUnavailable
	}

	tx := db.Begin()
	results := tx.Where("user_id = ? and height = ? and signature = ?", user.ID, user.Height, user.Signature).Delete(&User{})
	err := results.Error

	if err == nil {
		results = tx.Where("user_id = ?", user.ID).Delete(&UserSecret{})
		err = results.Error
	}

	if err == nil {
		userConsensusTableName := generateUserConsensusTableName(user.ID)
		tx_ := tx.Scopes(repository.TableName(userConsensusTableName))
		err = tx_.Migrator().DropTable(&UserConsensus{})
	}

	if err == nil {
		userDeviceTableName := generateUserDeviceTableName(user.ID)
		tx_ := tx.Scopes(repository.TableName(userDeviceTableName))
		err = tx_.Migrator().DropTable(&UserDevice{})
	}

	if err != nil {
		tx.Rollback()
	} else {
		results := tx.Commit()
		err = results.Error
	}
	return err
}
