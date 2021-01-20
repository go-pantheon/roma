package service

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/system/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
)

type SystemService struct {
	*climsg.UnimplementedSystemServiceServer

	uc  *biz.SystemUseCase
	log *log.Helper
}

func NewSystemService(logger log.Logger, uc *biz.SystemUseCase) climsg.SystemServiceServer {
	return &SystemService{
		log: log.NewHelper(log.With(logger, "module", "player/system/gate/service")),
		uc:  uc,
	}
}

func (s *SystemService) Heartbeat(ctx context.Context, cs *climsg.CSHeartbeat) (*climsg.SCHeartbeat, error) {
	return s.uc.Heartbeat(ctx.(core.Context), cs)
}
