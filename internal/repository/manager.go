package repository

import (
	"pan/internal/log"
	"path/filepath"
	"slices"
	"sync"
)

type Repository interface {
	SetupToRepository(db RepositoryDB) error
}

type RepositoryDBModule interface {
	Repository
	DBName() string
}

type RepositoryManager interface {
	AttachBaseModule(module RepositoryDBModule) error
	DetachBaseModule(module RepositoryDBModule) error

	AttachTempModule(module RepositoryDBModule) error
	DetachTempModule(module RepositoryDBModule) error
}

type stdRepositoryManager struct {
	logger log.Logger

	locker sync.Mutex

	baseDBPath  string
	baseDBMap   map[string]RepositoryDB
	baseModules []RepositoryDBModule
	baseLocker  sync.Mutex

	tempDBPath  string
	tempDBMap   map[string]RepositoryDB
	tempModules []RepositoryDBModule
	tempLocker  sync.Mutex
}

var _ = (RepositoryManager)((*stdRepositoryManager)(nil))

func (manager *stdRepositoryManager) Setup(config RepositoryConfig) error {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	manager.logger.Debug("RepositoryManager", "Setup")
	basePath := config.BasePath()
	tempPath := config.TempPath()
	err := manager.InitBaseDB(basePath)
	if err == nil {
		err = manager.InitTempDB(tempPath)
	}
	return err
}

func (manager *stdRepositoryManager) InitBaseDB(dbPath string) error {
	manager.baseLocker.Lock()
	defer manager.baseLocker.Unlock()
	if manager.baseDBPath == dbPath {
		return nil
	}

	baseDBMap := make(map[string]RepositoryDB)
	if len(manager.baseModules) > 0 {
		err := setupToModules(dbPath, baseDBMap, manager.baseModules...)
		if err != nil {
			if len(manager.baseDBPath) > 0 {
				setupToModules(manager.baseDBPath, manager.baseDBMap, manager.baseModules...)
			}
			return err
		}
	}

	manager.baseDBPath = dbPath
	manager.baseDBMap = baseDBMap
	return nil
}

func (manager *stdRepositoryManager) InitTempDB(dbPath string) error {
	manager.tempLocker.Lock()
	defer manager.tempLocker.Unlock()
	if manager.tempDBPath == dbPath {
		return nil
	}

	tempDBMap := make(map[string]RepositoryDB)
	if len(manager.tempModules) > 0 {
		err := setupToModules(dbPath, tempDBMap, manager.tempModules...)
		if err != nil {
			if len(manager.tempDBPath) > 0 {
				setupToModules(manager.tempDBPath, manager.tempDBMap, manager.tempModules...)
			}
			return err
		}
	}

	manager.tempDBPath = dbPath
	manager.tempDBMap = tempDBMap
	return nil
}

func (manager *stdRepositoryManager) AttachBaseModule(module RepositoryDBModule) error {
	manager.baseLocker.Lock()
	defer manager.baseLocker.Unlock()

	var err error
	if len(manager.baseDBPath) > 0 {
		err = setupToModules(manager.baseDBPath, manager.baseDBMap, module)
	}
	if err == nil {
		manager.baseModules = append(manager.baseModules, module)
	}
	return err
}

func (manager *stdRepositoryManager) DetachBaseModule(module RepositoryDBModule) error {
	manager.baseLocker.Lock()
	defer manager.baseLocker.Unlock()

	dbName := module.DBName()
	keepDB := false
	offset := -1
	for idx, modulesitory := range manager.baseModules {
		if modulesitory == module {
			offset = idx
			if keepDB {
				break
			}
			continue
		}

		if modulesitory.DBName() == dbName {
			keepDB = true
			if offset >= 0 {
				break
			}
		}
	}

	if !keepDB {
		delete(manager.baseDBMap, dbName)
	}

	if offset >= 0 {
		manager.baseModules = slices.Delete(manager.baseModules, offset, offset+1)
	}

	return module.SetupToRepository(nil)
}

func (manager *stdRepositoryManager) AttachTempModule(module RepositoryDBModule) error {
	manager.tempLocker.Lock()
	defer manager.tempLocker.Unlock()

	var err error
	if len(manager.tempDBPath) > 0 {
		err = setupToModules(manager.tempDBPath, manager.tempDBMap, module)
	}
	if err == nil {
		manager.tempModules = append(manager.tempModules, module)
	}
	return err
}

func (manager *stdRepositoryManager) DetachTempModule(module RepositoryDBModule) error {
	manager.tempLocker.Lock()
	defer manager.tempLocker.Unlock()

	dbName := module.DBName()
	keepDB := false
	offset := -1
	for idx, modulesitory := range manager.tempModules {
		if modulesitory == module {
			offset = idx
			if keepDB {
				break
			}
			continue
		}

		if modulesitory.DBName() == dbName {
			keepDB = true
			if offset >= 0 {
				break
			}
		}
	}

	if !keepDB {
		delete(manager.tempDBMap, dbName)
	}

	if offset >= 0 {
		manager.tempModules = slices.Delete(manager.tempModules, offset, offset+1)
	}

	return module.SetupToRepository(nil)
}

func setupToModules(dbPath string, dbMap map[string]RepositoryDB, modules ...RepositoryDBModule) error {

	var err error
	for _, module := range modules {

		name := module.DBName()
		if len(name) <= 0 {
			err = module.SetupToRepository(nil)
		}

		if err != nil {
			break
		}

		db, ok := dbMap[name]
		if !ok {
			db, err = newSqliteDB(filepath.Join(dbPath, name))
			if err != nil {
				break
			}
			dbMap[name] = db
		}

		err = module.SetupToRepository(db)
		if err != nil {
			break
		}
	}

	return err
}
