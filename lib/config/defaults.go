package config

func InitAsDefaults() error {
	err := initHostNameAsDefault()
	if err == nil {
		err = initRootPathAsDefault()
	}
	return err
}
