package config

func InitAsDefaults() error {
	err := InitHostName()
	if err == nil {
		err = InitRootPath()
	}
	return err
}
