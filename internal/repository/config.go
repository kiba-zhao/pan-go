package repository

import "pan/internal/config"

type RepositoryConfig interface {
	BasePath() string
	TempPath() string
}

type RepositoryConfigListener = config.ConfigurerListener[RepositoryConfig]
type RepositoryConfigurer = config.Configurer[RepositoryConfig]
