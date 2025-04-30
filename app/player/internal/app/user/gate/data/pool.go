package data

import (
	"context"
	"time"

	"github.com/go-pantheon/fabrica-util/multipool"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	userobj "github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/object"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/go-pantheon/roma/pkg/util/lru"
	"google.golang.org/protobuf/proto"
)

var (
	defaultThresholds = []int{256, 512, 1024, 4096, 16384}
)

var _ domain.UserProtoPool = (*UserProtoPool)(nil)

type UserProtoPool struct {
	pool  *multipool.MultiLayerPool
	cache *lru.LRU
}

func NewUserProtoPool() domain.UserProtoPool {
	pool := multipool.NewMultiLayerPool(
		func() multipool.Resetable {
			return userobj.NewUserProto()
		},
		func(obj multipool.Resetable) int {
			return userobj.UserProtoSize(obj.(*dbv1.UserProto))
		},
		multipool.WithThresholds(defaultThresholds),
	)

	c := &UserProtoPool{
		pool: pool,
		cache: lru.NewLRU(
			lru.WithCapacity(constants.WorkerSize),
			lru.WithTTL(constants.CacheExpiredDuration),
			lru.WithOnRemove(func(key int64, value proto.Message) {
				pool.Put(value.(*dbv1.UserProto))
			}),
		),
	}

	return c
}

func (c *UserProtoPool) Cache(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time) {
	c.cache.Put(uid, user, ctime)
}

func (c *UserProtoPool) Get() (ret *dbv1.UserProto, putFunc func()) {
	ret = c.pool.Get(defaultThresholds[0]).(*dbv1.UserProto)

	return ret, func() {
		c.pool.Put(ret)
	}
}

func (c *UserProtoPool) Load(ctx context.Context, uid int64, ctime time.Time) (ret *dbv1.UserProto) {
	o, ok := c.cache.Get(uid, ctime)
	if !ok {
		return nil
	}

	return o.(*dbv1.UserProto)
}

func (c *UserProtoPool) Put(user *dbv1.UserProto) {
	c.cache.Remove(user.Id)
	c.pool.Put(user)
}
