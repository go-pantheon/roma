package domain

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

type UserRepo interface {
	GetByID(ctx context.Context, user *dbv1.UserProto, mods []life.ModuleKey) error
	GetList(ctx context.Context, start, limit int64, conds map[life.ModuleKey]*dbv1.UserModuleProto, mods []life.ModuleKey) ([]*dbv1.UserProto, int64, error)
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

func (do *UserDomain) GetByID(ctx context.Context, user *dbv1.UserProto, mods []life.ModuleKey) (err error) {
	return do.repo.GetByID(ctx, user, mods)
}

func (do *UserDomain) GetList(ctx context.Context, start, limit int64, conds map[life.ModuleKey]*dbv1.UserModuleProto, mods []life.ModuleKey) ([]*dbv1.UserProto, int64, error) {
	return do.repo.GetList(ctx, start, limit, conds, mods)
}
