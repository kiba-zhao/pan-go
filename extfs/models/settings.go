package models

type Settings struct {
	TotalHeaderName string `mapstructure:"total_header_name"`
	DBFilePath      string `mapstructure:"db_file_path"`
}
