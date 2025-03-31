package config

type AppSettings = *Settings

type AppConfig = Config[AppSettings]

func New() AppConfig {
	settings := newDefaultSettings()
	cfg, err := NewConfig[AppSettings]("pan.toml")
	if err != nil {
		panic(err)
	}
	cfg.SetDefaults(settings)

	return cfg
}
