package self

import (
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/router/routetable/redis"
	"github.com/go-pantheon/roma/app/room/internal/conf"
	"github.com/go-pantheon/roma/pkg/data/redisdb"
)

type SelfRouteTable struct {
	routetable.ReNewalRouteTable
}

func NewSelfRouteTable(db *redisdb.DB, c *conf.Data) *SelfRouteTable {
	return &SelfRouteTable{
		ReNewalRouteTable: routetable.NewRenewalRouteTable(
			redis.New(db.DB),
			profile.ServiceName(),
			routetable.WithTTL(c.RouteTableAliveDuration.AsDuration()),
		),
	}
}
