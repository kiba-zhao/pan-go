package repository

func InitAsDefaults() error {
	err := initDBPathAsDefault()
	if err == nil {
		err = initTempDBPathAsDefault()
	}
	return err
}
