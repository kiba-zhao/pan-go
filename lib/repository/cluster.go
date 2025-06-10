package repository

import (
	"path/filepath"
	"slices"
	"sync"
)

type Repository interface {
	SetupToRepository(db RepositoryDB) error
}

type RepositoryClusterModule interface {
	Repository
	DBName() string
}

type RepositoryCluster interface {
	AttachBaseModule(module RepositoryClusterModule) error
	DetachBaseModule(module RepositoryClusterModule) error

	AttachTempModule(module RepositoryClusterModule) error
	DetachTempModule(module RepositoryClusterModule) error
}

type stdRepositoryCluster struct {
	baseDBPath  string
	baseDBMap   map[string]RepositoryDB
	baseModules []RepositoryClusterModule
	baseLocker  sync.Mutex

	tempDBPath  string
	tempDBMap   map[string]RepositoryDB
	tempModules []RepositoryClusterModule
	tempLocker  sync.Mutex
}

var _ = (RepositoryCluster)((*stdRepositoryCluster)(nil))

func (cluster *stdRepositoryCluster) InitBaseDB(dbPath string) error {
	cluster.baseLocker.Lock()
	defer cluster.baseLocker.Unlock()
	if cluster.baseDBPath == dbPath {
		return nil
	}

	baseDBMap := make(map[string]RepositoryDB)
	if len(cluster.baseModules) > 0 {
		err := setupToModules(dbPath, baseDBMap, cluster.baseModules...)
		if err != nil {
			if len(cluster.baseDBPath) > 0 {
				setupToModules(cluster.baseDBPath, cluster.baseDBMap, cluster.baseModules...)
			}
			return err
		}
	}

	cluster.baseDBPath = dbPath
	cluster.baseDBMap = baseDBMap
	return nil
}

func (cluster *stdRepositoryCluster) AttachBaseModule(module RepositoryClusterModule) error {
	cluster.baseLocker.Lock()
	defer cluster.baseLocker.Unlock()

	var err error
	if len(cluster.baseDBPath) > 0 {
		err = setupToModules(cluster.baseDBPath, cluster.baseDBMap, module)
	}
	if err == nil {
		cluster.baseModules = append(cluster.baseModules, module)
	}
	return err
}

func (cluster *stdRepositoryCluster) DetachBaseModule(module RepositoryClusterModule) error {
	cluster.baseLocker.Lock()
	defer cluster.baseLocker.Unlock()

	dbName := module.DBName()
	keepDB := false
	offset := -1
	for idx, modulesitory := range cluster.baseModules {
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
		delete(cluster.baseDBMap, dbName)
	}

	if offset >= 0 {
		cluster.baseModules = slices.Delete(cluster.baseModules, offset, offset+1)
	}

	return module.SetupToRepository(nil)
}

func (cluster *stdRepositoryCluster) InitTempDB(dbPath string) error {
	cluster.tempLocker.Lock()
	defer cluster.tempLocker.Unlock()
	if cluster.tempDBPath == dbPath {
		return nil
	}

	tempDBMap := make(map[string]RepositoryDB)
	if len(cluster.tempModules) > 0 {
		err := setupToModules(dbPath, tempDBMap, cluster.tempModules...)
		if err != nil {
			if len(cluster.tempDBPath) > 0 {
				setupToModules(cluster.tempDBPath, cluster.tempDBMap, cluster.tempModules...)
			}
			return err
		}
	}

	cluster.tempDBPath = dbPath
	cluster.tempDBMap = tempDBMap
	return nil
}

func (cluster *stdRepositoryCluster) AttachTempModule(module RepositoryClusterModule) error {
	cluster.tempLocker.Lock()
	defer cluster.tempLocker.Unlock()

	var err error
	if len(cluster.tempDBPath) > 0 {
		err = setupToModules(cluster.tempDBPath, cluster.tempDBMap, module)
	}
	if err == nil {
		cluster.tempModules = append(cluster.tempModules, module)
	}
	return err
}

func (cluster *stdRepositoryCluster) DetachTempModule(module RepositoryClusterModule) error {
	cluster.tempLocker.Lock()
	defer cluster.tempLocker.Unlock()

	dbName := module.DBName()
	keepDB := false
	offset := -1
	for idx, modulesitory := range cluster.tempModules {
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
		delete(cluster.tempDBMap, dbName)
	}

	if offset >= 0 {
		cluster.tempModules = slices.Delete(cluster.tempModules, offset, offset+1)
	}

	return module.SetupToRepository(nil)
}

func setupToModules(dbPath string, dbMap map[string]RepositoryDB, modules ...RepositoryClusterModule) error {

	var err error
	for _, module := range modules {

		name := module.DBName()
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
