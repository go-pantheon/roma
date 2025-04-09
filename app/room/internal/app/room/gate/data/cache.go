package data

import (
	"context"
	"sync"
	"time"

	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	roomobj "github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain/object"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/go-pantheon/roma/pkg/util/lru"
	"google.golang.org/protobuf/proto"
)

var _ domain.RoomProtoCache = (*protoCache)(nil)

type protoCache struct {
	caches *lru.LRU
}

func NewProtoCache() domain.RoomProtoCache {
	caches := lru.NewLRU(
		lru.WithCapacity(constants.WorkerSize),
		lru.WithTTL(constants.CacheExpiredDuration),
		lru.WithOnRemove(func(key int64, value proto.Message) {
			putInPool(value.(*dbv1.RoomProto))
		}),
	)
	return &protoCache{caches: caches}
}

func (c *protoCache) Put(ctx context.Context, uid int64, room *dbv1.RoomProto, ctime time.Time) {
	c.caches.Put(uid, room, ctime)
}

func (c *protoCache) Get(ctx context.Context, uid int64, ctime time.Time) *dbv1.RoomProto {
	o, ok := c.caches.Get(uid, ctime)
	if !ok || o == nil {
		return nil
	}
	return o.(*dbv1.RoomProto)
}

func (c *protoCache) Create() *dbv1.RoomProto {
	return protoPool.Get().(*dbv1.RoomProto)
}

func (c *protoCache) Reset(p *dbv1.RoomProto) {
	putInPool(p)
}

var protoPool = sync.Pool{
	New: func() any {
		return roomobj.NewRoomProto()
	},
}

func putInPool(p *dbv1.RoomProto) {
	proto.Reset(p)
	protoPool.Put(p)
}
