package user

import (
	"context"
	"iter"
	"pan/pkg/ptp"
)

type UserExtraService struct {
	UserDataService     *UserDataService
	UserExtraRepository UserExtraRepository
}

func (service *UserExtraService) ScanWithUserMeta(ctx context.Context, peerId ptp.PeerID, meta UserMeta) (iter.Seq2[UserExtra, error], error) {
	user, err := service.UserDataService.CheckWithUserMetaForTopic(ctx, peerId, meta)
	if err != nil {
		return nil, err
	}

	return service.UserExtraRepository.ScanWithUserID(ctx, user.ID)
}
