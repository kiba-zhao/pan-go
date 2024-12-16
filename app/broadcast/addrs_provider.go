package broadcast

import (
	"pan/app/config"
	"slices"
	"sync"
)

type broadcastAddrsProvider struct {
	module        *broadcastModule
	agent         *broadcastAgent
	addrs         []string
	addrsRW       sync.RWMutex
	publicAddrs   []string
	publicAddrsRW sync.RWMutex
}

func (provider *broadcastAddrsProvider) BroadcastServeAddrs() []string {
	provider.addrsRW.RLock()
	defer provider.addrsRW.RUnlock()
	return provider.addrs
}

func (provider *broadcastAddrsProvider) BroadcastDeliverAddrs() []string {
	return provider.BroadcastServeAddrs()
}

func (provider *broadcastAddrsProvider) BroadcastPublicAddrs() []string {
	provider.publicAddrsRW.RLock()
	defer provider.publicAddrsRW.RUnlock()
	return provider.publicAddrs
}

func (provider *broadcastAddrsProvider) OnConfigUpdated(settings config.AppSettings) {

	provider.addrsRW.Lock()
	if !slices.Equal(provider.addrs, settings.BroadcastAddress) {
		provider.addrs = settings.BroadcastAddress
		provider.module.Reload()
	}
	provider.addrsRW.Unlock()

	provider.publicAddrsRW.Lock()
	if !slices.Equal(provider.publicAddrs, settings.PublicAddress) {
		provider.publicAddrs = settings.PublicAddress
		provider.agent.Reload()
	}
	provider.publicAddrsRW.Unlock()

}
