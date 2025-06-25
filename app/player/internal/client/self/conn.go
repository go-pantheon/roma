package self

import (
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	xredis "github.com/go-pantheon/fabrica-kit/router/routetable/redis"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/redis/go-redis/v9"
)

type SelfRouteTable struct {
	routetable.ReNewalRouteTable
}

func NewSelfRouteTable(rdb redis.UniversalClient, c *conf.Data) *SelfRouteTable {
	return &SelfRouteTable{
		ReNewalRouteTable: routetable.NewRenewalRouteTable(
			xredis.New(rdb),
			profile.ServiceName(),
			routetable.WithTTL(c.RouteTableAliveDuration.AsDuration()),
		),
	}
}
