package user

import (
	"context"
	"iter"
)

type UserService struct {
	UserRepository UserRepository

	UserConsensusService *UserConsensusService
}

func (service *UserService) Scan(ctx context.Context) (iter.Seq2[User, error], error) {
	return service.UserRepository.Scan(ctx)
}

func parseUserMetaWithUser(user User) UserMeta {
	var meta UserMeta
	meta.Code = user.Code
	meta.GenesisSignature = user.GenesisSignature
	meta.Signature = user.Signature
	meta.Height = user.Height
	return meta
}
