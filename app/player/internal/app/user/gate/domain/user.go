package domain

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	userobj "github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/object"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"

	"github.com/pkg/errors"
)

type UserRepo interface {
	Create(ctx context.Context, uid int64, defaultUser *dbv1.UserProto, ctime time.Time) error
	QueryByID(ctx context.Context, uid int64, p *dbv1.UserProto) error
	UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error
	UpdateLoginTime(ctx context.Context, uid int64, loginAt, logoutAt time.Time) error
	Exist(ctx context.Context, uid int64) (bool, error)
	IncVersion(ctx context.Context, uid int64, newVersion int64) error
}

type UserProtoCache interface {
	Put(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time)
	Get(ctx context.Context, uid int64, ctime time.Time) (ret *dbv1.UserProto)
	Remove(ctx context.Context, uid int64)
}

type UserDomain struct {
	log        *log.Helper
	repo       UserRepo
	protoCache UserProtoCache
}

func NewUserDomain(pr UserRepo, logger log.Logger, cache UserProtoCache) *UserDomain {
	return &UserDomain{
		repo:       pr,
		protoCache: cache,
		log:        log.NewHelper(log.With(logger, "module", "player/user/gate/domain")),
	}
}

func (do *UserDomain) Create(ctx context.Context, uid int64, sid int64, ctime time.Time, p *dbv1.UserProto) (err error) {
	defaultUser := do.initDefaultUser(uid, sid)
	defaultUser.EncodeServer(p, nil, true)
	err = do.repo.Create(ctx, uid, p, ctime)
	if err != nil {
		return err
	}
	return nil
}

func (do *UserDomain) initDefaultUser(id int64, sid int64) *userobj.User {
	u := userobj.NewUser(id, sid, profile.Version())
	return u
}

func (do *UserDomain) Exist(ctx context.Context, uid int64) (bool, error) {
	return do.repo.Exist(ctx, uid)
}

func (do *UserDomain) Load(ctx context.Context, uid int64, p *dbv1.UserProto) (err error) {
	return do.repo.QueryByID(ctx, uid, p)
}

func (do *UserDomain) Login(ctx context.Context, uid int64, ctime, logoutTime time.Time) (err error) {
	return do.repo.UpdateLoginTime(ctx, uid, ctime, logoutTime)
}

func (do *UserDomain) Persist(ctx context.Context, uid int64, proto life.VersionProto) (err error) {
	p, ok := proto.(*dbv1.UserProto)
	if !ok {
		err = errors.Wrapf(xerrors.ErrDBRecordType, "uid=%d proto=%T", uid, proto)
		return
	}
	return do.repo.UpdateByID(ctx, uid, p)
}

func (do *UserDomain) IncVersion(ctx context.Context, uid int64, newVersion int64) (err error) {
	return do.repo.IncVersion(ctx, uid, newVersion)
}

func (do *UserDomain) Cache(ctx context.Context, uid int64, proto *dbv1.UserProto, ctime time.Time) {
	do.protoCache.Put(ctx, uid, proto, ctime)
}

func (do *UserDomain) Get(ctx context.Context, uid int64, ctime time.Time) (ret *dbv1.UserProto) {
	return do.protoCache.Get(ctx, uid, ctime)
}

func (do *UserDomain) Remove(ctx context.Context, uid int64) {
	do.protoCache.Remove(ctx, uid)
}
