package data

import (
	"context"
	"sync"
	"time"

	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	userobj "github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/object"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/go-pantheon/roma/pkg/util/lru"
	"google.golang.org/protobuf/proto"
)

var _ domain.UserOfflineCache = (*offlineCache)(nil)

type offlineCache struct {
	caches *lru.LRU
}

func NewProtoCache() domain.UserOfflineCache {
	return &offlineCache{
		caches: lru.NewLRU(
			lru.WithCapacity(constants.WorkerSize),
			lru.WithTTL(constants.CacheExpiredDuration),
			lru.WithOnRemove(func(key int64, value proto.Message) {
				putInPool(value.(*dbv1.UserProto))
			}),
		),
	}
}

func (c *offlineCache) Put(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time) {
	c.caches.Put(uid, user, ctime)
}

func (c *offlineCache) Get(ctx context.Context, uid int64, ctime time.Time) *dbv1.UserProto {
	o, ok := c.caches.Get(uid, ctime)
	if !ok || o == nil {
		return nil
	}
	return o.(*dbv1.UserProto)
}

func (c *offlineCache) Create() *dbv1.UserProto {
	return protoPool.Get().(*dbv1.UserProto)
}

func (c *offlineCache) Reset(p *dbv1.UserProto) {
	putInPool(p)
}

var protoPool = sync.Pool{
	New: func() any {
		return userobj.NewUserProto()
	},
}

func putInPool(p *dbv1.UserProto) {
	proto.Reset(p)
	protoPool.Put(p)
}
