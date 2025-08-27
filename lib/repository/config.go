package repository

import "pan/lib/config"

type RepositoryConfig interface {
	BasePath() string
	TempPath() string
}

type RepositoryConfigListener = config.ConfigurerListener[RepositoryConfig]
type RepositoryConfigurer = config.Configurer[RepositoryConfig]
