package core

import (
	"context"
	"time"

	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/router/routetable/redis"
	redisgo "github.com/redis/go-redis/v9"
)

const serviceName = "janus.gate.interface"

type RouteTableManager struct {
	routeTable routetable.ReadOnlyRouteTable
	cache      *RouteTableCache
}

func NewRouteTableManager(rdb redisgo.UniversalClient) (*RouteTableManager, func()) {
	rt := routetable.NewReadOnlyRouteTable(redis.New(rdb), serviceName)
	cache := NewRouteTableCache(rt, serviceName, 10*time.Second)

	return &RouteTableManager{
			routeTable: rt,
			cache:      cache,
		}, func() {
			_ = cache.Stop(context.Background())
		}
}

func (m *RouteTableManager) GetEndpoint(ctx context.Context, uid int64, color string) (string, error) {
	return m.cache.get(ctx, uid, color)
}

func (m *RouteTableManager) GetEndpoints(ctx context.Context, uids []int64, color string) (map[string][]int64, error) {
	return m.cache.batchGet(ctx, uids, color)
}

func (m *RouteTableManager) Stop(ctx context.Context) error {
	return m.cache.Stop(ctx)
}
