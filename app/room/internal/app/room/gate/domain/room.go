package domain

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/log"
	"github.com/pkg/errors"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/room/v1"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/life"
	"github.com/vulcan-frame/vulcan-kit/xerrors"
)

type RoomRepo interface {
	Create(ctx context.Context, p *dbv1.RoomProto, ctime time.Time) (err error)
	QueryByID(ctx context.Context, id int64, p *dbv1.RoomProto) error
	UpdateByID(ctx context.Context, id int64, user *dbv1.RoomProto) error
	Exist(ctx context.Context, id int64) (bool, error)
	IncVersion(ctx context.Context, id int64, newVersion int64) error
}

type RoomProtoCache interface {
	Put(ctx context.Context, uid int64, room *dbv1.RoomProto, ctime time.Time)
	Get(ctx context.Context, uid int64, ctime time.Time) *dbv1.RoomProto
	Create() *dbv1.RoomProto
	Reset(p *dbv1.RoomProto)
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

func (do *RoomDomain) UpdateOfflineCache(ctx context.Context, id int64, proto *dbv1.RoomProto, ctime time.Time) {
	do.cache.Put(ctx, id, proto, ctime)
}

func (do *RoomDomain) OfflineCache(ctx context.Context, id int64, ctime time.Time) *dbv1.RoomProto {
	return do.cache.Get(ctx, id, ctime)
}

func (do *RoomDomain) GetProtoFromPool() (p *dbv1.RoomProto) {
	return do.cache.Create()
}

func (do *RoomDomain) PutBackProtoIntoPool(p *dbv1.RoomProto) {
	do.cache.Reset(p)
}
