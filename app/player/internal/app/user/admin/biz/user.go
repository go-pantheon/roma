package biz

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/admin/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
)

type UserUseCase struct {
	log *log.Helper
	mgr *core.Manager
	do  *domain.UserDomain
}

func NewUserUseCase(mgr *core.Manager, do *domain.UserDomain, logger log.Logger) *UserUseCase {
	uc := &UserUseCase{
		log: log.NewHelper(log.With(logger, "module", "player/user/admin/biz")),
		mgr: mgr,
		do:  do,
	}

	return uc
}

func (uc *UserUseCase) GetByID(ctx context.Context, uid int64) (p *dbv1.UserProto, err error) {
	p, err = uc.do.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (uc *UserUseCase) GetList(ctx context.Context, index, size int32, cond *dbv1.UserProto) ([]*dbv1.UserProto, int64, error) {
	return uc.do.GetList(ctx, index, size, cond)
}

func (uc *UserUseCase) UpdateUserByID(ctx context.Context, uid int64, proto *dbv1.UserProto) (err error) {
	return uc.do.UpdateByID(ctx, uid, proto)
}
