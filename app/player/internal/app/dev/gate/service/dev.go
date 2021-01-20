package service

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds/cmdregistrar"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-kit/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DevService struct {
	climsg.UnimplementedDevServiceServer

	log *log.Helper
	uc  *biz.DevUseCase
}

func NewDevService(logger log.Logger, uc *biz.DevUseCase, rg *cmdregistrar.Registrar) climsg.DevServiceServer {
	return &DevService{
		log: log.NewHelper(log.With(logger, "module", "player/dev/gate/service")),
		uc:  uc,
	}
}

func (s *DevService) DevList(ctx context.Context, cs *climsg.CSDevList) (*climsg.SCDevList, error) {
	if !profile.IsDev() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot list dev commands in non-dev environment")
	}
	return s.uc.List(ctx.(core.Context))
}

func (s *DevService) DevExecute(ctx context.Context, cs *climsg.CSDevExecute) (*climsg.SCDevExecute, error) {
	if !profile.IsDev() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot execute command in non-dev environment")
	}

	return s.uc.Execute(ctx.(core.Context), cs.Mod, cs.Cmd, cs.Args)
}
