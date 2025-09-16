package user

import (
	"pan/lib/feature"
	"pan/lib/repository"
)

var subModuleNewFuncArray []feature.SubModuleNewFunc[*stdModule]

const (
	ModuleName = "user"
)

func New() interface{} {

	module := &stdModule{}
	module.FeatureModule = feature.New(module, subModuleNewFuncArray...)

	return module
}

type stdModule struct {
	*feature.FeatureModule
}

var _ = (repository.RepositoryDBModule)((*stdModule)(nil))

func (module *stdModule) SetupToRepository(db repository.RepositoryDB) error {
	err := db.AutoMigrate(
		&User{},
		&UserDevice{},
		&OwnerSecret{},
		&OwnerHistory{},
	)
	return err
}

func (module *stdModule) DBName() string {
	return ModuleName + ".db"
}
