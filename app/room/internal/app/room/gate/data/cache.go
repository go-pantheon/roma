package data

import (
	"context"
	"time"

	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/go-pantheon/roma/pkg/util/lru"
)

var _ domain.RoomProtoCache = (*protoCache)(nil)

type protoCache struct {
	caches *lru.LRU
}

func NewProtoCache() domain.RoomProtoCache {
	caches := lru.NewLRU(
		lru.WithCapacity(constants.WorkerSize),
		lru.WithTTL(constants.CacheExpiredDuration),
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

func (c *protoCache) Remove(ctx context.Context, uid int64) {
	c.caches.Remove(uid)
}
