package domain

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/pkg/errors"
)

type RoomRepo interface {
	Create(ctx context.Context, p *dbv1.RoomProto, ctime time.Time) (err error)
	QueryByID(ctx context.Context, id int64, p *dbv1.RoomProto) error
	UpdateByID(ctx context.Context, id int64, room *dbv1.RoomProto) error
	Exist(ctx context.Context, id int64) (bool, error)
	IncVersion(ctx context.Context, id int64, newVersion int64) error
}

type RoomProtoCache interface {
	Put(ctx context.Context, uid int64, room *dbv1.RoomProto, ctime time.Time)
	Get(ctx context.Context, uid int64, ctime time.Time) *dbv1.RoomProto
	Remove(ctx context.Context, uid int64)
}

type RoomDomain struct {
	log   *log.Helper
	repo  RoomRepo
	cache RoomProtoCache // offline object cache
}

func NewRoomDomain(pr RoomRepo, logger log.Logger, cache RoomProtoCache) *RoomDomain {
	return &RoomDomain{
		repo:  pr,
		cache: cache,
		log:   log.NewHelper(log.With(logger, "module", "room/room/gate/domain"))}
}

func (do *RoomDomain) Create(ctx context.Context, _ int64, ctime time.Time, p *dbv1.RoomProto) (err error) {
	return do.repo.Create(ctx, p, ctime)
}

func (do *RoomDomain) Load(ctx context.Context, id int64, p *dbv1.RoomProto) (err error) {
	return do.repo.QueryByID(ctx, id, p)
}

func (do *RoomDomain) TakeRoomProto(ctx context.Context, id int64, allowCreate bool) (ret *dbv1.RoomProto, newborn bool, err error) {
	p := do.cache.Get(ctx, id, time.Now())
	if p != nil {
		do.cache.Remove(ctx, id)
		return p, false, nil
	}

	p = dbv1.RoomProtoPool.Get()

	p.Id = id

	if err = do.Load(ctx, id, p); err != nil {
		if errors.Is(err, xerrors.ErrDBRecordNotFound) {
			if allowCreate {
				err = do.Create(ctx, id, time.Now(), p)
				newborn = true
			}
		}
	}

	if err != nil {
		return nil, false, err
	}

	return p, newborn, nil
}

func (do *RoomDomain) Persist(ctx context.Context, id int64, proto life.VersionProto) (err error) {
	data, ok := proto.(*dbv1.RoomProto)
	if !ok {
		err = errors.Wrapf(xerrors.ErrDBProtoEncode, "room domain persist failed. id=%d proto=%T", id, proto)
		return
	}

	return do.repo.UpdateByID(ctx, id, data)
}

func (do *RoomDomain) IncVersion(ctx context.Context, id int64, newVersion int64) (err error) {
	return do.repo.IncVersion(ctx, id, newVersion)
}

func (do *RoomDomain) Exist(ctx context.Context, id int64) (exist bool, err error) {
	return do.repo.Exist(ctx, id)
}

func (do *RoomDomain) OnLogout(ctx context.Context, id int64, proto *dbv1.RoomProto) error {
	do.cache.Put(ctx, id, proto, time.Now())
	return nil
}
