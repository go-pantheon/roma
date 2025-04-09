package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

type UserService struct {
	climsg.UnimplementedUserServiceServer
	log *log.Helper

	uc *biz.UserUseCase
}

func NewUserService(logger log.Logger, uc *biz.UserUseCase) climsg.UserServiceServer {
	return &UserService{
		log: log.NewHelper(log.With(logger, "module", "player/user/gate/service")),
		uc:  uc,
	}
}

func (s *UserService) Login(ctx context.Context, cs *climsg.CSLogin) (*climsg.SCLogin, error) {
	return s.uc.Login(ctx.(core.Context), cs)
}

func (s *UserService) UpdateName(ctx context.Context, cs *climsg.CSUpdateName) (*climsg.SCUpdateName, error) {
	return s.uc.UpdateName(ctx.(core.Context), cs)
}

func (s *UserService) SetGender(ctx context.Context, cs *climsg.CSSetGender) (*climsg.SCSetGender, error) {
	return s.uc.SetGender(ctx.(core.Context), cs)
}
