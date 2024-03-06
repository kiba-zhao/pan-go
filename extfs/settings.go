package extfs

// TODO: defines a struct named Settings with one fields: Enabled of type bool

type Settings struct {
	Enabled bool
}

func defaultSettings() *Settings {
	return &Settings{
		Enabled: true,
	}
}
