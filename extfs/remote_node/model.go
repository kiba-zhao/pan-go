package remotenode

import "time"

type RemoteNode struct {
	ID                 uint      `json:"id" form:"id"`
	PeerID             string    `json:"peerId" form:"peerId"`
	Name               string    `json:"name" form:"name"`
	Available          bool      `json:"available" form:"available"`
	TagQuantity        uint      `json:"tagQuantity" form:"tagQuantity"`
	PendingTagQuantity uint      `json:"pendingTagQuantity" form:"pendingTagQuantity"`
	CreatedAt          time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" form:"updatedAt"`
}
