package domain

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	userobj "github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

type UserRepo interface {
	Create(ctx context.Context, uid int64, defaultUser *dbv1.UserProto, ctime time.Time) error
	QueryByID(ctx context.Context, uid int64, p *dbv1.UserProto, mods []life.ModuleKey) error
	UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error
	IsExist(ctx context.Context, uid int64) (bool, error)
	IncVersion(ctx context.Context, uid int64, newVersion int64) error
}

type UserCache interface {
	CachePut(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time)
	CacheGet(ctx context.Context, uid int64, ctime time.Time) (ret *dbv1.UserProto)
	CacheRemove(ctx context.Context, uid int64)
}

type UserDomain struct {
	log        *log.Helper
	repo       UserRepo
	protoCache UserCache
}

func NewUserDomain(pr UserRepo, logger log.Logger, cache UserCache) *UserDomain {
	return &UserDomain{
		repo:       pr,
		protoCache: cache,
		log:        log.NewHelper(log.With(logger, "module", "player/user/gate/domain")),
	}
}

func (do *UserDomain) Create(ctx context.Context, uid int64, ctime time.Time, p *dbv1.UserProto) (err error) {
	defaultUser := do.initDefaultUser(uid)
	defaultUser.EncodeServer(p, userregister.AllModuleKeys())

	err = do.repo.Create(ctx, uid, p, ctime)
	if err != nil {
		return err
	}
	return nil
}

func (do *UserDomain) initDefaultUser(id int64) *userobj.User {
	u := userobj.NewUser(id, profile.Version())
	return u
}

func (do *UserDomain) Exist(ctx context.Context, uid int64) (bool, error) {
	return do.repo.IsExist(ctx, uid)
}

func (do *UserDomain) Load(ctx context.Context, uid int64, p *dbv1.UserProto) (err error) {
	return do.repo.QueryByID(ctx, uid, p, userregister.AllModuleKeys())
}

func (do *UserDomain) Persist(ctx context.Context, uid int64, p *dbv1.UserProto) (err error) {
	defer dbv1.UserProtoPool.Put(p)

	return do.repo.UpdateByID(ctx, uid, p)
}

func (do *UserDomain) IncVersion(ctx context.Context, uid int64, newVersion int64) (err error) {
	return do.repo.IncVersion(ctx, uid, newVersion)
}

func (do *UserDomain) CachePut(ctx context.Context, uid int64, proto *dbv1.UserProto, ctime time.Time) {
	do.protoCache.CachePut(ctx, uid, proto, ctime)
}

func (do *UserDomain) CacheGet(ctx context.Context, uid int64, ctime time.Time) (ret *dbv1.UserProto) {
	return do.protoCache.CacheGet(ctx, uid, ctime)
}

func (do *UserDomain) CacheRemove(ctx context.Context, uid int64) {
	do.protoCache.CacheRemove(ctx, uid)
}
