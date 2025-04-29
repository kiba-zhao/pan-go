package config

type AppSettings = *Settings

type AppConfig = Config[AppSettings]

func New() AppConfig {
	cfg, err := NewConfig[AppSettings]("pan.toml")
	if err != nil {
		panic(err)
	}
	settings := newDefaultSettings(cfg)
	cfg.SetDefaults(settings)

	return cfg
}
