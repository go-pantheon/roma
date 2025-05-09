package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
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

type UserProtoPool interface {
	Cache(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time)
	Load(ctx context.Context, uid int64, ctime time.Time) (ret *dbv1.UserProto)
	Get() (ret *dbv1.UserProto, putFunc func())
	Put(user *dbv1.UserProto)
}

type UserDomain struct {
	log       *log.Helper
	repo      UserRepo
	protoPool UserProtoPool
}

func NewUserDomain(pr UserRepo, logger log.Logger, cache UserProtoPool) *UserDomain {
	return &UserDomain{
		repo:      pr,
		protoPool: cache,
		log:       log.NewHelper(log.With(logger, "module", "player/user/gate/domain")),
	}
}

func (do *UserDomain) Create(ctx context.Context, uid int64, ctime time.Time, p *dbv1.UserProto) (err error) {
	defaultUser := do.initDefaultUser(uid)
	defaultUser.EncodeServer(p)
	err = do.repo.Create(ctx, uid, p, ctime)
	if err != nil {
		return err
	}
	return nil
}

func (do *UserDomain) initDefaultUser(id int64) *userobj.User {
	u := userobj.NewUser(id, fmt.Sprintf("player-%d", id))
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

func (do *UserDomain) CacheProto(ctx context.Context, uid int64, proto *dbv1.UserProto, ctime time.Time) {
	do.protoPool.Cache(ctx, uid, proto, ctime)
}

func (do *UserDomain) LoadCachedProto(ctx context.Context, uid int64, ctime time.Time) (ret *dbv1.UserProto) {
	return do.protoPool.Load(ctx, uid, ctime)
}

func (do *UserDomain) GetProto() (ret *dbv1.UserProto, putFunc func()) {
	return do.protoPool.Get()
}

func (do *UserDomain) PutBackProto(proto *dbv1.UserProto) {
	do.protoPool.Put(proto)
}
