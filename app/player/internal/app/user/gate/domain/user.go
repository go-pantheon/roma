package domain

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/errors"
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
	UpdateSID(ctx context.Context, uid int64, sid int64, version int64) error
	IncVersion(ctx context.Context, uid int64, newVersion int64) error
}

type UserCache interface {
	Put(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time)
	Get(ctx context.Context, uid int64, ctime time.Time) (ret *dbv1.UserProto)
	Remove(ctx context.Context, uid int64)
}

type UserDomain struct {
	log   *log.Helper
	repo  UserRepo
	cache UserCache
}

func NewUserDomain(pr UserRepo, logger log.Logger, cache UserCache) *UserDomain {
	return &UserDomain{
		log:   log.NewHelper(log.With(logger, "module", "player/user/gate/domain")),
		repo:  pr,
		cache: cache,
	}
}

func (do *UserDomain) Create(ctx context.Context, uid int64, ctime time.Time, p *dbv1.UserProto) (err error) {
	defaultUser := do.newDefaultUser(uid)
	defaultUser.EncodeServer(p, userregister.AllModuleKeys())

	err = do.repo.Create(ctx, uid, p, ctime)
	if err != nil {
		return err
	}
	return nil
}

func (do *UserDomain) newDefaultUser(id int64) *userobj.User {
	return userobj.NewUser(id, profile.Version())
}

func (do *UserDomain) Exist(ctx context.Context, uid int64) (bool, error) {
	return do.repo.IsExist(ctx, uid)
}

func (do *UserDomain) UpdateSID(ctx context.Context, uid int64, sid int64, version int64) (err error) {
	return do.repo.UpdateSID(ctx, uid, sid, version)
}

func (do *UserDomain) TakeUserProto(ctx context.Context, uid int64, allowCreate bool) (ret *dbv1.UserProto, newborn bool, err error) {
	p := do.cache.Get(ctx, uid, time.Now())
	if p != nil {
		do.cache.Remove(ctx, uid)
		return p, false, nil
	}

	p = dbv1.UserProtoPool.Get()

	p.Id = uid

	if err = do.Load(ctx, uid, p); err != nil {
		if errors.Is(err, xerrors.ErrDBRecordNotFound) {
			if allowCreate {
				err = do.Create(ctx, uid, time.Now(), p)
				newborn = true
			}
		}
	}

	if err != nil {
		return nil, false, err
	}

	return p, newborn, nil
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

func (do *UserDomain) OnLogout(ctx context.Context, uid int64, proto *dbv1.UserProto) {
	do.cache.Put(ctx, uid, proto, time.Now())
}
