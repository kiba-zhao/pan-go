//go:build !(android || ios)

package settings

import "pan/pkg/web"

type stdWebConfig struct {
	settings *HostSettings
}

func newWebConfig(settings *HostSettings) web.WebConfig {
	cfg := &stdWebConfig{}
	cfg.settings = settings
	return cfg
}

var _ = (web.WebConfig)((*stdWebConfig)(nil))

func (cfg *stdWebConfig) Addr() string {
	if !cfg.settings.WebEnabled {
		return ""
	}
	return cfg.settings.WebAddr
}
