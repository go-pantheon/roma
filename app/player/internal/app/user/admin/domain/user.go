package domain

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

type UserRepo interface {
	GetByID(ctx context.Context, uid int64) (*dbv1.UserProto, error)
	GetList(ctx context.Context, start, limit int64, cond *dbv1.UserProto) ([]*dbv1.UserProto, int64, error)
	UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error
}

type UserDomain struct {
	log  *log.Helper
	repo UserRepo
}

func NewUserDomain(pr UserRepo, logger log.Logger) *UserDomain {
	return &UserDomain{
		repo: pr,
		log:  log.NewHelper(log.With(logger, "module", "player/user/admin/domain")),
	}
}

func (do *UserDomain) GetByID(ctx context.Context, uid int64) (u *dbv1.UserProto, err error) {
	return do.repo.GetByID(ctx, uid)
}

func (do *UserDomain) GetList(ctx context.Context, start, limit int64, cond *dbv1.UserProto) ([]*dbv1.UserProto, int64, error) {
	return do.repo.GetList(ctx, start, limit, cond)
}

func (do *UserDomain) UpdateByID(ctx context.Context, uid int64, proto *dbv1.UserProto) (err error) {
	return do.repo.UpdateByID(ctx, uid, proto)
}
