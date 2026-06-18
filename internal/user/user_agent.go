package user

import (
	"bytes"
	"context"
	"pan/pkg/log"
	"pan/pkg/net"
	"slices"
	sync "sync"
)

type stdUserAgent struct {
	ordered      []net.PeerID
	pendings     []net.PeerID
	pendingsChan chan struct{}
	pendingsLock sync.RWMutex

	logger log.Logger
}

func (agent *stdUserAgent) ApplyConsistent(peerId net.PeerID) error {
	agent.pendingsLock.Lock()
	defer agent.pendingsLock.Unlock()

	idx, ok := slices.BinarySearchFunc(agent.ordered, peerId, bytes.Compare)
	if ok {
		return nil
	}

	agent.ordered = slices.Insert(agent.ordered, idx, peerId)
	agent.pendings = append(agent.pendings, peerId)
	if len(agent.pendings) == 1 {
		agent.pendingsChan <- struct{}{}
	}
	return nil
}

func (agent *stdUserAgent) MakeConsistent(ctx context.Context) error {
	agent.logger.Debug("UserAgent", "MakeConsistent begin")
	defer agent.logger.Debug("UserAgent", "MakeConsistent end")

	var err error
	var pendings []net.PeerID

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-agent.pendingsChan:
			agent.pendingsLock.Lock()
			pendings = agent.pendings
			agent.ordered = make([]net.PeerID, 0)
			agent.pendings = make([]net.PeerID, 0)
			agent.pendingsLock.Unlock()
		}

		if err != nil {
			break
		}

		for _, peerId := range pendings {
			println(peerId)
			// TODO: implement
		}
	}
	return err
}

// func (leader *stdLeader) MakeConsistent() {

// 	// TODO: implement

// 	if leader.UserDeviceService == nil {
// 		leader.logger.Error("user", "MakeConsistent Error:Leader Unavailable")
// 		return
// 	}

// 	userDeviceSeq, err := leader.UserDeviceService.ScanWithUserEnabled()
// 	if err != nil {
// 		leader.logger.Error("user", "MakeConsistent Error: "+err.Error())
// 		return
// 	}

// 	var peerId peer.PeerID
// 	remoteUsers := make([]*RemoteUserFields, 0)
// 	for device := range userDeviceSeq {
// 		peerId_, err := peer.DecodePeerID(device.PeerID)
// 		if err != nil {
// 			leader.logger.Error("user", "MakeConsistent Error: "+err.Error())
// 			continue
// 		}
// 		if !bytes.Equal(peerId_, peerId) {
// 			if len(remoteUsers) > 0 {
// 				leader.makeUserConsistent(peerId, remoteUsers)
// 				remoteUsers = make([]*RemoteUserFields, 0)
// 			}
// 			peerId = peerId_
// 		}

// 		var remoteUser RemoteUserFields
// 		remoteUser.Code = device.User.Code
// 		remoteUser.GenesisSignature = device.User.GenesisSignature
// 		remoteUsers = append(remoteUsers, &remoteUser)
// 	}

// 	if len(remoteUsers) > 0 {
// 		leader.makeConsistent(peerId, remoteUsers)
// 	}
// }
