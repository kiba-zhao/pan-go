package broadcast

import (
	"net"
	"pan/lib/log"
	"pan/lib/peer"
	"pan/lib/quic"
	"slices"
	"sync"
)

type stdBroadcastRuntime struct {
	logger log.Logger

	serveModules   []BroadcastServeModule
	serveModulesRW sync.RWMutex

	addrs   []string
	addrsRW sync.RWMutex

	deliverLimitSize   int
	deliverLimitSizeRW sync.RWMutex

	mtu   int
	mtuRW sync.RWMutex

	quicCluster   quic.QuicCluster
	quicClusterRW sync.RWMutex

	peerSettings   *peer.PeerSettings
	peerSettingsRW sync.RWMutex

	store   BroadcastStore
	storeRW sync.RWMutex

	server *stdBroadcastServer
	agent  *stdBroadcastAgent
}

var _ = (BroadcastServerRuntime)((*stdBroadcastRuntime)(nil))

func (runtime *stdBroadcastRuntime) Addrs() []string {
	runtime.addrsRW.RLock()
	defer runtime.addrsRW.RUnlock()

	if len(runtime.addrs) <= 0 {
		return runtime.addrs
	}

	addrs := make([]string, 0)
	for _, addr := range runtime.addrs {
		if idx, ok := slices.BinarySearch(addrs, addr); !ok {
			addrs = slices.Insert(addrs, idx, addr)
		}
	}
	return addrs
}

func (runtime *stdBroadcastRuntime) setAddrs(addrs []string) {
	runtime.logger.Debug("BroadcastRuntime", "SetAddrs")

	runtime.addrsRW.Lock()
	defer runtime.addrsRW.Unlock()
	if slices.Equal(runtime.addrs, addrs) {
		return
	}
	runtime.addrs = addrs
	server := runtime.server
	if server != nil {
		server.Reload()
	}
}

var _ = (BroadcastClusterRuntime)((*stdBroadcastRuntime)(nil))

func (runtime *stdBroadcastRuntime) ServeModules() []BroadcastServeModule {
	runtime.serveModulesRW.RLock()
	defer runtime.serveModulesRW.RUnlock()
	return runtime.serveModules
}

func (runtime *stdBroadcastRuntime) RegisterServeModule(module BroadcastServeModule) {
	runtime.serveModulesRW.Lock()
	defer runtime.serveModulesRW.Unlock()
	runtime.serveModules = append(runtime.serveModules, module)
}

func (runtime *stdBroadcastRuntime) UnregisterServeModule(module BroadcastServeModule) {
	runtime.serveModulesRW.Lock()
	defer runtime.serveModulesRW.Unlock()
	for i, m := range runtime.serveModules {
		if m == module {
			runtime.serveModules = append(runtime.serveModules[:i], runtime.serveModules[i+1:]...)
			break
		}
	}
}

func (runtime *stdBroadcastRuntime) DeliverLimitSize() int {
	runtime.deliverLimitSizeRW.RLock()
	defer runtime.deliverLimitSizeRW.RUnlock()
	return runtime.deliverLimitSize
}

func (runtime *stdBroadcastRuntime) setDeliverLimitSize(size int) {
	runtime.deliverLimitSizeRW.Lock()
	defer runtime.deliverLimitSizeRW.Unlock()
	runtime.deliverLimitSize = size
}

func (runtime *stdBroadcastRuntime) MTU() int {
	runtime.mtuRW.RLock()
	defer runtime.mtuRW.RUnlock()
	return runtime.mtu
}

func (runtime *stdBroadcastRuntime) setMTU(mtu int) {
	runtime.mtuRW.Lock()
	defer runtime.mtuRW.Unlock()
	runtime.mtu = mtu
}

func (runtime *stdBroadcastRuntime) DeliverConn() *net.UDPConn {
	quicCluster := runtime.QuicCluster()
	if quicCluster == nil {
		return nil
	}
	return quicCluster.ServeUDPConn()
}

var _ = (BroadcastAgentRuntime)((*stdBroadcastRuntime)(nil))

func (runtime *stdBroadcastRuntime) QuicCluster() quic.QuicCluster {
	runtime.quicClusterRW.RLock()
	defer runtime.quicClusterRW.RUnlock()
	return runtime.quicCluster
}

func (runtime *stdBroadcastRuntime) setQuicCluster(cluster quic.QuicCluster) {
	runtime.quicClusterRW.Lock()
	defer runtime.quicClusterRW.Unlock()
	runtime.quicCluster = cluster
}

func (runtime *stdBroadcastRuntime) PeerSettings() *peer.PeerSettings {
	runtime.peerSettingsRW.RLock()
	defer runtime.peerSettingsRW.RUnlock()
	return runtime.peerSettings
}

func (runtime *stdBroadcastRuntime) setPeerSettings(settings *peer.PeerSettings) {
	runtime.peerSettingsRW.Lock()
	defer runtime.peerSettingsRW.Unlock()
	runtime.peerSettings = settings

	agent := runtime.agent
	if agent != nil {
		agent.Reload()
	}

}

func (runtime *stdBroadcastRuntime) Store() BroadcastStore {
	runtime.storeRW.RLock()
	defer runtime.storeRW.RUnlock()
	return runtime.store
}

func (runtime *stdBroadcastRuntime) setStore(store BroadcastStore) {
	runtime.storeRW.Lock()
	defer runtime.storeRW.Unlock()
	runtime.store = store

	agent := runtime.agent
	if agent != nil {
		agent.Reload()
	}
}
