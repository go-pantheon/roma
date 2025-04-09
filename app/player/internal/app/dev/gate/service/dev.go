package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/cmdregistrar"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
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
