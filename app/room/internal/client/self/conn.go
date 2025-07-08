package self

import (
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/router/routetable/redis"
	"github.com/go-pantheon/roma/app/room/internal/conf"
	goredis "github.com/redis/go-redis/v9"
)

type SelfRouteTable struct {
	routetable.ReNewalRouteTable
}

func NewSelfRouteTable(db goredis.UniversalClient, c *conf.Data) *SelfRouteTable {
	return &SelfRouteTable{
		ReNewalRouteTable: routetable.NewRenewalRouteTable(
			redis.New(db),
			profile.ServiceName(),
			routetable.WithTTL(c.RouteTableAliveDuration.AsDuration()),
		),
	}
}
