package data

import (
	"context"
	"time"

	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/go-pantheon/roma/pkg/util/lru"
	"google.golang.org/protobuf/proto"
)

var _ domain.UserCache = (*UserProtoCache)(nil)

type UserProtoCache struct {
	cache *lru.LRU
}

func NewUserProtoCache() domain.UserCache {
	c := &UserProtoCache{
		cache: lru.NewLRU(
			lru.WithCapacity(constants.WorkerSize),
			lru.WithTTL(constants.CacheExpiredDuration),
			lru.WithOnRemove(func(key int64, value proto.Message) {
				user := value.(*dbv1.UserProto)
				dbv1.UserProtoPool.Put(user)
			}),
		),
	}

	return c
}

func (c *UserProtoCache) Put(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time) {
	c.cache.Put(uid, user, ctime)
}

func (c *UserProtoCache) Get(ctx context.Context, uid int64, ctime time.Time) (ret *dbv1.UserProto) {
	o, ok := c.cache.Get(uid, ctime)
	if !ok {
		return nil
	}

	return o.(*dbv1.UserProto)
}

func (c *UserProtoCache) Remove(ctx context.Context, uid int64) {
	c.cache.Remove(uid)
}
