// Define broadcast addrs provider
//
// It is used to provide the addresses that the broadcast module can serve and deliver broadcast messages.
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

// BroadcastServeAddrs returns the list of addresses that the broadcast module should serve.
//
// The addresses are used to serve the broadcast message when the Serve method is called.
// If the parameter list is empty, the message is delivered to all online peers.
func (provider *broadcastAddrsProvider) BroadcastServeAddrs() []string {
	provider.addrsRW.RLock()
	defer provider.addrsRW.RUnlock()
	return provider.addrs
}

// BroadcastDeliverAddrs returns the list of addresses that the broadcast module can deliver.
//
// The addresses are used to deliver the broadcast message when the Deliver method is called.
// If the parameter list is empty, the message is delivered to all online peers.
func (provider *broadcastAddrsProvider) BroadcastDeliverAddrs() []string {
	return provider.BroadcastServeAddrs()
}

// BroadcastPublicAddrs returns the list of public addresses that the broadcast module can use.
// The addresses are used to broadcast messages publicly.

func (provider *broadcastAddrsProvider) BroadcastPublicAddrs() []string {
	provider.publicAddrsRW.RLock()
	defer provider.publicAddrsRW.RUnlock()
	return provider.publicAddrs
}

// OnConfigUpdated updates the broadcast and public addresses of the provider
// based on the given application settings. If the addresses are changed,
// it triggers a reload of the respective module and agent components to
// apply the new configuration. Locks are used to ensure thread safety
// while updating the addresses.

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
