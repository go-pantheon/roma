package service

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
)

type StorageService struct {
	climsg.UnimplementedStorageServiceServer

	uc  *biz.StorageUseCase
	log *log.Helper
}

func NewStorageService(logger log.Logger, uc *biz.StorageUseCase) climsg.StorageServiceServer {
	return &StorageService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "player/storage/gate/service")),
	}
}

func (s *StorageService) UsePack(ctx context.Context, req *climsg.CSUsePack) (*climsg.SCUsePack, error) {
	return s.uc.UsePack(ctx.(core.Context), req)
}
