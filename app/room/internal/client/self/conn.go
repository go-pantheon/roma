package self

import (
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/router/routetable/redis"
	"github.com/go-pantheon/roma/app/room/internal/data"
)

type SelfRouteTable struct {
	routetable.ReNewalRouteTable
}

func NewSelfRouteTable(db *data.Data) *SelfRouteTable {
	return &SelfRouteTable{
		ReNewalRouteTable: routetable.NewRenewalRouteTable(redis.New(db.Rdb), profile.ServiceName()),
	}
}
