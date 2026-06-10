package user

import (
	"iter"
	"pan/internal/net"
)

type UserExtraService struct {
	UserDataService     *UserDataService
	UserExtraRepository UserExtraRepository
}

func (service *UserExtraService) ScanWithUserMeta(peerId net.PeerID, meta UserMeta) (iter.Seq2[UserExtra, error], error) {
	user, err := service.UserDataService.CheckWithUserMetaForTopic(peerId, meta)
	if err != nil {
		return nil, err
	}

	return service.UserExtraRepository.ScanWithUserID(user.ID)
}
