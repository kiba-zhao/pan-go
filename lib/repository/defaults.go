package repository

func InitAsDefaults() error {
	err := InitDBPath()
	if err == nil {
		err = InitTempDBPath()
	}
	return err
}
