package room

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-pantheon/fabrica-kit/router/balancer"
	"github.com/go-pantheon/fabrica-kit/router/conn"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/router/routetable/redis"
	"github.com/go-pantheon/roma/pkg/data/redisdb"
)

const (
	serviceName = "roma.room.service"
)

type Conn struct {
	*conn.Conn
}

func NewConn(logger log.Logger, rt *RouteTable, r registry.Discovery) (*Conn, error) {
	conn, err := conn.NewConn(serviceName, balancer.TypeReader, logger, rt, r)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Conn: conn,
	}, nil
}

type RouteTable struct {
	routetable.ReadOnlyRouteTable
}

func NewRouteTable(db *redisdb.DB) *RouteTable {
	return &RouteTable{
		ReadOnlyRouteTable: routetable.NewReadOnlyRouteTable(redis.New(db.DB), serviceName),
	}
}
